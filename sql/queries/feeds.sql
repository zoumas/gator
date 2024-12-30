-- name: CreateFeed :one
insert into feeds (id, created_at, updated_at, name, url, user_id)
values ($1, $2, $3, $4, $5, $6)
returning *;

-- name: GetFeeds :many
select
  f.name,
  f.url,
  u.name as owner
from feeds f
join users u on f.user_id = u.id;

-- name: GetFeedByURL :one
select *
from feeds
where url = $1;

-- name: MarkFeedFetched :exec
update feeds
set
  updated_at = now() at time zone 'utc',
  last_fetched_at = now() at time zone 'utc'
where id = $1;

-- name: GetNextFeedToFetch :one
select *
from feeds
order by last_fetched_at asc nulls first
limit 1;
