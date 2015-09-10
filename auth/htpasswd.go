//
// Copyright (c) 2015 SUSE LLC. All rights reserved.
//
// Based on https://github.com/abbot/go-http-auth
// Copyright 2012-2013 Lev Shamardin
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
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"encoding/csv"
	"os"
	"strings"
)

/*
 SecretProvider is used by authenticators. Takes user name and realm
 as an argument, returns secret required for authentication (HA1 for
 digest authentication, properly encrypted password for basic).

 Returning an empty string means failing the authentication.
*/
type SecretProvider func(user, realm string) string

/*
 Common functions for file auto-reloading
*/
type File struct {
	Path string
	Info os.FileInfo
	/* must be set in inherited types during initialization */
	Reload func()
}

func (f *File) ReloadIfNeeded() {
	info, err := os.Stat(f.Path)
	if err != nil {
		panic(err)
	}
	if f.Info == nil || f.Info.ModTime() != info.ModTime() {
		f.Info = info
		f.Reload()
	}
}

/*
 Structure used for htdigest file authentication. Users map realms to
 maps of users to their HA1 digests.
*/
type HtdigestFile struct {
	File
	Users map[string]map[string]string
}

func reload_htdigest(hf *HtdigestFile) {
	r, err := os.Open(hf.Path)
	if err != nil {
		panic(err)
	}
	csv_reader := csv.NewReader(r)
	csv_reader.Comma = ':'
	csv_reader.Comment = '#'
	csv_reader.TrimLeadingSpace = true

	records, err := csv_reader.ReadAll()
	if err != nil {
		panic(err)
	}

	hf.Users = make(map[string]map[string]string)
	for _, record := range records {
		_, exists := hf.Users[record[1]]
		if !exists {
			hf.Users[record[1]] = make(map[string]string)
		}
		hf.Users[record[1]][record[0]] = record[2]
	}
}

/*
 SecretProvider implementation based on htdigest-formated files. Will
 reload htdigest file on changes. Will panic on syntax errors in
 htdigest files.
*/
func HtdigestFileProvider(filename string) SecretProvider {
	hf := &HtdigestFile{File: File{Path: filename}}
	hf.Reload = func() { reload_htdigest(hf) }
	return func(user, realm string) string {
		hf.ReloadIfNeeded()
		_, exists := hf.Users[realm]
		if !exists {
			return ""
		}
		digest, exists := hf.Users[realm][user]
		if !exists {
			return ""
		}
		return digest
	}
}

/*
 Structure used for htdigest file authentication. Users map users to
 their salted encrypted password
*/
type HtpasswdFile struct {
	File
	Users map[string]string
}

func reload_htpasswd(h *HtpasswdFile) {
	r, err := os.Open(h.Path)
	if err != nil {
		panic(err)
	}
	csv_reader := csv.NewReader(r)
	csv_reader.Comma = ':'
	csv_reader.Comment = '#'
	csv_reader.TrimLeadingSpace = true

	records, err := csv_reader.ReadAll()
	if err != nil {
		panic(err)
	}

	h.Users = make(map[string]string)
	for _, record := range records {
		h.Users[record[0]] = record[1]
	}
}

/*
 SecretProvider implementation based on htpasswd-formated files. Will
 reload htpasswd file on changes. Will panic on syntax errors in
 htpasswd files. Realm argument of the SecretProvider is ignored.
*/
func HtpasswdFileProvider(filename string) SecretProvider {
	h := &HtpasswdFile{File: File{Path: filename}}
	h.Reload = func() { reload_htpasswd(h) }
	return func(user, realm string) string {
		h.ReloadIfNeeded()
		password, exists := h.Users[user]
		if !exists {
			return ""
		}
		return password
	}
}

func NewHtpasswdValidator(filename string) Validator {
	provider := HtpasswdFileProvider(filename)

	return func(username string, passwd string) bool {
		// realm is ignored
		hashedPw := provider(username, "")

		if strings.HasPrefix(hashedPw, "{SHA}") {
			d := sha1.New()
			d.Write([]byte(passwd))

			if subtle.ConstantTimeCompare([]byte(hashedPw[5:]), []byte(base64.StdEncoding.EncodeToString(d.Sum(nil)))) != 1 {
				return false
			}
		} else {
			e := NewMD5Entry(hashedPw)
			if e == nil {
				return false
			}
			if subtle.ConstantTimeCompare([]byte(hashedPw), MD5Crypt([]byte(passwd), e.Salt, e.Magic)) != 1 {
				return false
			}
		}
		return true
	}
}
