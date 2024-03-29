package main

import (
	"flag"
	"log"
	"net/http"

	"umigame-api/models"
	"umigame-api/routers"
	"umigame-api/tasks"
	"umigame-api/utils"

	"github.com/kelseyhightower/envconfig"
)
func enableCors(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // ここにCORS設定を記述
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token")

        // OPTIONSメソッドに対するプリフライトリクエストへの対応
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func main() {
	port := flag.String("p", ":8080", "HTTP network port")
	flag.Parse()

	var env models.Env
	envconfig.Process("", &env)

	db, err := utils.ConnectDB(env)
	if err != nil {
		log.Println(err)
		return
	}

	r := routers.NewRouter(db, *port)
	wrappedHandler := enableCors(r)
    http.Handle("/", wrappedHandler)
	go tasks.ExeTasks(db)

	log.Printf("server start at port %s", *port)
	log.Fatal(http.ListenAndServe(*port, wrappedHandler))
}
