package model

import "time"

type ActivityType string

const (
	ActivityAdded     ActivityType = "ADDED"
	ActivityModified  ActivityType = "MODIFIED"
	ActivityCompleted ActivityType = "COMPLETED"
	ActivityRemoved   ActivityType = "REMOVED"
)

type ActivityEvent struct {
	ID        string
	TaskID    string
	UserID    string
	Type      ActivityType
	Timestamp time.Time
	Details   string
}
