package store

//logic for set,get,store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Store struct {
	logFile *os.File //Logfile for persistent data
	mu      sync.RWMutex
	data    map[string]string
}

// constructor
func NewStore(path string) (*Store, error) {
	//if /data does not exist ->create data dir , if exist -> do nothing
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	//open or create Kairo .log
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	// Initialize the empty store
	s := &Store{
		logFile: file,
		data:    make(map[string]string),
	}

	//Replay the log into the map
	if err := s.logReplay(path); err != nil {
		file.Close() //since relay failed we close the file before returning err
		return nil, err
	}

	//return the recovered map
	return s, nil

}

// read
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, exist := s.data[key]
	return value, exist
}

// writre - SET
func (s *Store) Set(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	//update the log
	newlog := logEntry{Op: "SET", Key: key, Value: value}
	newjson, err := json.Marshal(newlog)
	// we append '\n' so that we always get the log in a new line

	if err != nil {
		fmt.Println("error converting newlog into json ")
		return err
	}
	newjson = append(newjson, '\n')

	//now we append the newjson into the file
	byteCount, err := s.logFile.Write(newjson) //.Write returns int n,err error ,here n is how many bytes were actully written

	if err != nil {
		fmt.Println("error Writing logFile")
		return err
	}

	//chances are very low but it is possible that short write than requested so we cross check :
	if byteCount != len(newjson) {
		fmt.Println("Short Write")
		return fmt.Errorf("short write")
	}

	//sync -> durability
	err = s.logFile.Sync()
	if err != nil {
		fmt.Println("Error in Sync")
		return fmt.Errorf("Error in Sync : %w", err)
	}

	//make sure the log is update and only then we update in-memory map

	s.data[key] = value
	return nil
}

// write : since it also modifies the map
func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	newlog := logEntry{Op: "DELETE", Key: key}
	newjson, err := json.Marshal(newlog)
	if err != nil {
		fmt.Println("error converting newlog into json ")
		return err
	}
	newjson = append(newjson, '\n')

	byteCount, err := s.logFile.Write(newjson)

	if err != nil {
		fmt.Println("error Writing logFile")
		return err
	}

	if byteCount != len(newjson) {
		fmt.Println("Short Write")
		return fmt.Errorf("short Write")
	}
	//now that we confimed the log file is updated

	//sync before updating the map
	err = s.logFile.Sync()
	if err != nil {
		fmt.Println("Error in Sync")
		return fmt.Errorf("Error in Sync : %w", err)
	}

	delete(s.data, key)
	return nil
}

//since we are keeping the file open while the program is running , it's better to ahve close fun

func (s *Store) Close() error {
	return s.logFile.Close()
}

// operations for logreplay
func (s *Store) logReplayOperationHelper(entry logEntry) error {
	if entry.Op == "SET" {
		s.data[entry.Key] = entry.Value
	} else if entry.Op == "DELETE" {
		//Delete
		delete(s.data, entry.Key)
	} else {
		//unknown operation or typo
		return fmt.Errorf("Unknown Operation detected %q", entry.Op)
	}

	return nil
}

//Log compaction

func (s *Store) Compact() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tempPath := s.logFile.Name() + ".tmp"
	logPath := s.logFile.Name()

	tempFile, err := os.OpenFile(
		tempPath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0600,
	)
	if err != nil {
		return err
	}

	// Write current state to temporary log
	for key, value := range s.data {

		entry := logEntry{
			Op:    "SET",
			Key:   key,
			Value: value,
		}

		data, err := json.Marshal(entry)
		if err != nil {
			tempFile.Close()
			return err
		}

		data = append(data, '\n')

		byteCount, err := tempFile.Write(data)
		if err != nil {
			tempFile.Close()
			return err
		}

		if byteCount != len(data) {
			tempFile.Close()
			return fmt.Errorf("short write")
		}
	}

	// Make sure the temporary log is durable
	if err := tempFile.Sync(); err != nil {
		tempFile.Close()
		return err
	}

	// Close temp file before replacing old log
	if err := tempFile.Close(); err != nil {
		return err
	}

	// Close the old log
	if err := s.logFile.Close(); err != nil {
		return err
	}

	// Replace old log with compacted log
	if err := os.Rename(tempPath, logPath); err != nil {
		// Try to restore a usable file handle
		file, reopenErr := os.OpenFile(
			logPath,
			os.O_CREATE|os.O_APPEND|os.O_RDWR,
			0600,
		)

		if reopenErr == nil {
			s.logFile = file
		}

		return err
	}

	// Open the new compacted log
	file, err := os.OpenFile(
		logPath,
		os.O_CREATE|os.O_APPEND|os.O_RDWR,
		0600,
	)
	if err != nil {
		return err
	}

	s.logFile = file

	return nil
}
