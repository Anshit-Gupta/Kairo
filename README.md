# Kairo

> A persistent, concurrent key-value store built from scratch in Go.

Kairo is a small storage engine built to explore the fundamentals behind persistent key-value databases.

It started as an in-memory `map[string]string` and evolved into a store with concurrency control, append-only persistence, crash recovery, and durable writes.

The project is intentionally small, with a focus on understanding the systems behind the abstractions.

---

## Features

- In-memory key-value storage
- `Set`, `Get`, and `Delete` operations
- Concurrent access with `sync.RWMutex`
- Append-only JSON log for persistence
- Automatic state recovery on startup
- Detection of corrupted middle log entries
- Handling of incomplete final log entries
- Durable writes using `File.Sync()`
- Configurable log file path
- Tests for persistence and recovery
- Race-detector compatible

---

## Architecture

```text
                    Kairo
                      │
             ┌────────┴────────┐
             │                 │
        In-Memory Store     Append-Only Log
       map[string]string          │
             │                    │
        Get / Set / Delete    Persistence
                                  │
                              Recovery
```

Writes follow:

```text
Operation → JSON Log → Sync → In-Memory State
```

On startup:

```text
Log → Replay → Reconstructed State
```

---

## Storage Format

Operations are stored as newline-delimited JSON:

```json
{"Op":"SET","Key":"name","Value":"Anshit"}
{"Op":"SET","Key":"language","Value":"Go"}
{"Op":"DELETE","Key":"language"}
```

The log is replayed on startup to reconstruct the current in-memory state.

---

## Crash Recovery

Kairo distinguishes between corrupted log entries and an incomplete final write.

If the final record is truncated because of a crash:

```json
{"Op":"SET","Key":"city","Value":
```

the incomplete final entry is ignored during recovery.

A malformed entry in the middle of the log causes recovery to fail rather than silently skipping potentially lost operations.

---

## Durability

Every successful mutation follows:

```text
Write → Sync → Update Memory
```

`File.Sync()` is used to provide a stronger durability guarantee before the operation is considered successful.

The current implementation synchronizes every mutation. Performance and synchronization strategies are being explored through benchmarks.

---

## Project Structure

```text
kairo/
├── cmd/
│   └── kairo/
│       └── main.go
├── internal/
│   └── store/
│       ├── store.go
│       ├── log.go
│       ├── store_test.go
│       └── store_bench_test.go
├── data/
├── go.mod
└── README.md
```

---

## Usage

```go
package main

import (
	"fmt"

	"github.com/Anshit-Gupta/kairo/internal/store"
)

func main() {
	s, err := store.NewStore("./data/kairo.log")
	if err != nil {
		panic(err)
	}
	defer s.Close()

	if err := s.Set("name", "Anshit"); err != nil {
		panic(err)
	}

	value, exists := s.Get("name")
	if exists {
		fmt.Println(value)
	}

	if err := s.Delete("name"); err != nil {
		panic(err)
	}
}
```

---

## Running

```sh
go run ./cmd/kairo
```

### Tests

```sh
go test ./...
```

### Race Detector

```sh
go test -race ./...
```

### Benchmarks

```sh
go test -bench=. ./internal/store
```

---


## Why Kairo?

Kairo is primarily a learning project.

The goal is to understand storage systems by building one instead of treating databases as black boxes.

The project focuses on going deeper into a smaller system rather than adding a large number of unrelated features.

---

