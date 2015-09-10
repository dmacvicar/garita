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

import "testing"

func Test_MD5Crypt(t *testing.T) {
	test_cases := [][]string{
		{"apache", "$apr1$J.w5a/..$IW9y6DR0oO/ADuhlMF5/X1"},
		{"pass", "$1$YeNsbWdH$wvOF8JdqsoiLix754LTW90"},
		{"topsecret", "$apr1$JI4wh3am$AmhephVqLTUyAVpFQeHZC0"},
	}
	for _, tc := range test_cases {
		e := NewMD5Entry(tc[1])
		result := MD5Crypt([]byte(tc[0]), e.Salt, e.Magic)
		if string(result) != tc[1] {
			t.Fatalf("MD5Crypt returned '%s' instead of '%s'", string(result), tc[1])
		}
		t.Logf("MD5Crypt: '%s' (%s%s$) -> %s", tc[0], e.Magic, e.Salt, result)
	}
}
