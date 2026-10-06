package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
    "syscall"
	"github.com/Anshit-Gupta/kairo/internal/server"
	"github.com/Anshit-Gupta/kairo/internal/store"
)


func main() {

   sigChan := make(chan os.Signal,1)
   //Create store

   dataPath := os.Getenv("KAIRO_DATA_PATH")
   if dataPath == "" {
    dataPath = "./data/kairo.log"
   }

	s, err := store.NewStore("./data/kairo.log")
	if err != nil {
		panic(err)
	}
	defer s.Close()
  //create Server(store)
    srv:=server.NewServer(s)
   
   //start https server 
   

   port := os.Getenv("PORT")

   if port==""{
      port = "8080"
   }

    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

     httpServer := &http.Server{
      Addr: ":" + port,
      Handler: srv,
     }

    go func() {
    if err := httpServer.ListenAndServe(); err != nil &&
        err != http.ErrServerClosed {
        fmt.Println("Server error:", err)
    }
   }()

     sig := <-sigChan //wait untill someone send a value into sigChan
    fmt.Println("Received:", sig)
     //after getting the signal , we tell http server too shut donw gracefully 
     
     //The shutdown process has at most 5 seconds
    ctx,cancel := context.WithTimeout(context.Background() , 5*time.Second) 
    defer cancel()

     if err := httpServer.Shutdown(ctx);err!=nil{
      fmt.Println("Server shutdown error: ", err)
     }
     



     
     


}
