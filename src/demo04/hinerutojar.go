package hinerutojar

import (
	"net/http"
	"github.com/gorilla/mux"
	"appengine"
  "encoding/json"
  "errors"
  "strconv"
	"log"
)

type Artist struct {
  Id    int     `json:"id"`
  Name  string  `json:"name"`
  Part  string  `json:"part"`
}

func (a *Artist)Marshal(c appengine.Context) ([]byte, error) {
  js, err := json.Marshal(a)
  if err != nil {
    c.Errorf("err=%s", err.Error())
    return js, err
  }
  // log
  c.Debugf("js=%s", js)

  return js, err
}

type Artists struct {
  Items  map[int]*Artist
}

var count int
var artists *Artists

func (a *Artists)Create(name string, part string) (*Artist) {
  count++
  item := &Artist{count, name, part}
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
  if name != "" {
    item.Name = name
  }
  if part != "" {
    item.Part = part
  }
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

func (a *Artists)MarshalIndent(c appengine.Context) ([]byte, error) {
  list := a.List()
  js, err := json.MarshalIndent(list, "", "  ")
  if err != nil {
    c.Errorf("err=%s", err.Error())
    return js, err
  }
  logtext, _ := json.Marshal(list)
  // log
  c.Debugf("js=%s", logtext)

  return js, err
}

var data = []byte(`[
{"id":0, "name":"Jhon", "part":"Guitar"}
,{"id":0, "name":"Paul", "part":"Bass"}
,{"id":0, "name":"George", "part":"Guitar"}
,{"id":0, "name":"Ringo", "part":"Drums"}
]`)

var router = mux.NewRouter()

func init() {
	// init data
	var items []Artist
  err := json.Unmarshal(data, &items)
  if err != nil {
		// log only localhost
    log.Printf("init: %s", err.Error())
//    os.Exit(1)
  }
  artists = new(Artists)
  artists.Items = make(map[int]*Artist)
  count = -1
  for i, _ := range items {
		// log only localhost
		log.Printf("init: i=%d, name=%s, part=%s", i, items[i].Name, items[i].Part)
    count++
    items[i].Id = i
    artists.Items[i] = &items[i]
  }

	// handlers
	router.HandleFunc("/artist/{id:[0-9]+}", get).Methods("GET")
  router.HandleFunc("/artist", post).Methods("POST")
  router.HandleFunc("/artist/{id:[0-9]+}", put).Methods("PUT")
  router.HandleFunc("/artist/{id:[0-9]+}", del).Methods("DELETE")
  router.HandleFunc("/artist/list", list).Methods("GET")
	http.Handle("/", router)
}


// /artist?name=<var>&part=<var>
func post(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
  name := r.FormValue("name")
  part := r.FormValue("part")
  c.Debugf("post: name=%s, part=%s", name, part)
  item := artists.Create(name, part)
  js, _ := item.Marshal(c)
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

// /artist/id
func get(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
  params := mux.Vars(r)
  id := params["id"]
  c.Debugf("get: id=%s", id)
  i, err := strconv.Atoi(id)
  if err != nil {
    // bad argument
    c.Errorf("get: error: %s", err.Error())
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  item, err := artists.Read(i)
  if err != nil {
    c.Errorf("get: error: %s", err.Error())
    http.Error(w, err.Error(), http.StatusNotFound)
    return
  }
  js, _ := item.Marshal(c)
  w.WriteHeader(http.StatusCreated)
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

// /artist/id&name=<var>&part=<var>
func put(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
  params := mux.Vars(r)
  id := params["id"]
  i, err := strconv.Atoi(id)
  if err != nil {
    // bad argument
    c.Errorf("put: error: %s", err.Error())
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  name := r.FormValue("name")
  part := r.FormValue("part")
  item, err := artists.Update(i, name, part)
  if err != nil {
    c.Errorf("put: error: %s", err.Error())
    http.Error(w, err.Error(), http.StatusNotFound)
    return
  }
  js, _ := item.Marshal(c)
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

// artist/id
func del(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
  params := mux.Vars(r)
  id := params["id"]
  i, err := strconv.Atoi(id)
  if err != nil {
    // bad argument
    c.Errorf("get: error: %s", err.Error())
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  err = artists.Delete(i)
  if err != nil {
    c.Errorf("get: error: %s", err.Error())
    http.Error(w, err.Error(), http.StatusNotFound)
    return
  }
  // response 204
  w.WriteHeader(http.StatusNoContent)
}

func list(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
  js, _ := artists.MarshalIndent(c)
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}
