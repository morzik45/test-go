package repository

import (
	"database/sql"
	"fmt"

	exam "github.com/morzik45/test-go"
)

type TaskPostgres struct {
	db *sql.DB
}

func NewTaskPostgres(db *sql.DB) *TaskPostgres {
	return &TaskPostgres{db: db}
}

func (r *TaskPostgres) GetTaskById(variantId, taskId int) (exam.Task, error) {
	task := exam.Task{
		VariantID: variantId,
	}
	query := fmt.Sprintf("SELECT id, question, correct, answers FROM %s WHERE variant_id=$1 ORDER BY id OFFSET $2 LIMIT 1", tasksTable)
	row := r.db.QueryRow(query, task.VariantID, taskId-1)
	if err := row.Scan(&task.Id, &task.Question, &task.Correct, &task.Answers); err != nil {
		return task, err
	}
	return task, nil
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
