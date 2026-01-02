package http

import (
	"net/http"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Test Handler Response"))
}
func HelloHandler() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	http.ListenAndServe(":9999", nil)

}
