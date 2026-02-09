-- name: CreateCampaign :one
INSERT INTO campaigns (
    user_id, title, description, status, start_date, end_date, budget
) VALUES (
             $1, $2, $3, $4, $5, $6, $7
         ) RETURNING id, user_id, title, description, status, budget, created_at;

-- name: GetCampaign :one
SELECT * FROM campaigns
WHERE id = $1 AND user_id = $2
    LIMIT 1;

-- name: ListCampaigns :many
SELECT id, user_id, title, description, status, start_date, end_date, budget, created_at
FROM campaigns
WHERE user_id = $1
ORDER BY created_at DESC
    LIMIT $2 OFFSET $3;

-- name: UpdateCampaign :one
UPDATE campaigns
SET
    title = COALESCE(sqlc.narg('title'), title),
    description = COALESCE(sqlc.narg('description'), description),
    status = COALESCE(sqlc.narg('status'), status),
    start_date = COALESCE(sqlc.narg('start_date'), start_date),
    end_date = COALESCE(sqlc.narg('end_date'), end_date),
    budget = COALESCE(sqlc.narg('budget'), budget),
    updated_at = NOW()
WHERE id = $1 AND user_id = $2
    RETURNING *;

-- name: DeleteCampaign :exec
DELETE FROM campaigns
WHERE id = $1 AND user_id = $2;