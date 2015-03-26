// Copyright 2015 Yahoo!
// Author:  David Leon Gil (dgil@yahoo-inc.com)
// License: Apache 2
package kauth

import (
	"github.com/golang/glog"
	"github.com/square/go-jose"
)

// A stub for a proper privileged-separated key authority.
type Kauth struct {
	signer jose.Signer
}

// Sign submits a message to the key authority for signing.
// (In this case, it just signs it...)
func (a *Kauth) Sign(msg []byte) (b []byte, err error) {
	glog.Infof("msg: %s", msg)
	obj, err := a.signer.Sign(msg)
	if err != nil {
		glog.Errorf("error signing message: %s", err)
		return
	}
	s, err := obj.CompactSerialize()
	if err != nil {
		glog.Errorf("error serializing object: %s", err)
	}
	return []byte(s), nil
}

// New initializes a new key authority from a PEM file
// containing the authority's private key.
func New(kauthPem []byte) (ka *Kauth, err error) {
	priv, err := jose.LoadPrivateKey(kauthPem)
	if err != nil {
		return
	}
	signer, err := jose.NewSigner(jose.ES256, priv)
	if err != nil {
		return
	}
	ka = &Kauth{signer: signer}
	return
}
