# Time to GO

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
- 並列処理が楽に書ける→go routine
- JSONと構造体の相互変換が楽 (デモ)
- 最初からテンプレートエンジンを持っている
- ネットワーク関連ライブラリが充実
 - CGIやサーブレットではなくhttpdを直接書ける
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
 - ex. Google App Engine

# デモ

## インストール

    brew install go

## HTTPサーバー

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

## GOPATHとgo get

    mkdir gocode
    cd gocode
    export GOPATH=`pwd`:~/git/time-to-go/src/demo01
	go get github.com/gorilla/mux
	go run ~/git/time-to-go/src/demo01/httpd.go
	go build ~/git/time-to-go/src/demo01/httpd.go

## JSONと構造体の相互変換
## テンプレートエンジン
## Google App Engine
