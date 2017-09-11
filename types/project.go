package types

import "time"

// Project ci project
type Project struct {
	Name   string  `json:"name,omitempty"`
	Repo   string  `json:"repo,omitempty"`
	Builds []Build `json:"builds,omitempty"`
}

// Build build entity
type Build struct {
	ID        int       `json:"id"`
	StartTime time.Time `json:"starting_time,omitempty"`
	Duration  string    `json:"diration,omitempty"`
}
