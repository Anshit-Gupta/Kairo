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
	if err := os.MkdirAll(filepath.Dir(path),0755); err!=nil{
		return nil,err
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
	if err!=nil{
		fmt.Println("Error in Sync")
		return fmt.Errorf("Error in Sync : %w",err)
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
	if err!=nil{
		fmt.Println("Error in Sync")
		return fmt.Errorf("Error in Sync : %w",err)
	}

   
	delete(s.data, key) 
	return nil
}

//since we are keeping the file open while the program is running , it's better to ahve close fun

func (s *Store) Close() error {
	return s.logFile.Close()
}


//operations for logreplay 
func (s* Store)logReplayOperationHelper(entry logEntry)error{
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