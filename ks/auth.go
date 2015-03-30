// Copyright 2015 Yahoo
// Author:  David Leon Gil (dgil@yahoo-inc.com)
// License: Apache 2
package ks

import (
	"net/http"

	"github.com/golang/glog"
)

type handler func(w http.ResponseWriter, r *http.Request)

func requireAuth(f handler, forwrite bool) handler {
	if Config.SkipAuth {
		glog.Infof("requireAuth: skipping auth due to configuration")
		return func(w http.ResponseWriter, r *http.Request) {
			glog.Infof("NOAUTH: request %+v", r)
			f(w, r)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// When authentication is required, minimize what's logged to prevent
		// logging usable authentication information.

		// This is where you'd implement some sort of authentication scheme.
		// Sorry, no implementation for Yahoo-external users provided just
		// yet.
		f(w, r)
		return
	}
}
