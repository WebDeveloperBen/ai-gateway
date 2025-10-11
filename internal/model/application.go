package model

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID          uuid.UUID
	OrgID       uuid.UUID
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ApplicationConfig struct {
	ID          uuid.UUID
	AppID       uuid.UUID
	OrgID       uuid.UUID
	Environment string
	Config      map[string]interface{} // JSON blob
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
