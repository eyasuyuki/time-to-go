package main

import (
  "github.com/gorilla/mux"
  "log"
  "net/http"
  "encoding/json"
  "errors"
  "strconv"
)

type Artist struct {
  Id    int     `json:id`
  Name  string  `json:name`
  Part  string  `json:part`
}

type Artists struct {
  Items  map[int]*Artist
}

var count int
var artists *Artists

func (a *Artists)Create(name string, part string) (*Artist) {
  count++
  item := new(Artist)
  item.Id = count
  item.Name = name
  item.Part = part
  a.Items[item.Id] = item
  return item
}

func (a *Artists)Read(id int) (*Artist, error) {
  item := a.Items[id]
  if item == nil {
    return nil, errors.New("not found")
  }
  return item, nil
}

func (a *Artists)Update(id int, name string, part string) (*Artist, error) {
  item, err := a.Read(id)
  if err != nil {
    return nil, err
  }
  item.Name = name
  item.Part = part
  return item, nil
}

func (a *Artists)Delete(id int) (error) {
  _, err := a.Read(id)
  if err != nil {
    return err
  }
  delete(a.Items, id)
  return nil
}

func (a *Artists)List() ([]*Artist) {
  list := make([]*Artist, 0, len(a.Items))
  for _, item := range a.Items {
    list = append(list, item)
  }
  return list
}

func init() {
  count = 0
  artists = new(Artists)
}

// /artist?name=<var>&part=<var>
func post(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  name := params["name"]
  part := params["part"]
  log.Printf("post: name=%s, part=%s\n", name, part)
  item := artists.Create(name, part)
  js, _ := json.Marshal(item)
  log.Printf("post: json=%s\n", js)
  w.Write(js)
}

// /artist/id
func get(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  id := params["id"]
  log.Printf("get: id=%s\n", id)
  i, err := strconv.Atoi(id)
  if err != nil {
    // bad argument
    log.Printf("get: error: %s\n", err.Error())
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  item, err := artists.Read(i)
  if err != nil {
    log.Printf("get: error: %s\n", err.Error())
    http.Error(w, err.Error(), http.StatusNotFound)
    return
  }
  js, _ := json.Marshal(item)
  log.Printf("get: json=%s\n", js)
  w.WriteHeader(http.StatusCreated)
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

// /artist/id&name=<var>&part=<var>
func put(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  id := params["id"]
  i, err := strconv.Atoi(id)
  if err != nil {
    // bad argument
    log.Printf("put: error: %s\n", err.Error())
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  name := params["name"]
  part := params["part"]
  item, err := artists.Update(i, name, part)
  if err != nil {
    log.Printf("put: error: %s\n", err.Error())
    http.Error(w, err.Error(), http.StatusNotFound)
    return
  }
  js, _ := json.Marshal(item)
  log.Printf("put: json=%s\n", js)
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

// artist/id
func del(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  id := params["id"]
  i, err := strconv.Atoi(id)
  if err != nil {
    // bad argument
    log.Printf("get: error: %s\n", err.Error())
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  err = artists.Delete(i)
  if err != nil {
    log.Printf("get: error: %s\n", err.Error())
    http.Error(w, err.Error(), http.StatusNotFound)
    return
  }
  // response 204
  w.WriteHeader(http.StatusNoContent)
}

func list(w http.ResponseWriter, r *http.Request) {
  js, _ := json.Marshal(artists.List())
  log.Printf("list: json=%s\n", js)
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

func main() {
  router := mux.NewRouter()
  router.Queries("name", "", "part", "")
  router.HandleFunc("/artist/{id:[0-9]+}", get).Methods("GET")
  router.HandleFunc("/artist", post).Methods("POST")
  router.HandleFunc("/artist/{id:[0-9]+}", put).Methods("PUT")
  router.HandleFunc("/artist/{id:[0-9]+}", del).Methods("DELETE")
  router.HandleFunc("/artist/list", list).Methods("GET")

  http.Handle("/", router)

  log.Println("Listening...")
  http.ListenAndServe(":3000", nil)
}
