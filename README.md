<img src="/slot-logo.png">

# Slot

A minimal Redis-compatible in-memory key-value store written in Go.

## What is this

`slot` implements a subset of the [Redis protocol](https://redis.io/docs/reference/protocol-spec/) (RESP) and provides an in-memory store with string and hash data structures. Built as a learning exercise to understand how Redis works under the hood.

## Features

- **RESP parser** — handles Redis Serialization Protocol directly
- **String commands** — `SET`, `GET`
- **Hash commands** — `HSET`, `HGET`, `HGETALL`
- **PING** — connectivity check
- **Thread-safe** — uses `sync.RWMutex` for concurrent access

## Usage

```bash
go run main.go
```

Listens on port `6379` by default (Redis default). Override with the `PORT` environment variable:

```bash
PORT=6380 go run main.go
```

## Test with redis-cli

```bash
# ping
redis-cli -p 6379 PING

# string operations
redis-cli -p 6379 SET foo bar
redis-cli -p 6379 GET foo

# hash operations
redis-cli -p 6379 HSET myhash key1 value1
redis-cli -p 6379 HGET myhash key1
redis-cli -p 6379 HGETALL myhash
```

## Status

This is an educational implementation — not production-ready. Missing persistence, replication, clustering, and most Redis features. But it works.
