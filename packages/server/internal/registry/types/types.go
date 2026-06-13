package types

import "time"

type Suite struct {
	id          string
	name        string
	description string
	createdAt   time.Time
	updatedAt   time.Time
}