package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var pageTmpl = template.Must(template.New("page").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Famous Quotes</title>
  <style>
    body { font-family: sans-serif; margin: 2rem; }
    h1    { margin-bottom: 1rem; }
    ul    { list-style: disc inside; }
  </style>
</head>
<body>
  <h1>Famous Quotes</h1>
  <ul>
    {{range .}}
      <li>{{.}}</li>
    {{end}}
  </ul>
</body>
</html>
`))

func main() {
	// Read env-vars (injected by your Deployment)
	host := getenv("QUOTES_HOSTNAME", "localhost")
	dbName := getenv("QUOTES_DATABASE", "quotes")
	user := getenv("QUOTES_USER", "myuser")
	pass := os.Getenv("QUOTES_PASSWORD") // empty if missing

	// Force TCP to avoid socket fallback
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=true",
		user, pass, host, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT message FROM quotes")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var quotes []string
		for rows.Next() {
			var msg string
			if scanErr := rows.Scan(&msg); scanErr == nil {
				quotes = append(quotes, msg)
			}
		}
		if err = rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pageTmpl.Execute(w, quotes)
	})

	log.Println("Serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
