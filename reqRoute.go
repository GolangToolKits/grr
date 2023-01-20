package grrt

// Copyright 2022 GolangToolKits Authors. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

import (
	"log"
	"net/http"
)

// ReqRoute ReqRoute
type ReqRoute struct {
	//namedRoutes map[string]*Route
	handler      http.Handler
	host         string
	path         string
	matcher      Matcher
	active       bool
	pathVarsUsed bool
}

// New New
func (t *ReqRoute) New() Route {
	var m pathMatcher
	t.matcher = m.New()
	return t
}

// Handler Handler
func (t *ReqRoute) Handler(handler http.Handler) Route {
	if t.active {
		t.handler = handler
	}
	return t
}

// HandlerFunc HandlerFunc
func (t *ReqRoute) HandlerFunc(f func(http.ResponseWriter, *http.Request)) Route {
	return t.Handler(http.HandlerFunc(f))
}

// Path Path
func (t *ReqRoute) Path(p string) Route {
	if t.chechPath(p) && t.matcher.addPath(p) {
		t.path = p
		t.active = true
	}
	return t
}

// Host Host
func (t *ReqRoute) Host(h string) Route {
	return nil
}

// GetHandler GetHandler
func (t *ReqRoute) GetHandler() http.Handler {
	return t.handler
}

func (t *ReqRoute) chechPath(p string) bool {
	var rtn = false
	if t.chechCurlys(p) && t.chechBackSlash(p) {
		rtn = true
	} else {
		t.printError(p, "Problem with path")
	}
	return rtn
}

func (t *ReqRoute) chechCurlys(p string) bool {
	var rtn = false
	open := rune('{')
	closed := rune('}')
	var cl int = 0
	for _, c := range p {
		if c == open {
			cl++
		} else if c == closed {
			cl--
		}
	}
	if cl == 0 {
		rtn = true
	} else {
		t.printError(p, "Mismatched curly brackets")
	}
	return rtn
}

func (t *ReqRoute) chechBackSlash(p string) bool {
	var rtn = true
	bs := rune('/')
	var lastBs int = 0
	for i, c := range p {
		if i == 0 && c != bs {
			rtn = false
			t.printError(p, "Bad backslash")
			break
		} else if i != 0 {
			if c == bs && i == lastBs+1 {
				rtn = false
				t.printError(p, "Bad backslash")
				break
			} else if c == bs && i != len(p)-1 {
				lastBs = i
			} else if c == bs && i == len(p)-1 {
				rtn = false
			}
		}
	}
	return rtn
}

func (t *ReqRoute) printError(p string, cause string) {
	log.Println("Path not added to route:")
	log.Println(p)
	log.Println(cause)
}
