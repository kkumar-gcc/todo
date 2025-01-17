package constants

import "github.com/goravel/framework/support/color"

const (
	PriorityLow = iota + 1
	PriorityMedium
	PriorityHigh
)

// PriorityMap - Maps string priority to integer constants
var PriorityMap = map[string]int{
	"low":    PriorityLow,
	"medium": PriorityMedium,
	"high":   PriorityHigh,
}

var PriorityLabels = map[int]string{
	PriorityLow:    "Low Priority",
	PriorityMedium: "Medium Priority",
	PriorityHigh:   "High Priority",
}

var PriorityColors = map[int]string{
	PriorityLow:    color.Sprint("<fg=green>Low</>"),
	PriorityMedium: color.Sprint("<fg=yellow>Medium</>"),
	PriorityHigh:   color.Sprint("<fg=red>High</>"),
}
