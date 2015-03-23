// Code from: http://www.alexedwards.net/blog/a-mux-showdown

package main

import (
  "github.com/gorilla/mux"
  "log"
  "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  name := params["name"]
  w.Write([]byte("Hello " + name))
}

func main() {
  router := mux.NewRouter()
  router.HandleFunc("/user/{name:[a-z]+}/profile", handler).Methods("GET")

  http.Handle("/", router)

  log.Println("Listening...")
  http.ListenAndServe(":3000", nil)
}
