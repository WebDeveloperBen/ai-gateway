package repository

import "github.com/google/uuid"

func DerefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ParseUUID(s string) uuid.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return u
}
