# gator 🐊

This project is part of [boot.dev](https://www.boot.dev/courses/build-blog-aggregator-golang)

## Disclaimer

This project is for learning purposes, it's not meant to be cloned.

While this was great practice, I don't see myself using the program.
If I start using RSS seriously perhaps I would create a simple single-user viewer,
but definitely not a long-running service.

## What is this?

A multi-user command line RSS feed aggregator and viewer.

A user can add a feed to be processed and the
feed's items will be made into posts that can be browsed.
Users can follow feeds that other users have added.

## Primary Learning Goals

- Interact with a SQL Database ([PostgreSQL](https://www.postgresql.org/))
  - Use [goose](https://pressly.github.io/goose/) for schema migrations
  - Use [sqlc](https://sqlc.dev/) to generate type-safe Go code for queries
- Modern command line subcommands approach
  - Build a system for handling subcommands from scratch
- RSS
  - Fetch and parse RSS feeds
  - Mark feeds for aggregation and do work periodically
