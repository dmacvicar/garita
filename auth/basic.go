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
package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func BasicAuth(pass AuthenticatedHandlerFunc, realm string, validator Validator) http.HandlerFunc {

	sendUnauthorized := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
		http.Error(w, "authorization failed", http.StatusUnauthorized)
	}

	return func(w http.ResponseWriter, r *http.Request) {

		if len(r.Header["Authorization"]) < 1 {
			sendUnauthorized(w, r)
			return
		}

		authHeader := r.Header["Authorization"][0]
		auth := strings.SplitN(authHeader, " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(w, "bad syntax", http.StatusBadRequest)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !validator(pair[0], pair[1]) {
			sendUnauthorized(w, r)
			return
		}

		//w.Header().Set("X-Authorized-Username", pair[0])
		ar := &AuthenticatedRequest{Request: *r, Username: pair[0]}
		pass(w, ar)
	}
}
