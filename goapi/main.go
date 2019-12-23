package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func revert(txt string) string {
	words := strings.Split(txt, " ")

	newWords := make([]string, len(words))

	for i := 0; i < len(words); i++ {
		newWords[i] = words[len(words)-1-i]
	}

	return strings.Join(newWords, " ")
}

type reverseRequest struct {
	MessageContent string `json:"content"`
}

type reverseResponse struct {
	ResponseResult string `json:"result"`
}

func main() {
	log.Print("GOAPI TEST, version ", VERSION, " (", BUILDDATE, ")")

	server := http.NewServeMux()

	// simple api: getting {"content": "some text"}, it returns {"result": "texts some"}
	server.HandleFunc("/revert", func(w http.ResponseWriter, r *http.Request) {
		log.Println("revert invoked - ", r.Method)
		log.Printf("headers -> %v", r.Header)

		if r.Method == http.MethodPost {
			msg := reverseRequest{}

			if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
				log.Println("bad request, error: ", err)
				w.WriteHeader(http.StatusBadRequest)
			} else if msg.MessageContent != "" {
				rsp := reverseResponse{ResponseResult: revert(msg.MessageContent)}

				log.Println("reverting:", msg.MessageContent, "->", rsp.ResponseResult)

				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(rsp)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}

		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	// listining
	port := getEnv("PORT", "8001")
	log.Print("Running server at :", port)

	log.Fatal(http.ListenAndServe(":"+port, server))
}
