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
package api

import (
	"encoding/json"
	"github.com/dmacvicar/garita/auth"
	"github.com/dmacvicar/garita/token"
	"github.com/dmacvicar/garita/utils"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"os"
)

type TokenResponse struct {
	Token string `json:"token"`
}

type tokenHandler struct {
	keyPath      string
	htpasswdPath string
}

func createAuthTokenFunc(keyPath string) func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	return func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
		service := r.URL.Query().Get("service")
		scope, _ := token.ParseScope(r.URL.Query().Get("scope"))

		token, err := token.NewJwtToken(r.Username, service, scope, keyPath)
		log.Println(utils.PrettyPrint(token.Claim()))

		if err != nil {
			log.Println("error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		signed, err := token.SignedString()
		if err != nil {
			log.Println("error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		js, err := json.Marshal(TokenResponse{Token: signed})
		if err != nil {
			log.Println("error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func NewGaritaTokenHandler(htpasswdPath string, keyPath string) http.Handler {
	validator := auth.NewHtpasswdValidator(htpasswdPath)
	tokenHandler := auth.BasicAuth(createAuthTokenFunc(keyPath), "realm", validator)
	logHandler := handlers.LoggingHandler(os.Stdout, auth.GetOnly(tokenHandler))
	return logHandler
}
