// Copyright 2015 Yahoo
// Author:  David Leon Gil (dgil@yahoo-inc.com)
// License: Apache 2

// Package yenc implements Yahoo-internal baseX encodings, as well
// as some other specialized baseX encodings used at Yahoo.
// (This variant is cleaned up for vendoring into Keyshop.)
package yenc

import (
	"github.com/yahoo/keyshop/yenc/base64"
)

const (
	DotPadding rune = '.' // Dot padding for Closure64
)

// Closure64 is the Closure library's dot-padded URL-safe base64
var Closure64 = base64.RawURLEncoding.WithPadding(DotPadding)

// Some aliases to make code more concise:

// Std64 is an alias for encoding/base64.StdEncoding
var Std64 = base64.StdEncoding

// URL64 is an alias for encoding/base64.URLEncoding
var URL64 = base64.URLEncoding

// RawStd64 is an alias for encoding/base64.RawStdEncoding
var Raw64 = base64.StdEncoding

// RawURL64 is an alias encoding/base64.RawURLEncoding
var RawURL64 = base64.URLEncoding
