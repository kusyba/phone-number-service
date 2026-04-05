-- name: InsertPhone :exec
INSERT INTO phones (number, country, region, provider, source, created_at)
VALUES ($1, $2, $3, $4, $5, NOW())
ON CONFLICT (number) DO NOTHING;

-- name: CheckPhoneExists :one
SELECT EXISTS(SELECT 1 FROM phones WHERE number = $1);

-- name: GetPhonesWithFilters :many
SELECT id, number, country, region, provider, source, created_at
FROM phones
WHERE ($1 = '' OR number LIKE '%' || $1 || '%')
  AND ($2 = '' OR country = $2)
  AND ($3 = '' OR region = $3)
  AND ($4 = '' OR provider = $4)
ORDER BY id DESC
LIMIT $5 OFFSET $6;

-- name: CountPhonesWithFilters :one
SELECT COUNT(*)
FROM phones
WHERE ($1 = '' OR number LIKE '%' || $1 || '%')
  AND ($2 = '' OR country = $2)
  AND ($3 = '' OR region = $3)
  AND ($4 = '' OR provider = $4);
