package db

import "github.com/google/uuid"

// StubUserID identifies the single pre-seeded user (see migration
// 00002_initial_settings.sql) that every request currently operates
// on. There is no real sign-up yet (FR-29 is not implemented), so
// features that need a persisted User use this fixed ID instead of a
// per-visitor account — see DEBT-1 for the reasoning and the
// condition to remove it.
var StubUserID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
