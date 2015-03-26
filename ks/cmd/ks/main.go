// Copyright 2015 Yahoo!
// Author:  David Leon Gil (dgil@yahoo-inc.com)
// License: Apache 2
package main

import (
	"crypto/tls"
	"encoding/pem"
	"flag"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/yahoo/keyshop/ks"
	"golang.org/x/crypto/sha3"
)

// ServeIzkp returns an http.Handler that reads an input file and
// computes an interactive zero-knowledge proof-of-posession protocol.
// (This is completely unused, but isn't it cool?)
func ServeIzkp(fn string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadFile(fn)
		if err != nil {
			glog.Errorf("error reading file %s: %s", fn, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		chalString := r.Header.Get("x-izkp-challenge")
		if chalString == "" {
			glog.Infof("didn't receive a challenge, so using a raw hash")
			d := make([]byte, 64)
			sha3.ShakeSum256(d, b)
			w.Write(d)
			return
		}
		challenge := []byte(chalString)
		glog.Infof("received a challenge of length %d", len(challenge))
		h := sha3.New512()
		h.Write(challenge)
		h.Write(b)
		d := make([]byte, 64)
		h.Sum(d)
		w.Write(d)
		return
	}
}

func serveBytes(b []byte) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(b)
		return
	}
}

func serveFile(fn, mime string) func(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		panic("couldn't open file: " + fn)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", mime)
		w.Write(b)
		return
	}
}

func readPem(fn string) (raw, decoded []byte) {
	fn = ks.Config.TLSPrefix + fn
	raw, err := ioutil.ReadFile(fn)
	if err != nil {
		glog.Fatalf("error loading PEM from %s", fn)
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		glog.Fatalf("couldn't parse PEM from %s", fn)
	}
	decoded = block.Bytes
	return
}

func main() {
	// Parse flags for Glog.
	flag.Parse()

	var err error
	runtime.GOMAXPROCS(16)

	chainPem, chainDer := readPem("chain.pem")

	// Handle certificate paths
	p := mux.NewRouter()

	c := p.PathPrefix("/-").Subrouter()
	c.HandleFunc("/chain.pem", serveBytes(chainPem)).Methods("GET")
	c.HandleFunc("/chain.der", serveBytes(chainDer)).Methods("GET")

	// Set up the subrouter for the keyshop
	r := p.PathPrefix("/v1/k").Subrouter()
	r.HandleFunc("/{userid}", ks.Get).Methods("GET")
	r.HandleFunc("/{userid}/{deviceid}", ks.Post).Methods("POST")

	s := &http.Server{
		Addr:           ks.Config.Addr,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	glog.Infof("before handle")
	http.Handle("/", p)

	if ks.Config.UseTLS {
		prefix := ks.Config.TLSPrefix
		s.TLSConfig = &tls.Config{
			SessionTicketsDisabled: true,
			MinVersion:             tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				/*
					tls.CurveP384, // the high security-strength curves don't have
					tls.CurveP521, // constant-time implementations in Go at present
				*/
			},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			},
			PreferServerCipherSuites: true,
		}
		glog.Infof("starting to serve %s:\n%+v\n", ks.Config.Addr, ks.Config, s.TLSConfig)

		err = s.ListenAndServeTLS(prefix+"chain.pem", prefix+"privatekey.pem")
	} else {
		glog.Infof("starting to serve raw http")
		err = s.ListenAndServe()
	}
	if err != nil {
		glog.Fatalf("Error starting server: %s\n", err)
	}
}
