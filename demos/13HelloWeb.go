package demos

import (
	"fmt"
	"net/http"
)

func handle(writer http.ResponseWriter, reader *http.Request) {
	_, err := fmt.Fprintf(writer, "Hello Web")
	if err != nil {
		return
	}
}

func Web() {
	server := &http.Server{Addr: "127.0.0.1:8080"}
	http.HandleFunc("/", handle)
	err := server.ListenAndServe()
	if err != nil {
		return
	}
	fmt.Println("WebServer started !")
}
