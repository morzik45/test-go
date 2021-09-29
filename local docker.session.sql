INSERT INTO results (test_id, percent)
VALUES (
        3,
        (
            SELECT c.count * 100 / a.count as percent
            FROM (
                    SELECT COUNT(a.id) as count
                    FROM user_answers AS a
                        RIGHT JOIN tasks as t ON (
                            t.id = a.task_id
                            AND t.correct = a.answer
                        )
                    WHERE test_id = 3
                ) as c,
                (
                    SELECT COUNT(id) as count
                    FROM tasks
                    WHERE variant_id = 2
                ) as a
        )
    )
RETURNING percent;