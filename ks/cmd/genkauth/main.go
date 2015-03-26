// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A lightly modified clone of tls/generate_cert.go from
// the Go stdlib. (dlg)

package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const prefix = "data/"

var (
	ecdsaCurve = flag.String("ecdsa-curve", "", "ECDSA curve to use to generate a key. Valid values are P224, P256, P384, P521")
)

func pemBlockForPrivateKey(priv *ecdsa.PrivateKey) *pem.Block {
	b, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
		os.Exit(2)
	}
	return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
}

func pemBlockForPublicKey(pub *ecdsa.PublicKey) *pem.Block {
	b, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA public key: %v", err)
		os.Exit(2)
	}
	return &pem.Block{Type: "EC PUBLIC KEY", Bytes: b}
}

// FIXME(dlg): This is hideous.
func arrayForPublicKey(pub *ecdsa.PublicKey) string {
	raw := elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	vals := make([]string, len(raw))
	for i, b := range raw {
		vals[i] = strconv.Itoa(int(b))
	}
	return "[" + strings.Join(vals, ", ") + "]"
}

func main() {
	flag.Parse()

	var priv *ecdsa.PrivateKey
	var err error
	switch *ecdsaCurve {
	case "P224":
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		log.Printf("using P256 by default")
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	}
	if err != nil {
		log.Fatalf("failed to generate private key: %s", err)
	}

	keyOut, err := os.OpenFile(prefix+"kauth.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Print("failed to open kauth.pem for writing:", err)
		return
	}
	pem.Encode(keyOut, pemBlockForPrivateKey(priv))
	keyOut.Close()
	log.Print("wrote kauth.pem\n")

	keyOut, err = os.OpenFile(prefix+"kauth.pem.pub", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Print("failed to open kauth.pem.pub for writing:", err)
		return
	}
	pem.Encode(keyOut, pemBlockForPublicKey(&priv.PublicKey))
	keyOut.Close()

	arrayOut, err := os.OpenFile(prefix+"kauth.pub.js", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	arrayOut.Write([]byte(arrayForPublicKey(&priv.PublicKey)))
	arrayOut.Close()

	log.Print("wrote kauth.pem.pub\n")
}
