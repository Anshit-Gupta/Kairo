package main

import (
	"fmt"
   "net/http"
	"github.com/Anshit-Gupta/kairo/internal/server"
	"github.com/Anshit-Gupta/kairo/internal/store"
)



func main() {
   //Create store
	s, err := store.NewStore("./data/kairo.log")
	if err != nil {
		panic(err)
	}
	defer s.Close()
  //create Server(store)
    srv:=server.NewServer(s)
   
   //start https server 
   err=http.ListenAndServe(":8080",srv)
   fmt.Println("Server started succesfully")


   fmt.Println("satrting the main code ")
	if err := s.Set("name", "Ansh"); err != nil {
    panic(err)
   }

      if err := s.Set("language", "Go"); err != nil {
    panic(err)
   }

   if err := s.Close(); err != nil {
    panic(err)
   }
}
