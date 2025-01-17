package repositories

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/kkumar-gcc/todo/models"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

// TaskRepository defines the methods that the Task repository should implement.
type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id int) error
	DeleteBulk(ctx context.Context, ids []int) error
	GetAll(ctx context.Context, status, priority int, sort string) ([]models.Task, error)
	GetByID(ctx context.Context, id int) (*models.Task, error)
	Update(ctx context.Context, id int, updateFunc func(task *models.Task) (*models.Task, error)) error
}

type TaskRepositoryImpl struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &TaskRepositoryImpl{
		db: db,
	}
}

func (r *TaskRepositoryImpl) Create(ctx context.Context, task *models.Task) error {
	query := `INSERT INTO tasks (title, status, completed_at, priority, tags)
              VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, task.Title, task.Status, task.CompletedAt, task.Priority, task.Tags)
	return err
}

func (r *TaskRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM tasks WHERE id = ?"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (r *TaskRepositoryImpl) DeleteBulk(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	query := "DELETE FROM tasks WHERE id IN (" + strings.Repeat("?,", len(ids)-1) + "?)"

	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (r *TaskRepositoryImpl) GetAll(ctx context.Context, status, priority int, sort string) ([]models.Task, error) {
	query := "SELECT id, title, status, created_at, completed_at, priority, tags FROM tasks WHERE 1=1"

	var args []any
	if status != 0 {
		query += " AND status = ?"
		args = append(args, status)
	}

	if priority != 0 {
		query += " AND priority = ?"
		args = append(args, priority)
	}

	if sort != "" {
		query += " ORDER BY " + sort
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Status, &task.CreatedAt, &task.CompletedAt, &task.Priority, &task.Tags)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskRepositoryImpl) GetByID(ctx context.Context, id int) (*models.Task, error) {
	query := "SELECT id, title, status, created_at, completed_at, priority, tags FROM tasks WHERE id = ?"
	row := r.db.QueryRowContext(ctx, query, id)

	var task models.Task
	err := row.Scan(&task.ID, &task.Title, &task.Status, &task.CreatedAt, &task.CompletedAt, &task.Priority, &task.Tags)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepositoryImpl) Update(ctx context.Context, id int, updateFunc func(task *models.Task) (*models.Task, error)) error {
	task, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if task == nil {
		return ErrTaskNotFound
	}

	updatedTask, err := updateFunc(task)
	if err != nil {
		return err
	}

	query := `UPDATE tasks SET title = ?, status = ?, completed_at = ?, priority = ?, tags = ? WHERE id = ?`
	_, err = r.db.ExecContext(ctx, query, updatedTask.Title, updatedTask.Status, updatedTask.CompletedAt, updatedTask.Priority, updatedTask.Tags, id)
	return err
}
