package settings_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LinnikD/learn-the-grammar/backend/internal/db"
	"github.com/LinnikD/learn-the-grammar/backend/internal/db/dbtest"
	"github.com/LinnikD/learn-the-grammar/backend/internal/settings"
)

func TestGet_NewUserDefaults(t *testing.T) {
	pool := dbtest.StartPostgres(t)
	ctx := context.Background()

	_, err := pool.Exec(ctx, `INSERT INTO topics (id, name, active) VALUES ($1, 'Family', true), ($2, 'Retired', false)`,
		uuid.New(), uuid.New())
	require.NoError(t, err)

	svc := settings.NewService(pool)

	got, err := svc.Get(ctx, db.StubUserID)
	require.NoError(t, err)

	assert.Equal(t, "A1", got.Level)
	assert.False(t, got.OnboardingCompleted)
	require.Len(t, got.Topics, 1, "only the active Topic should be listed")
	assert.Equal(t, "Family", got.Topics[0].Name)
}

func TestGet_TopicActivatedAfterFirstReadIsEnabledAutomatically(t *testing.T) {
	pool := dbtest.StartPostgres(t)
	ctx := context.Background()

	svc := settings.NewService(pool)

	before, err := svc.Get(ctx, db.StubUserID)
	require.NoError(t, err)
	assert.Empty(t, before.Topics)

	_, err = pool.Exec(ctx, `INSERT INTO topics (id, name, active) VALUES ($1, 'Travel', true)`, uuid.New())
	require.NoError(t, err)

	after, err := svc.Get(ctx, db.StubUserID)
	require.NoError(t, err)
	require.Len(t, after.Topics, 1, "a Topic added after the first read must still show up, enabled")
	assert.Equal(t, "Travel", after.Topics[0].Name)
}

func TestSave_CompletesOnboardingOnceAndUpdatesLevel(t *testing.T) {
	pool := dbtest.StartPostgres(t)
	ctx := context.Background()

	svc := settings.NewService(pool)

	first, err := svc.Get(ctx, db.StubUserID)
	require.NoError(t, err)
	require.False(t, first.OnboardingCompleted)

	saved, err := svc.Save(ctx, db.StubUserID, "B1")
	require.NoError(t, err)
	assert.Equal(t, "B1", saved.Level)
	assert.True(t, saved.OnboardingCompleted)

	// A later Save changes the Level but must not un-complete onboarding.
	savedAgain, err := svc.Save(ctx, db.StubUserID, "A2")
	require.NoError(t, err)
	assert.Equal(t, "A2", savedAgain.Level)
	assert.True(t, savedAgain.OnboardingCompleted)
}
