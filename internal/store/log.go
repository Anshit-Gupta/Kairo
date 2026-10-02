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

	//insterad of using a for loop , we are doing each log manually
	if !scanner.Scan() {
		return nil
	}

	currentLine := scanner.Text()

	for {

		if scanner.Scan() { //returns a bool , we check if the next line  exist
			var entry logEntry
			//convert the json into go using json.Unmarshal
			err := json.Unmarshal([]byte(currentLine), &entry) //daya,desitination
			if err != nil {
				fmt.Println("error in json->go")
				return err
			}

			err = s.logReplayOperationHelper(entry)
			if err != nil {
				return fmt.Errorf("error doing SET/DELETE operation: %w", err)
			}
			currentLine = scanner.Text() //move the current to next

		} else {
			//here we know that the currentLine is the last one
			var entry logEntry
			if err := json.Unmarshal([]byte(currentLine), &entry); err != nil {
				break
			}

			err = s.logReplayOperationHelper(entry)
			if err != nil {
				return fmt.Errorf("error doing SET/DELETE operation: %w", err)
			}
			break

		}

	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
