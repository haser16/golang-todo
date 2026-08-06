package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/haser16/golang-todo/internal/core/domain"
	core_postgres_pool "github.com/haser16/golang-todo/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) GetTask(
	ctx context.Context,
	id int,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeOut())
	defer cancel()

	query := `SELECT
    id, 
    version,
    title, 
    description, 
    completed, 
    created_at, 
    completed_at, 
    author_user_id
	FROM todoapp.tasks WHERE id = $1;`

	row := r.pool.QueryRow(ctx, query, id)

	var taskModel TaskModel

	err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CreatedAt,
		&taskModel.CompletedAt,
		&taskModel.AuthorUserID,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf("task with id='%d': %w", id, err)
		}
		return domain.Task{}, fmt.Errorf("scan row: %w", err)
	}
	taskDomain := taskDomainFromModel(taskModel)
	return taskDomain, nil
}
