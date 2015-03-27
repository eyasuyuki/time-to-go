package main

import (
  "github.com/gorilla/mux"
  "log"
  "net/http"
  "encoding/json"
  "html/template"
)

type Customer struct {
  Email string  `json:"email"`
  Name  string  `json:"name"`
  Item  string  `json:"item"`
}

var data = []byte(`[
  {"email":"suzuki@example.com",  "name":"鈴木", "item":"iPad mini"}
  ,{"email":"sato@example.com",   "name":"佐藤", "item":"Xperia Z"}
  ,{"email":"tanaka@example.com", "name":"田中", "item":"Surface Pro"}
]`)

var page = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>顧客リスト</title>
  </head>
  <body>
    <table>
      <thead>
        <tr><th>email</th><th>名前</th><th>購入品目</th></tr>
      </thead>
      <tbody>
      {{range .}}
        <tr><td>{{.Email}}</td><td>{{.Name}}</td><td>{{.Item}}</td></tr>
      {{end}}
      </tbody>
    </table>
  </body>
</html>
`

var customers []Customer

func init() {
  err := json.Unmarshal(data, &customers)
  if err != nil {
    log.Fatalf("init: error=%S", err.Error())
  }
}

func handler(w http.ResponseWriter, r *http.Request) {
  t := template.Must(template.New("page").Parse(page))

  w.Header().Set("Content-Type", "text/html")
  t.Execute(w, &customers)
}

func main() {
  router := mux.NewRouter()
  router.HandleFunc("/customers", handler).Methods("GET")

  http.Handle("/", router)

  log.Println("Listening...")
  http.ListenAndServe(":3000", nil)
}
