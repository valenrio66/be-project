package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func StringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func PtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ToPgTimestamp(t *time.Time) pgtype.Timestamptz {
	if t != nil {
		return pgtype.Timestamptz{Time: *t, Valid: true}
	}
	return pgtype.Timestamptz{Valid: false}
}

func ToPgText(s *string) pgtype.Text {
	if s != nil {
		return pgtype.Text{String: *s, Valid: true}
	}
	return pgtype.Text{Valid: false}
}
