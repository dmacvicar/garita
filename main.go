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
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"log"
	"net/http"
	"github.com/gorilla/mux"

	api "github.com/dmacvicar/garita/api"
)

var router *mux.Router

func main() {

	pwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	port := flag.Int("port", 80, "port to listen to")
	htpasswdPath := flag.String("htpasswd", path.Join(pwd, "htpasswd"), "password file")
	keyPath := flag.String("key", path.Join(pwd, "server.key"), "token secret key")
	flag.Parse()

	if _, err := os.Stat(*htpasswdPath); os.IsNotExist(err) {
		fmt.Printf("no such file or directory: %s", *htpasswdPath)
		return
	}

	if _, err := os.Stat(*keyPath); os.IsNotExist(err) {
		fmt.Printf("no such file or directory: %s", *keyPath)
		return
	}

	tokenHandler := api.NewGaritaTokenHandler(*htpasswdPath, *keyPath)

	router = mux.NewRouter()
	router.Handle("/v2/token", tokenHandler)
	log.Printf("Listening...:%d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), router))
}

