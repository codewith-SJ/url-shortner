
package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
	"net/http"
)

func generateCode(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func main() {
	// Initialize DB connection
	InitDB()

	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			URL string `json:"url"`
		}

		var req Request
		json.NewDecoder(r.Body).Decode(&req)

		code := generateCode(6) // temporary

		_, err := DB.Exec(
			"INSERT INTO short_urls (code, original_url) VALUES ($1, $2)",
			code, req.URL,
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Write([]byte("Short URL created: " + code))
	})

	http.HandleFunc("/r/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Path[len("/r/"):]

		var original string
		err := DB.QueryRow(
			"SELECT original_url FROM short_urls WHERE code=$1",
			code,
		).Scan(&original)

		if err != nil {
			http.Error(w, "Not found", 404)
			return
		}

		http.Redirect(w, r, original, http.StatusFound)
	})

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}

