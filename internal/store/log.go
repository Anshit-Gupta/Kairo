package store

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type logEntry struct {
	Op    string
	Key   string
	Value string
}

func (s *Store) logReplay(path string) error {
	//open file in read
	File, err := os.Open(path)

	if err != nil {
		fmt.Println("error in opening file")
		return err
	}
	//gotta close the file as well
	defer File.Close()

	//read line by line and update the map
	scanner := bufio.NewScanner(File)

	for scanner.Scan() {
		//read one line
		line := scanner.Text()
		var entry logEntry
		//convert the json into go using json.Unmarshal
		err := json.Unmarshal([]byte(line), &entry) //daya,desitination
		if err != nil {
			fmt.Println("error in json->go")
			return err
		}

		if entry.Op == "SET" {
			s.data[entry.Key] = entry.Value
		} else if entry.Op == "DELETE" {
			//Delete
			delete(s.data, entry.Key)
		} else {
			//unknown operation or typo
			return fmt.Errorf("Unknown Operation detected %q", entry.Op)
		}

	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
