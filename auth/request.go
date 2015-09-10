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

import "net/http"

/*
 Request handlers must take AuthenticatedRequest instead of http.Request
*/
type AuthenticatedRequest struct {
	http.Request
	/*
	 Authenticated user name. Current API implies that Username is
	 never empty, which means that authentication is always done
	 before calling the request handler.
	*/
	Username string
}

/*
 AuthenticatedHandlerFunc is like http.HandlerFunc, but takes
 AuthenticatedRequest instead of http.Request
*/
type AuthenticatedHandlerFunc func(http.ResponseWriter, *AuthenticatedRequest)
