package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeJsnResponse(w http.ResponseWriter, v any) {
	out, err := json.Marshal(v)
	if err != nil {
		log.Println("handler error at writing json:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
