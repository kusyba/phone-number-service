-- name: CreateGroup :one
INSERT INTO groups (name, description, flags)
VALUES ($1, $2, $3)
RETURNING id, name, description, flags, created_at, updated_at;

-- name: GetGroupByID :one
SELECT id, name, description, flags, created_at, updated_at
FROM groups
WHERE id = $1;

-- name: GetGroupsWithFilters :many
SELECT id, name, description, flags, created_at, updated_at
FROM groups
WHERE 
    ($1::text = '' OR name ILIKE '%' || $1 || '%')
    AND ($2::text = '' OR name ILIKE '%' || $2 || '%' OR description ILIKE '%' || $2 || '%')
ORDER BY 
    CASE WHEN $3 = 'name_asc' THEN name END ASC,
    CASE WHEN $3 = 'name_desc' THEN name END DESC,
    CASE WHEN $3 = 'created_at_asc' THEN created_at END ASC,
    CASE WHEN $3 = 'created_at_desc' THEN created_at END DESC,
    CASE WHEN $3 = '' THEN id END DESC
LIMIT $4 OFFSET $5;

-- name: CountGroupsWithFilters :one
SELECT COUNT(*)
FROM groups
WHERE 
    ($1::text = '' OR name ILIKE '%' || $1 || '%')
    AND ($2::text = '' OR name ILIKE '%' || $2 || '%' OR description ILIKE '%' || $2 || '%');

-- name: UpdateGroup :exec
UPDATE groups
SET name = $2, description = $3, flags = $4
WHERE id = $1;

-- name: DeleteGroup :exec
DELETE FROM groups WHERE id = $1;

-- name: AddUserToGroup :exec
INSERT INTO user_groups (user_id, group_id)
VALUES ($1, $2)
ON CONFLICT (user_id, group_id) DO NOTHING;

-- name: RemoveUserFromGroup :exec
DELETE FROM user_groups
WHERE user_id = $1 AND group_id = $2;

-- name: GetUserGroups :many
SELECT g.id, g.name, g.description, g.flags, g.created_at, g.updated_at
FROM groups g
JOIN user_groups ug ON g.id = ug.group_id
WHERE ug.user_id = $1
ORDER BY g.name;
