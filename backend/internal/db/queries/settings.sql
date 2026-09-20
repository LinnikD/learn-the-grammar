-- name: GetUserSettings :one
SELECT level, onboarding_completed_at
FROM users
WHERE id = $1;

-- name: SaveUserSettings :one
-- Onboarding, once completed, stays completed: a later Save never
-- clears onboarding_completed_at.
UPDATE users
SET level = $2,
    onboarding_completed_at = COALESCE(onboarding_completed_at, now())
WHERE id = $1
RETURNING level, onboarding_completed_at;

-- name: ListActiveTopics :many
SELECT id, name
FROM topics
WHERE active
ORDER BY name;
