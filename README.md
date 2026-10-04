# Kairo

> A persistent, concurrent key-value store built from scratch in Go.

Kairo is a lightweight storage engine built to explore the fundamentals of
persistent key-value databases.

It started as an in-memory map[string]string and evolved into a store with concurrency control, append-only persistence, crash recovery, durable writes, and log compaction.

The project is intentionally small, with a focus on understanding the systems behind the abstractions.

## Features

- In-memory key-value storage
- `Set`, `Get`, and `Delete` operations
- Concurrent access with `sync.RWMutex`
- Append-only JSON log persistence
- Automatic recovery after restart
- Durable writes with `File.Sync()`
- Log compaction
- Configurable log file path

## Architecture

```text
                           Kairo
                             │
                ┌────────────┴────────────┐
                │                         │
         In-Memory Store            Append-Only Log
        map[string]string                  │
                │                          │
        Get / Set / Delete             Persistence
                                           │
                                       Recovery
                                           │
                                      Compaction
```

### Write flow

```text
Operation → Append JSON Record → Sync Log → Update Memory
```

### Recovery flow

```text
Log File → Replay Records → Reconstructed In-Memory State
```

### Compaction flow

```text
Current State → Temporary Log → Sync → Replace Old Log
```

## Usage

```go
s, err := store.NewStore("./data/kairo.log")
if err != nil {
	log.Fatal(err)
}
defer s.Close()

if err := s.Set("name", "Kairo"); err != nil {
	log.Fatal(err)
}

value, ok := s.Get("name")
if ok {
	fmt.Println(value)
}
```

## Project Structure

```text
cmd/kairo/          Application entry point
internal/store/     Storage engine implementation and tests
data/               Local data directory
```

## Development

```sh
# Run
go run ./cmd/kairo

# Test
go test ./...

# Test with the race detector
go test -race ./...

# Run benchmarks
go test -run=^$ -bench=. -benchmem ./internal/store
```

This project is currently intended for learning and experimentation.
