package constants

import "github.com/goravel/framework/support/color"

const (
	StatusPending = iota + 1
	StatusInProgress
	StatusCompleted
)

// StatusMap - Maps string status to integer constants
var StatusMap = map[string]int{
	"pending":     StatusPending,
	"in-progress": StatusInProgress,
	"completed":   StatusCompleted,
}

var StatusLabels = map[int]string{
	StatusPending:    "Pending",
	StatusInProgress: "In Progress",
	StatusCompleted:  "Completed",
}

var StatusColors = map[int]string{
	StatusPending:    color.Sprint("<fg=blue>Pending</>"),
	StatusInProgress: color.Sprint("<fg=cyan>In Progress</>"),
	StatusCompleted:  color.Sprint("<fg=magenta>Completed</>"),
}
