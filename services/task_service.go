package services

import (
	"context"
	"errors"

	"github.com/kkumar-gcc/todo/models"
	"github.com/kkumar-gcc/todo/repositories"
)

var (
	ErrInvalidID          = errors.New("id must be a positive integer")
	ErrEmptyTitle         = errors.New("title cannot be empty")
	ErrInvalidStatus      = errors.New("status must be a non-negative integer")
	ErrInvalidPriority    = errors.New("priority must be a non-negative integer")
	ErrTaskNotFound       = errors.New("task not found")
	ErrTaskCreationFailed = errors.New("failed to create task")
	ErrTaskUpdateFailed   = errors.New("failed to update task")
	ErrTaskDeleteFailed   = errors.New("failed to delete task")
)

type TaskService interface {
	CreateTask(ctx context.Context, title string, status, priority int, tags string) error
	DeleteTask(ctx context.Context, id int) error
	DeleteTasks(ctx context.Context, ids []int) error
	GetAllTasks(ctx context.Context, status, priority int, sort string) ([]models.Task, error)
	GetTaskByID(ctx context.Context, id int) (*models.Task, error)
	UpdateTask(ctx context.Context, id int, updateFunc func(task *models.Task) (*models.Task, error)) error
}

type TaskServiceImpl struct {
	repository repositories.TaskRepository
}

// NewTaskService creates a new instance of TaskService
func NewTaskService(repo repositories.TaskRepository) TaskService {
	return &TaskServiceImpl{
		repository: repo,
	}
}

func (r *TaskServiceImpl) CreateTask(ctx context.Context, title string, status, priority int, tags string) error {
	if title == "" {
		return ErrEmptyTitle
	}
	if status < 0 {
		return ErrInvalidStatus
	}
	if priority < 0 {
		return ErrInvalidPriority
	}

	task := &models.Task{
		Title:    title,
		Status:   status,
		Priority: priority,
		Tags:     tags,
	}

	if err := r.repository.Create(ctx, task); err != nil {
		return ErrTaskCreationFailed
	}

	return nil
}

func (r *TaskServiceImpl) DeleteTask(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrInvalidID
	}

	if err := r.repository.Delete(ctx, id); err != nil {
		return ErrTaskDeleteFailed
	}

	return nil
}

func (r *TaskServiceImpl) DeleteTasks(ctx context.Context, ids []int) error {
	if len(ids) <= 0 {
		return ErrInvalidID
	}

	if err := r.repository.DeleteBulk(ctx, ids); err != nil {
		return ErrTaskDeleteFailed
	}

	return nil
}

func (r *TaskServiceImpl) GetAllTasks(ctx context.Context, status, priority int, sort string) ([]models.Task, error) {
	tasks, err := r.repository.GetAll(ctx, status, priority, sort)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskServiceImpl) GetTaskByID(ctx context.Context, id int) (*models.Task, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}

	task, err := r.repository.GetByID(ctx, id)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	return task, nil
}

func (r *TaskServiceImpl) UpdateTask(ctx context.Context, id int, updateFunc func(task *models.Task) (*models.Task, error)) error {
	if id <= 0 {
		return ErrInvalidID
	}

	if err := r.repository.Update(ctx, id, updateFunc); err != nil {
		return ErrTaskUpdateFailed
	}

	return nil
}
