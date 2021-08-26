package tasks

import (
	"context"
	"database/sql"
)

type InSQL struct {
	db *sql.DB
}

func NewInSQL(db *sql.DB) (*InSQL, error) {
	const query = `CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY,
		list_id TEXT NOT NULL,
		text TEXT NOT NULL,
		completed INTEGER DEFAULT 0,
		created_at DATETIME
	)`

	_, err := db.Exec(query)

	if err != nil {
		return nil, err
	}

	return &InSQL{
		db: db,
	}, nil
}

func (sql *InSQL) CreateTask(ctx context.Context, task *Task) (*Task, error) {
	const query = "INSERT INTO tasks (id, list_id, text, completed, created_at) VALUES (?, ?, ?, ?, ?)"
	_, err := sql.db.ExecContext(ctx, query, task.ID, task.ListID, task.Text, task.Completed, task.CreatedAt)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (sql *InSQL) UpdateTask(ctx context.Context, task *Task) (*Task, error) {
	const query = "UPDATE tasks SET text = ?, completed = ? WHERE id = ?"
	_, err := sql.db.ExecContext(ctx, query, task.Text, task.Completed, task.ID)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (sql *InSQL) GetTask(ctx context.Context, id string) (*Task, error) {
	const query = "SELECT id, list_id, text, completed, created_at FROM tasks WHERE id = ?"

	task := &Task{}

	err := sql.db.QueryRowContext(ctx, query, id).Scan(&task.ID, &task.ListID, &task.Text, &task.Completed, &task.CreatedAt)

	if err != nil {
		return nil, err
	}

	return task, err
}

func (sql *InSQL) GetTasks(ctx context.Context, listsID string) ([]*Task, error) {
	const query = "SELECT id, list_id, text, completed, created_at FROM tasks WHERE list_id = ?"

	rows, err := sql.db.QueryContext(ctx, query, listsID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := []*Task{}

	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.ListID, &task.Text, &task.Completed, &task.CreatedAt)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (sql *InSQL) GetAllTasks(ctx context.Context) ([]*Task, error) {
	const query = "SELECT id, list_id, text, completed, created_at FROM tasks"

	rows, err := sql.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := []*Task{}

	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.ListID, &task.Text, &task.Completed, &task.CreatedAt)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (sql *InSQL) DeleteTask(ctx context.Context, taskID string) error {
	const query = "DELETE FROM tasks WHERE id = ?"
	_, err := sql.db.ExecContext(ctx, query, taskID)

	if err != nil {
		return err
	}

	return nil
}

func (sql *InSQL) DeleteTasks(ctx context.Context, listID string) error {
	const query = "DELETE FROM tasks WHERE list_id = ?"
	_, err := sql.db.ExecContext(ctx, query, listID)

	if err != nil {
		return err
	}

	return nil
}
