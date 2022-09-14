package main

import (
	"net/http"
	"os"
  "proxer/client"
	"github.com/gorilla/handlers"
  "fmt"
)

func main() {
  fmt.Println("- Starting PROXER -")
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, client.New()))
	http.ListenAndServe(":8080", nil)
  fmt.Println("- Proxer STARTED on port 8080 -")
}
