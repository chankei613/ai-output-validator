// cmd/ovserve はAI Output Validator APIをlocalhostで提供する単体サーバー。
//
//	go run ./cmd/ovserve -addr :8426 -db ai-output-validator.db
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/chankei613/ai-output-validator/internal/api"
	"github.com/chankei613/ai-output-validator/internal/db"
)

func main() {
	addr := flag.String("addr", ":8426", "待ち受けアドレス")
	dbPath := flag.String("db", "ai-output-validator.db", "SQLiteファイル")
	flag.Parse()

	conn, err := db.Init(*dbPath)
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	router := api.NewRouter(conn)
	log.Printf("ai-output-validator backend listening on %s", *addr)
	if err := http.ListenAndServe(*addr, router); err != nil {
		log.Fatal(err)
	}
}
