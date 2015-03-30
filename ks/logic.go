// Copyright 2015 Yahoo
// Author:  David Leon Gil (dgil@yahoo-inc.com)
// License: Apache 2
package ks

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/golang/glog"
	"golang.org/x/crypto/openpgp"
)

func validYmail(email string) bool {
	// FIXME(OSS): Should be obvious where this goes, and
	// how you should modify it. Note, however, that Match
	// does require the ^$ bracketing in Go -- unlike in
	// Python. Yipes.
	r := regexp.MustCompile(`^[a-z][a-z0-9]{2,64}@(yahoo-inc\.com|uk\.yahoo-inc\.com|yahoo-corp\.jp|tumblr\.com|yahoo\.com)$`)
	return r.MatchString(email)
}

func validKeyForUser(userid, email string, key []byte) (ok bool) {
	el, err := openpgp.ReadKeyRing(bytes.NewBuffer(key))
	if err != nil {
		glog.Errorf("error reading keyring: %s", err)
		return false
	}
	// Check that there's only one keypair included,
	if len(el) != 1 {
		glog.Errorf("Expected one entity, got %d.\n%+v", len(el), el)
		return false
	}
	// that there's only one UID packet for the keypair,
	identities := el[0].Identities
	if len(identities) != 1 {
		glog.Errorf("Expected one identity, got %d.\n%+v", len(identities), identities)
		return false
	}
	var uidEmail string
	for _, v := range identities {
		// This loop will only execute once...
		u := v.UserId
		if u.Name != "" || u.Comment != "" {
			glog.Errorf("too many fields filled (names and comments prohibited): got %+v", u)
			return false
		}
		uidEmail = u.Email
		if uidEmail == "" || uidEmail != email {
			glog.Errorf("email address in identity did not agree with email address passed in: got %s, wanted %s", uidEmail, email)
			return false
		}
	}

	// and, finally,
	// FIXME(OSS): authentication-mechanism-specific checks.
	if userid == email {
		err = nil
	} else {
		err = fmt.Errorf("userid != email")
	}

	if err != nil {
		glog.Error(err)
		return false
	}
	return true
}
