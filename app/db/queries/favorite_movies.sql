-- name: ListMovies :many
SELECT * FROM favorite_movies;

-- name: AddMovie :one
INSERT INTO favorite_movies(tmdb_id, title) VALUES ($1, $2) RETURNING *;

