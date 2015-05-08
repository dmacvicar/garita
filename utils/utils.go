//
// Copyright (c) 2015 SUSE LLC. All rights reserved.
//
// Contains code based on libtrust
// Copyright 2014 Docker, Inc.
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
package utils

import (
	"strings"
	"encoding/base32"
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/json"
)

func keyIDEncode(b []byte) string {
	s := strings.TrimRight(base32.StdEncoding.EncodeToString(b), "=")
	var buf bytes.Buffer
	var i int
	for i = 0; i < len(s)/4-1; i++ {
		start := i * 4
		end := start + 4
		buf.WriteString(s[start:end] + ":")
	}
	buf.WriteString(s[i*4:])
	return buf.String()
}

func KeyIDFromCryptoKey(pubKey crypto.PublicKey) (string, error) {
	// Generate and return a 'libtrust' fingerprint of the public key.
	// For an RSA key this should be:
	//   SHA256(DER encoded ASN1)
	// Then truncated to 240 bits and encoded into 12 base32 groups like so:
	//   ABCD:EFGH:IJKL:MNOP:QRST:UVWX:YZ23:4567:ABCD:EFGH:IJKL:MNOP
	derBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", err
	}
	hasher := crypto.SHA256.New()
	hasher.Write(derBytes)
	return keyIDEncode(hasher.Sum(nil)[:30]), nil
}

func PrettyPrint(x interface{}) string {
	b, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}
