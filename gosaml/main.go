package main

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/crewjam/saml/samlsp"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func revertString(txt string) string {
	return "this is the default implementation"
}

type reverseRequest struct {
	MessageContent string `json:"content"`
}

type reverseResponse struct {
	ResponseResult string `json:"result"`
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", samlsp.Token(r.Context()).Attributes.Get("givenName"))
}

// simple api: getting {"content": "some text"}, it returns {"result": "texts some"}
func revert(w http.ResponseWriter, r *http.Request) {
	log.Println("revert invoked - ", r.Method)
	log.Printf("headers -> %v", r.Header)

	if r.Method == http.MethodPost {
		msg := reverseRequest{}

		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			log.Println("bad request, error: ", err)
			w.WriteHeader(http.StatusBadRequest)
		} else if msg.MessageContent != "" {
			rsp := reverseResponse{ResponseResult: revertString(msg.MessageContent)}

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
}

func main() {
	log.Print("GOSAML TEST, version ", VERSION, " (", BUILDDATE, ")")

	certificate := getEnv("CERTFILE", "myservice.cert")
	privateKey := getEnv("KEYFILE", "myservice.key")
	idpURL := getEnv("IDP_URL", "http://localhost:8080/auth/realms/demo/protocol/saml/descriptor")
	port := getEnv("PORT", "8000")
	redirectURL := getEnv("REDIRECT_URL", "http://localhost:8000")

	keyPair, err := tls.LoadX509KeyPair(certificate, privateKey)
	if err != nil {
		panic(err) // TODO handle error
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		panic(err) // TODO handle error
	}

	idpMetadataURL, err := url.Parse(idpURL)
	if err != nil {
		panic(err) // TODO handle error
	}

	rootURL, err := url.Parse(redirectURL)
	if err != nil {
		panic(err) // TODO handle error
	}

	samlSP, _ := samlsp.New(samlsp.Options{
		URL:            *rootURL,
		Key:            keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate:    keyPair.Leaf,
		IDPMetadataURL: idpMetadataURL,
	})

	// hello handler
	appHandler := http.HandlerFunc(hello)

	// revert handler
	revertHandler := http.HandlerFunc(revert)

	// static contents
	staticHandler := http.FileServer(http.Dir("public/"))

	// http handle
	http.Handle("/hello", samlSP.RequireAccount(appHandler))
	http.Handle("/revert", samlSP.RequireAccount(revertHandler))
	http.Handle("/", samlSP.RequireAccount(staticHandler))
	http.Handle("/saml/", samlSP)

	log.Print("Running server at :", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
