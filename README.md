# Kairo

> A persistent, concurrent key-value store built from scratch in Go.

Kairo is a lightweight storage engine built to explore the fundamentals of
persistent key-value databases.

It started as an in-memory `map[string]string` and evolved into a store with
concurrency control, append-only persistence, crash recovery, durable writes,
log compaction, and an HTTP API.

The project is intentionally small, with a focus on understanding the systems
behind the abstractions.

## Features

- In-memory key-value storage
- `Set`, `Get`, and `Delete` operations
- Concurrent access with `sync.RWMutex`
- Append-only JSON log persistence
- Automatic recovery after restart
- Durable writes with `File.Sync()`
- Log compaction with safe log replacement
- Configurable log file path
- HTTP API
- Graceful HTTP server shutdown
- Unit and HTTP tests
- Race detector testing
- Benchmarks

## Architecture

```text
                         Kairo
                           │
              ┌────────────┴────────────┐
              │                         │
       In-Memory Store            Append-Only Log
       map[string]string                  │
              │                           │
       Get / Set / Delete            Persistence
              │                           │
              │                        Recovery
              │                           │
              │                       Compaction
              │
              └────────── HTTP API ──────────┐
                                             │
                                  GET /{key}  │
                                  PUT /{key}  │
                               DELETE /{key}  │
                              POST /compact   │
```

### Write flow

```text
Operation → Append JSON Record → Sync Log → Update Memory
```

A write is persisted to the append-only log and synchronized before the
in-memory state is updated.

### Recovery flow

```text
Log File → Replay Records → Reconstructed In-Memory State
```

When Kairo starts, it replays the log to reconstruct the current state.

### Compaction flow

```text
Current State → Temporary Log → Sync → Replace Old Log
```

Compaction removes obsolete historical records and keeps only the current
state.

## Persistence

Kairo uses an append-only JSON log.

For example:

```text
{"Op":"SET","Key":"name","Value":"Kairo"}
{"Op":"SET","Key":"language","Value":"Go"}
{"Op":"DELETE","Key":"name"}
```

On startup, these records are replayed in order to rebuild the in-memory store.

Kairo treats a malformed final log entry as a potentially incomplete write and
ignores it, while malformed entries in the middle of the log cause recovery to
fail.

## Concurrency

The in-memory store is protected using `sync.RWMutex`.

- `Get` uses a read lock
- `Set` and `Delete` use a write lock
- Compaction uses a write lock

This allows multiple concurrent readers while keeping mutations safe.

## HTTP API

Kairo exposes a simple HTTP API.

| Method | Endpoint | Description |
| --- | --- | --- |
| `PUT` | `/{key}` | Create or update a key |
| `GET` | `/{key}` | Retrieve a value |
| `DELETE` | `/{key}` | Delete a key |
| `POST` | `/compact` | Compact the append-only log |

### Examples

Set a value:

```sh
curl -X PUT http://localhost:8080/name -d "Anshit"
```

Get a value:

```sh
curl http://localhost:8080/name
```

Delete a value:

```sh
curl -X DELETE http://localhost:8080/name
```

Compact the log:

```sh
curl -X POST http://localhost:8080/compact
```

## Project Structure

```text
cmd/kairo/          Application entry point
internal/store/     Storage engine implementation and tests
internal/server/    HTTP server implementation and tests
data/               Local data directory
```

## Development

Run locally:

```sh
go run ./cmd/kairo
```

Run all tests:

```sh
go test ./...
```

Run with the race detector:

```sh
go test -race ./...
```

Run benchmarks:

```sh
go test -run=^$ -bench=. -benchmem ./internal/store
```

## Benchmarks

Benchmarks cover `Get` and `Set` operations. `Set` includes JSON encoding, log
writing, and `Sync()`, so its cost is dominated by persistence rather than the
in-memory map operation.

Results from the development machine:

| Operation | Time | Memory | Allocations |
| --- | ---: | ---: | ---: |
| `Get` | 15.20 ns/op | 0 B/op | 0 allocs/op |
| `Set` | 0.9 ms/op | 147 B/op | 3 allocs/op |

Benchmark results are machine and filesystem dependent and are intended for
comparison rather than absolute performance claims.

## Deployment

Kairo can be deployed as a Go web service. It reads the `PORT` environment
variable and binds to all interfaces.

The current free deployment uses ephemeral storage, so data written to the
local log may be lost after a restart or redeploy.

## Why Kairo?

Kairo was built while learning Go to understand backend systems by building a
concurrent, persistent key-value store from scratch.

The goal is to explore the engineering decisions behind storage systems rather
than to build a production database.
