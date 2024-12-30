-- name: CreateFeedFollow :one
with inserted_feed_follow as (
  insert into feed_follows (id, created_at, updated_at, user_id, feed_id)
  values ($1, $2, $3, $4, $5)
  returning *
)
select 
  ff.*,
  f.name as feed_name,
  u.name as user_name
from inserted_feed_follow ff
join feeds f on ff.feed_id = f.id
join users u on ff.user_id = u.id;

-- name: GetFeedFollowsForUser :many
select
  ff.*,
  f.name as feed_name,
  u.name as user_name
from feed_follows ff
join users u on ff.user_id = u.id
join feeds f on ff.feed_id = f.id
where ff.user_id = $1;

-- name: UnfollowFeedForUser :exec
delete from feed_follows ff
using feeds f
where ff.feed_id = f.id
and ff.user_id = $1
and f.url = $2;
