-- +goose Up
create table posts (
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    title text not null,
    url text not null unique,
    description text not null,
    published_at text not null,
    feed_id uuid not null references feeds (id) on delete cascade
);

-- +goose Down
drop table posts;
