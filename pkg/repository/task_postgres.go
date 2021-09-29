package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	exam "github.com/morzik45/test-go"
)

type TaskPostgres struct {
	db             *sql.DB
	saveAnswerLock *sync.Mutex
}

func NewTaskPostgres(db *sql.DB, saveAnswerLock *sync.Mutex) *TaskPostgres {
	return &TaskPostgres{db: db, saveAnswerLock: saveAnswerLock}
}

func (r *TaskPostgres) GetTaskById(variantId, taskId int, username string) (exam.Task, error) {
	task := exam.Task{
		VariantID: variantId,
		Id:        taskId,
	}
	var answersStr string
	var resultID, testID, userID sql.NullInt64
	query := fmt.Sprintf(`SELECT
							q.question question,
							q.answers answers,
							t.id t_id,
							r.id r_id,
							u.id u_id
						FROM (
								SELECT variant_id,
									question,
									answers
								FROM %s
								WHERE variant_id = $1
								ORDER BY id OFFSET $2
								LIMIT 1
							) q
							LEFT JOIN %s u ON (u.username = $3)
							LEFT JOIN %s t ON (
								t.id = (
									SELECT id
									from %s
									WHERE variant_id = q.variant_id
										AND user_id = u.id
									ORDER BY start_at DESC
									LIMIT 1
								)
							)
							LEFT JOIN %s r ON (r.test_id = t.id)
						LIMIT 1;`, tasksTable, usersTable, testsTable, testsTable, resultsTable)
	row := r.db.QueryRow(query, task.VariantID, taskId-1, username)
	if err := row.Scan(&task.Question, &answersStr, &testID, &resultID, &userID); err != nil {
		return task, err
	}
	json.Unmarshal([]byte(answersStr), &task.Answers)

	if testID.Valid && !resultID.Valid { // если тест с таким вариантом начат, но не закончен
		task.TestID = int(testID.Int64)
	} else {
		startQuery := fmt.Sprintf("INSERT INTO %s (user_id, variant_id) VALUES ($1, $2) RETURNING id;", testsTable)
		row := r.db.QueryRow(startQuery, userID.Int64, variantId)
		if err := row.Scan(&task.TestID); err != nil {
			return task, err
		}
	}
	return task, nil
}

func (r *TaskPostgres) SaveAnswer(answer exam.Answer, username string) (bool, int, error) {
	r.saveAnswerLock.Lock()
	defer r.saveAnswerLock.Unlock()

	var finished bool
	checkQuery := fmt.Sprintf(`
			SELECT q.id AS q_id,
			u.id AS u_id,
			t.id AS t_id,
			r.id AS r_id,
			a.id AS a_id
		FROM (
				SELECT id, variant_id
				FROM %s
				WHERE variant_id = $1
				ORDER BY id OFFSET $2
				LIMIT 1
			) AS q
			LEFT JOIN %s AS u ON (u.username = $3)
			LEFT JOIN %s AS t ON (
				t.user_id = u.id
				AND t.variant_id = $1
				AND t.id = $4
			)
			LEFT JOIN %s AS r ON (r.test_id = t.id)
			LEFT JOIN %s AS a ON (a.test_id = t.id AND a.task_id = q.id)
		LIMIT 1;
	`, tasksTable, usersTable, testsTable, resultsTable, userAnswersTable)
	row := r.db.QueryRow(checkQuery, answer.VariantID, answer.TaskID-1, username, answer.TestID)
	var q_id, u_id, t_id, r_id, a_id sql.NullInt64
	if err := row.Scan(&q_id, &u_id, &t_id, &r_id, &a_id); err != nil {
		return finished, 0, err
	}
	if !t_id.Valid || !u_id.Valid {
		return finished, 0, errors.New("requested test not started")
	}
	if r_id.Valid {
		return finished, 0, errors.New("test already finished")
	}
	if a_id.Valid {
		return finished, 0, errors.New("question already answered")
	}
	r.db.QueryRow(fmt.Sprintf("INSERT INTO %s (test_id, task_id, answer) VALUES ($1, $2, $3)", userAnswersTable), t_id, q_id, answer.AnswerID)

	row = r.db.QueryRow(fmt.Sprintf(`
		SELECT *
		FROM (
				SELECT COUNT(id)
				FROM %s
				WHERE variant_id = $1
			) AS t,
			(
				SELECT COUNT(id)
				FROM %s
				WHERE test_id = $2
			) AS a
	`, tasksTable, userAnswersTable), answer.VariantID, answer.TestID)
	var tasks_count, answers_count sql.NullInt64
	if err := row.Scan(&tasks_count, &answers_count); err != nil {
		return finished, 0, err
	}
	if tasks_count.Valid && answers_count.Valid && tasks_count.Int64 == answers_count.Int64 {
		row = r.db.QueryRow(fmt.Sprintf(`
			INSERT INTO %s (test_id, percent)
			VALUES (
					$1,
					(
						SELECT c.count * 100 / a.count as percent
						FROM (
								SELECT COUNT(a.id) as count
								FROM %s AS a
									RIGHT JOIN %s as t ON (
										t.id = a.task_id
										AND t.correct = a.answer
									)
								WHERE test_id = $1
							) as c,
							(
								SELECT COUNT(id) as count
								FROM %s
								WHERE variant_id = $2
							) as a
					)
				)
			RETURNING percent;
		`, resultsTable, userAnswersTable, tasksTable, tasksTable), answer.TestID, answer.VariantID)
		var percent sql.NullInt64
		if err := row.Scan(&percent); err != nil {
			return finished, 0, err
		}
		if percent.Valid {
			finished = true
			return finished, int(percent.Int64), nil
		}
	}
	return finished, 0, nil
}

func (r *TaskPostgres) CreateVariant(name string) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name) values ($1) RETURNING id", variantsTable)
	row := r.db.QueryRow(query, name)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *TaskPostgres) GetAllVariants() ([]exam.Variant, error) {
	query := fmt.Sprintf(`SELECT id, name FROM %s`, variantsTable)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var variants []exam.Variant
	for rows.Next() {
		var variant exam.Variant
		if err := rows.Scan(&variant.Id, &variant.Name); err != nil {
			return variants, err
		}
		variants = append(variants, variant)
	}
	if err := rows.Err(); err != nil {
		return variants, err
	}
	return variants, nil
}

func (r *TaskPostgres) StartTest(userId, variantId int) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (user_id, variant_id) values ($1, $2) RETURNING id", testsTable)

	row := r.db.QueryRow(query, userId, variantId)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
