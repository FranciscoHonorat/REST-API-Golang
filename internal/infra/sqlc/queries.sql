-- name: GetBookByID :one
SELECT id, name, author, year, masterpiece, created_at, updated_at
FROM books
WHERE id = $1;

-- name: ListBooks :many
SELECT id, name, author, year, masterpiece, created_at, updated_at
FROM books
WHERE 1=1
  AND (CASE WHEN sqlc.narg('filter_author')::text IS NOT NULL THEN author ILIKE '%' || sqlc.narg('filter_author') || '%' ELSE true END)
  AND (CASE WHEN sqlc.narg('filter_year')::int IS NOT NULL THEN year = sqlc.narg('filter_year')::int ELSE true END)
  AND (CASE WHEN sqlc.narg('filter_masterpiece')::boolean IS NOT NULL THEN masterpiece = sqlc.narg('filter_masterpiece')::boolean ELSE true END)
ORDER BY 
  CASE WHEN @sort_by::text = 'name' THEN name END,
  CASE WHEN @sort_by::text = 'author' THEN author END,
  CASE WHEN @sort_by::text = 'masterpiece' THEN masterpiece::text END
LIMIT @page_limit OFFSET @page_offset;

-- name: CountBooks :one
SELECT COUNT(*)
FROM books
WHERE 1=1
  AND (CASE WHEN sqlc.narg('filter_author')::text IS NOT NULL THEN author ILIKE '%' || sqlc.narg('filter_author') || '%' ELSE true END)
  AND (CASE WHEN sqlc.narg('filter_year')::int IS NOT NULL THEN year = sqlc.narg('filter_year')::int ELSE true END)
  AND (CASE WHEN sqlc.narg('filter_masterpiece')::boolean IS NOT NULL THEN masterpiece = sqlc.narg('filter_masterpiece')::boolean ELSE true END);

-- name: CreateBook :one
INSERT INTO books (name, author, year, masterpiece)
VALUES ($1, $2, $3, $4)
RETURNING id, name, author, year, masterpiece, created_at, updated_at;

-- name: UpdateBook :one
UPDATE books
SET name = $1, author = $2, year = $3, masterpiece = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $5
RETURNING id, name, author, year, masterpiece, created_at, updated_at;

-- name: DeleteBook :execrows
DELETE FROM books
WHERE id = $1;