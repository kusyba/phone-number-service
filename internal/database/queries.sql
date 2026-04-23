-- name: InsertPhone :exec
INSERT INTO phones (number, country, region, provider, source, created_at)
VALUES ($1, $2, $3, $4, $5, NOW())
ON CONFLICT (number) DO NOTHING;

-- name: GetPhonesWithFilters :many
SELECT id, number, country, region, provider, source, created_at
FROM phones
WHERE (true)
  AND ($1::text = '' OR number LIKE '%' || $1 || '%')
  AND ($2::text = '' OR country = $2)
  AND ($3::text = '' OR region = $3)
  AND ($4::text = '' OR provider = $4)
ORDER BY id DESC
LIMIT $5 OFFSET $6;

-- name: CountPhonesWithFilters :one
SELECT COUNT(*)
FROM phones
WHERE (true)
  AND ($1::text = '' OR number LIKE '%' || $1 || '%')
  AND ($2::text = '' OR country = $2)
  AND ($3::text = '' OR region = $3)
  AND ($4::text = '' OR provider = $4);
