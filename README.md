# Time to GO

タイムインターメディア技術部会資料

# Website

http://www.golang.org/

# Tutorial

http://go-tour-jp.appspot.com/

# 文法

- 割愛します。（C/C++, Java経験者なら楽勝)
- オブジェクト指向ではない
 - 構造体
 - メソッド
 - インターフェース

# 良い点

- Makefileがいらない
 - GOPATHとgo get (デモ)
- 速い→コンパイラ言語
- 並列処理が楽に書ける
 - go routine
 - channel
- JSONと構造体の相互変換が楽 (デモ)
- 最初からテンプレートエンジンを持っている (デモ)
- ネットワーク関連ライブラリが充実
 - ex. CGIやサーブレットではなくhttpdを直接書ける
- ガベージコレクションがある
- ポータビリティ→クロスコンパイルも楽

# 悪い点

- 検索しにくい→Golangで検索
- バイナリのサイズが大きい (Cなどに比べて)
- アセンブラやC/C++より遅い
- interfaceがややこしい

# イベント

- GoCon https://github.com/GoCon/GoCon
- 電車でGo! http://gocon.connpass.com/event/2108/
- ヒカルのGo! http://connpass.com/event/12179/
- etc.

# キラーアプリ

- Docker http://docker.com/
- Revel http://revel.github.io/

# 結論

- サーバーサイドなら最適な選択
 - ex. Google App Engine (デモ)

# デモ

## インストール

    brew install go

## HTTPサーバー

### ソースコード

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

### gorilla muxライブラリ

https://github.com/gorilla/mux

### GOPATHとgo get

    export GOPATH=~/gocode:~/git/time-to-go/src/demo01
    go get github.com/gorilla/mux

### 実行とビルド

    go run ~/git/time-to-go/src/demo01/httpd.go
    go build ~/git/time-to-go/src/demo01/httpd.go

### アクセス

    curl -i http://127.0.0.1:3000/user/yasuyuki/profile

## JSONと構造体の相互変換

### 構造体の定義

    type Artist struct {
    	Id   int    `json:"id"`
    	Name string `json:"name"`
    	Part string `json:"part"`
    }

### JSONデータ

    var data = `
    [{"id":0, "name":"Jhon", "part":"Guitar"}
    ,{"id":0, "name":"Paul", "part":"Bass"}
    ,{"id":0, "name":"George", "part":"Guitar"}
    ,{"id":0, "name":"Ringo", "part":"Drums"}]
    `

### encoding/jsonライブラリ

    import (
    	"encoding/json"
    	// (略)
    }

### JSONのパーズ

    var items []Artist
    err := json.Unmarshal([]byte(data), &items)

### 構造体からJSONを出力

    func (a *Artist) Marshal() ([]byte, error) {
    	js, err := json.Marshal(a)
    	// (略)
    	return js, err
    }

### RESTサーバー

    func main() {
    	router := mux.NewRouter()
    	router.HandleFunc("/artist/{id:[0-9]+}", get).Methods("GET")
    	router.HandleFunc("/artist", post).Methods("POST")
    	router.HandleFunc("/artist/{id:[0-9]+}", put).Methods("PUT")
    	router.HandleFunc("/artist/{id:[0-9]+}", del).Methods("DELETE")
    	router.HandleFunc("/artist/list", list).Methods("GET")

    	http.Handle("/", router)

    	log.Println("Listening...")
    	http.ListenAndServe(":3000", nil)
    }

### 起動

    export GOPATH=~/gocode:~/git/time-to-go/src/demo02
    go run ~/git/time-to-go/src/demo02/band.go

### 接続テスト

    curl -i http://127.0.0.1:3000/artist/list
    curl -i http://127.0.0.1:3000/artist/0
    curl -i -X POST http://127.0.0.1:3000/artist -d "name=Jake" -d "part=Ukulele"
    curl -i -X PUT http://127.0.0.1:3000/artist/0 -d "name=Brian"
    curl -i -X DELETE http://127.0.0.1:3000/artist/4

## テンプレートエンジン

### template/htmlライブラリ

    import (
    	// (略)
    	"html/template"
    	// (略)
    )

### 構造体の定義

    type Customer struct {
    	Email string `json:"email"`
    	Name  string `json:"name"`
    	Item  string `json:"item"`
    }

### HTMLテンプレート

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

### データと初期化

#### JSONデータ

    var data = []byte(`[
      {"email":"suzuki@example.com",  "name":"鈴木", "item":"iPad mini"}
      ,{"email":"sato@example.com",   "name":"佐藤", "item":"Xperia Z"}
      ,{"email":"tanaka@example.com", "name":"田中", "item":"Surface Pro"}
    ]`)

#### 初期化

    var customers []Customer

    func init() {
    	err := json.Unmarshal(data, &customers)
    	if err != nil {
    		log.Fatalf("init: error=%S", err.Error())
    	}
    }

### テンプレートのパーズ

    func handler(w http.ResponseWriter, r *http.Request) {
    	t := template.Must(template.New("page").Parse(page))

### テンプレートの適用

    	w.Header().Set("Content-Type", "text/html")
    	t.Execute(w, &customers)
    }


## Google App Engine

### インストール

    brew install go-app-engine-64

### app.yaml

    application: hinerutojar-hrd
    version: 1
    runtime: go
    api_version: go1

    handlers:
    - url: /.*
      script: _go_app

### ソースコード(一部)

    var router = mux.NewRouter()

    func init() {
    	// init data
    	var items []Artist
    	err := json.Unmarshal(data, &items)
    	if err != nil {
    		// log only localhost
    		log.Printf("init: %s", err.Error())
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

### ローカル環境での実行

    export GOPATH=~/gocode:~/git/time-to-go/src/demo04
    cd time-to-go/src/demo04
    goapp serve

### 接続テスト

    curl -i http://127.0.0.1:8080/artist/list
    curl -i http://127.0.0.1:8080/artist/0
    curl -i -X POST http://127.0.0.1:8080/artist -d "name=Jake" -d "part=Ukulele"
    curl -i -X PUT http://127.0.0.1:8080/artist/0 -d "name=Brian"
    curl -i -X DELETE http://127.0.0.1:8080/artist/4

### Google App Engineへのデプロイ

    goapp deploy


### 接続テスト

    curl -i http://hinerutojar-hrd.appspot.com/artist/list
    curl -i http://hinerutojar-hrd.appspot.com/artist/0
    curl -i -X POST http://hinerutojar-hrd.appspot.com/artist -d "name=Jake" -d "part=Ukulele"
    curl -i -X PUT http://hinerutojar-hrd.appspot.com/artist/0 -d "name=Brian"
    curl -i -X DELETE http://hinerutojar-hrd.appspot.com/artist/4
