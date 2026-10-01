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
