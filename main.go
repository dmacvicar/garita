//
// Copyright (c) 2015 SUSE LLC. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// package main
//
package main
import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/dgrijalva/jwt-go"
	"github.com/nu7hatch/gouuid"
)

var router *mux.Router

type AccessItem struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Actions []string `json:"actions"`
}

type TokenResponse struct {
    Token string `json:"token"`
}

func Token(w http.ResponseWriter, r *http.Request) {
	service := r.URL.Query().Get("service")
	scopeEsc := r.URL.Query().Get("scope")
	account := r.URL.Query().Get("account")

	scope, err := url.QueryUnescape(scopeEsc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	parts := strings.Split(scope, ":")
	if len(parts) != 3 {
		http.Error(w, "Invalid scope string", 500)
		return
	}

	if parts[0] != "repository" {
		http.Error(w, "Only repository access is supported", 500)
		return
	}
	repository := parts[1]

	token := jwt.New(jwt.SigningMethodHS256)
	// audience
	token.Claims["aud"] = service
	token.Claims["sub"] = account

	now := time.Now()
	// expiration
	token.Claims["exp"] = now.Add(time.Hour * 72).Unix()
	// not before
	token.Claims["nbf"] = now
	token.Claims["iat"] = now

	jti, err := uuid.NewV4()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	token.Claims["jti"] = jti
	token.Claims["access"] = [...]map[string]interface{}{
		{
			"type": "repository",
			"name": repository,
			"actions" : []string{"pull", "push"},
		},
	}

	// Sign and get the complete encoded token as a string
	pem, err := ioutil.ReadFile("/vagrant/vagrant/conf/ca_bundle/server.key")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	signed, err := token.SignedString(pem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(TokenResponse{Token: signed})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	router = mux.NewRouter()
	router.Handle("/v2/token", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(Token)))
	log.Printf("Listening...:%s", port)
	log.Fatal(http.ListenAndServe(":" +  port, router))
}

