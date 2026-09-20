// Package settings implements FR-1 (initial learning settings): the
// User's Level, onboarding-completion state, and the active Topics
// visible in Settings.
//
// Every call currently operates on the single stub user (db.StubUserID)
// — see DEBT-1 for why, and the condition to remove it.
package settings

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/LinnikD/learn-the-grammar/backend/internal/db"
)

// Topic is an active vocabulary Topic, as shown in Settings.
type Topic struct {
	ID   uuid.UUID
	Name string
}

// Settings is a User's current learning settings.
type Settings struct {
	Level               string
	OnboardingCompleted bool
	Topics              []Topic
}

// Service reads and saves Settings.
type Service struct {
	queries *db.Queries
}

// NewService returns a Service backed by db.
func NewService(pool db.DBTX) *Service {
	return &Service{queries: db.New(pool)}
}

// Get returns userID's current settings.
func (s *Service) Get(ctx context.Context, userID uuid.UUID) (Settings, error) {
	row, err := s.queries.GetUserSettings(ctx, toPgUUID(userID))
	if err != nil {
		return Settings{}, fmt.Errorf("getting user settings: %w", err)
	}

	return s.withTopics(ctx, row.Level, row.OnboardingCompletedAt)
}

// Save sets userID's Level and completes onboarding, then returns the
// resulting settings. Onboarding, once completed, is never
// un-completed by a later Save.
func (s *Service) Save(ctx context.Context, userID uuid.UUID, level string) (Settings, error) {
	row, err := s.queries.SaveUserSettings(ctx, db.SaveUserSettingsParams{
		ID:    toPgUUID(userID),
		Level: level,
	})
	if err != nil {
		return Settings{}, fmt.Errorf("saving user settings: %w", err)
	}

	return s.withTopics(ctx, row.Level, row.OnboardingCompletedAt)
}

func (s *Service) withTopics(ctx context.Context, level string, onboardingCompletedAt pgtype.Timestamptz) (Settings, error) {
	rows, err := s.queries.ListActiveTopics(ctx)
	if err != nil {
		return Settings{}, fmt.Errorf("listing active topics: %w", err)
	}

	topics := make([]Topic, len(rows))
	for i, row := range rows {
		topics[i] = Topic{ID: uuid.UUID(row.ID.Bytes), Name: row.Name}
	}

	return Settings{
		Level:               level,
		OnboardingCompleted: onboardingCompletedAt.Valid,
		Topics:              topics,
	}, nil
}

func toPgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}
