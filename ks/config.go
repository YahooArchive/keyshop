// Copyright 2015 Yahoo
// Author:  David Leon Gil (dgil@yahoo-inc.com)
// License: Apache 2

package ks

type config struct {
	Addr      string
	DbFn      string
	KauthFn   string
	SkipAuth  bool
	TLSPrefix string
	UseTLS    bool
}

var (
	// Config contains the default configuration for the keyshop.
	Config = &config{
		// Whether to do something sane or not.
		UseTLS:   true,
		SkipAuth: false,
		// Listen on localhost only. FIXME(OSS): This may or may not
		// be particularly useful, unless you only want to talk to
		// yourself.
		Addr: "localhost:25519",
		// The location of data files.
		DbFn:      "data/25519.db",
		KauthFn:   "data/kauth/kauth.pem",
		TLSPrefix: "data/tls/localhost.",
	}
)
