package lists

import (
	"context"
	"database/sql"
)

type InSQL struct {
	db *sql.DB
}

func NewInSQL(db *sql.DB) (*InSQL, error) {
	const query = `CREATE TABLE IF NOT EXISTS lists (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
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

func (sql *InSQL) CreateList(ctx context.Context, list *List) (*List, error) {
	const query = "INSERT INTO lists (id, name, created_at) VALUES (?, ?, ?)"
	_, err := sql.db.ExecContext(ctx, query, list.ID, list.Name, list.CreatedAt)

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (sql *InSQL) UpdateList(ctx context.Context, list *List) (*List, error) {
	const query = "UPDATE lists SET name = ? WHERE id = ?"
	_, err := sql.db.ExecContext(ctx, query, list.Name, list.ID)

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (sql *InSQL) GetList(ctx context.Context, id string) (*List, error) {
	const query = "SELECT id, name, created_at FROM lists WHERE id = ?"
	list := &List{}

	err := sql.db.QueryRowContext(ctx, query, id).Scan(&list.ID, &list.Name, &list.CreatedAt)

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (sql *InSQL) GetLists(ctx context.Context) ([]*List, error) {
	const query = "SELECT id, name, created_at FROM lists"
	lists := []*List{}

	rows, err := sql.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		list := &List{}
		err := rows.Scan(&list.ID, &list.Name, &list.CreatedAt)

		if err != nil {
			return nil, err
		}

		lists = append(lists, list)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return lists, nil
}

func (sql *InSQL) DeleteList(ctx context.Context, id string) error {
	const query = "DELETE FROM lists WHERE id = ?"
	_, err := sql.db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}

func (sql *InSQL) DeleteLists(ctx context.Context) error {
	const query = "DELETE FROM lists"
	_, err := sql.db.ExecContext(ctx, query)

	if err != nil {
		return err
	}

	return nil

}
