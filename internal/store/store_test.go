package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRecoveryWithIncompleteMiddleEntry(t *testing.T) {
	logContent := `{"Op":"SET","Key":"name","Value":"Anshit"}
    {"Op":"SET","Key":"age","Value":
    {"Op":"SET","Key":"city","Value":"Bangalore"}`
	dir := t.TempDir()
	logPath := filepath.Join(dir, "TestRecoveryWithIncompleteFinalEntry.log")
	err := os.WriteFile(logPath, []byte(logContent), 0600)
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewStore(logPath)

	if err == nil {
		t.Fatalf("expected recovery to fail")
	}

}

func TestRecoveryWithIncompleteFinalEntry(t *testing.T) {
	//here we will manully add 2 complete + 1 currupted log
	logContent := `{"Op":"SET","Key":"name","Value":"Anshit"}
    {"Op":"SET","Key":"age","Value":"21"}
    {"Op":"SET","Key":"city","Value":`

	dir := t.TempDir()
	logPath := filepath.Join(dir, "TestRecoveryWithIncompleteFinalEntry.log")
	err := os.WriteFile(logPath, []byte(logContent), 0600)
	if err != nil {
		t.Fatal(err)
	}

	s, err := NewStore(logPath)
	if err != nil {
		t.Fatalf("Error creating new store : %v", err)
	}

	defer s.Close()
	//now we check if name and age are read into the map or not

	value, exists := s.Get("name")
	if !exists {
		t.Fatalf("Expected key to exist")
	}

	if value != "Anshit" {
		t.Fatalf("Expected value : %s got : %s", "Anshit", value)
	}

	value, exists = s.Get("age")
	if !exists {
		t.Fatalf("Expected key to exist")
	}

	if value != "21" {
		t.Fatalf("Expected value : %s got : %s", "21", value)
	}
}

func TestPersistence(t *testing.T) {
	//testing if the data surive when store is closed

	dir := t.TempDir()
	path := filepath.Join(dir, "TestPersistece.log")

	s, err := NewStore(path)
	if err != nil {
		t.Fatalf("Error creating new store : %v", err)
	}

	//set key
	if err := s.Set("key", "value"); err != nil {
		t.Fatalf("Failed to set key %v ", err)
	}
	//close the store
	if err := s.Close(); err != nil {
		t.Fatalf("Failed to close store: %v", err)
	}
	//now repone the same path and check if the data still persist
	s, err = NewStore(path)
	if err != nil {
		t.Fatalf("Error Reopening the same store : %v", err)
	}
	defer s.Close()
	//now we check it by using get
	value, exists := s.Get("key")
	if !exists {
		t.Fatalf("Expected key to exist")
	}

	if value != "value" {
		t.Fatalf("Expected value : %s got : %s", "value", value)
	}

}

func TestPersistenceAfterUpdatesAndDeletes(t *testing.T) {

	dir := t.TempDir()
	path := filepath.Join(dir, "TestPersistenceAfterUpdatesAndDeletes.log")

	s, err := NewStore(path)
	if err != nil {
		t.Fatalf("Error creating new store : %v", err)
	}

	//set key
	if err := s.Set("key1", "value1"); err != nil {
		t.Fatalf("Failed to set key %v ", err)
	}

	if err := s.Set("key2", "value2"); err != nil {
		t.Fatalf("Failed to set key %v ", err)
	}

	//update
	if err := s.Set("key1", "Newvalue"); err != nil {
		t.Fatalf("Failed to set key %v ", err)
	}
	//delete
	err = s.Delete("key2")
	if err != nil {
		t.Fatalf("failed to delete key: %v", err)
	}

	//close the store
	if err := s.Close(); err != nil {
		t.Fatalf("Failed to close store: %v", err)
	}
	//Reopen it

	s, err = NewStore(path)
	if err != nil {
		t.Fatalf("Error Reopening the same store : %v", err)
	}
	defer s.Close()
	//check if the value is updated
	value, exists := s.Get("key1")
	if !exists {
		t.Fatalf("Expected key to exist")
	}

	if value != "Newvalue" {
		t.Fatalf("Expected value : %s got : %s", "Newvalue", value)
	}
	// check deleted

	value, exists = s.Get("key2")
	if exists {
		t.Fatalf("Expected key to Not exist")
	}

	if value != "" {
		t.Fatalf("Expected value to be empty string")
	}

}

func TestSetandGet(t *testing.T) {
	//we test set and get

	//to store log we will use temporary dir which will get removed automatically
	dir := t.TempDir()
	path := filepath.Join(dir, "kairo.log")

	//call NewStore with path
	s, err := NewStore(path)
	if err != nil {
		t.Fatalf("Error creating new store : %v", err)
	}
	//dont forget to close the file
	defer s.Close()

	//send set
	if err = s.Set("name", "Anshit"); err != nil {
		t.Fatalf("Failed to set key %v ", err)
	}
	//test get and all
	value, exists := s.Get("name")
	if !exists {
		t.Fatalf("Expected key to exist")
	}
	if value != "Anshit" {
		t.Fatalf("Expected value : %s got : %s", "Anshit", value)
	}

}

func TestUpdateExistingKey(t *testing.T) {

	dir := t.TempDir()
	path := filepath.Join(dir, "UpdateExistingKey.log")

	s, err := NewStore(path)
	if err != nil {
		t.Fatalf("Error creating new store : %v", err)
	}

	defer s.Close()

	//send set
	if err = s.Set("name", "Anshit"); err != nil {
		t.Fatalf("Failed to set key %v ", err)
	}

	//update exsting key
	if err = s.Set("name", "NewName"); err != nil {
		t.Fatalf("Expected to update key")
	}

	value, exists := s.Get("name")
	if !exists {
		t.Fatalf("Expected key to exist")
	}
	if value != "NewName" {
		t.Fatalf("Expected value : %s got : %s", "NewName", value)
	}

}
func TestGetMissingKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "MissingKey.log")

	s, err := NewStore(path)
	if err != nil {
		t.Fatalf("Error creating new store : %v", err)
	}
	defer s.Close()
	//without setting we call get
	value, exists := s.Get("key")
	if exists {
		t.Fatalf("NO Existing key found %s", value)
	}
	if value != "" {
		t.Fatalf("Expected to be empty string got %s", value)
	}
}

func TestDeleteNonExistingKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "DeleteNonExistingKey.log")

	s, err := NewStore(path)
	if err != nil {
		t.Fatalf("Error creating a new store : %v", err)
	}
	defer s.Close()

	//deleting a non-exsting key
	if err = s.Delete("key"); err != nil {
		t.Fatalf("Delete returned unexpected error: %v", err)
	}
}

func TestDeleteExistingKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "DeleteExistingKey.log")

	s, err := NewStore(path)
	if err != nil {
		t.Fatalf("Error creating a new store")
	}
	defer s.Close()
	//send set
	if err = s.Set("key", "value"); err != nil {
		t.Fatalf("Failed to set key %v ", err)
	}
	//deleting a exsting key
	if err = s.Delete("key"); err != nil {
		t.Fatalf("Error Deleting exsiting key")
	}

	//get on the delted key to confirm
	_, exists := s.Get("key")
	if exists {
		t.Fatalf("Expected to be deleted  ")
	}
}
