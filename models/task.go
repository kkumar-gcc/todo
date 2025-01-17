package models

import "time"

type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Status      int        `json:"status"` // Use constants: constants.StatusPending, constants.StatusInProgress, constants.StatusCompleted
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Priority    int        `json:"priority"` // Use constants: constants.PriorityLow, constants.PriorityMedium, constants.PriorityHigh
	Tags        string     `json:"tags"`     // Tags for categorization
}
