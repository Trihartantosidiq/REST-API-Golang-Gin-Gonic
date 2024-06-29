package helpers

import (
	"fmt"
	"log"
	"net/http"
)

func LogError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"error": "%s"}`, message)
}
