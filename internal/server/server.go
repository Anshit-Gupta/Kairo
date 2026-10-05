package server

import (
	"fmt"
	"github.com/Anshit-Gupta/kairo/internal/store"
	"io"
	"net/http"
	"strings"
)

type server struct {
	store *store.Store
}

// consttructor
// The server owns a reference to Kairo, but doesn't implement storage itself.
func NewServer(s *store.Store) *server {
	return &server{
		store: s,
	}
}

//HTTP handler

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {



   //compact
	if r.URL.Path == "/compact"{
		if r.Method != "POST" {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
      }

		err := s.store.Compact()

		if err!=nil{
			http.Error(w,"compaction failed",http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w,"compaction succesfull")
		return
	}


	
	key := strings.TrimPrefix(r.URL.Path,"/")
	//edge case : empty key 
	if key == "" {
    http.Error(w, "key cannot be empty", http.StatusBadRequest)
    return
}

	if r.Method == "GET" {

		value, exists := s.store.Get(key)

		if !exists {
			http.Error(w, "key not found", http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, value)
		return
	}

	// update
	if r.Method == "PUT" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		value := string(body) //since we get body in bytes by io.readall , we convert it into string since our func expects string
		err = s.store.Set(key, value)
		if err != nil {
			http.Error(w, "failed to update", http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "key updated successfully")
		return

	}

	//delete
	if r.Method == "DELETE" {

		err := s.store.Delete(key)
		if err != nil {
			http.Error(w, "failed to delete key", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
        return

	}
   
	//since go's default is 200 , for unsupported methods we will return error 
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	

}
