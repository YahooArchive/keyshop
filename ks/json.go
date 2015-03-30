// Copyright 2015 Yahoo
// Author:  David Leon Gil (dgil@yahoo-inc.com)
// License: Apache 2
package ks

// A DKey represents the key for a single device.
type DKey struct {
	DeviceID  string `json:"deviceid"`
	Key       string `json:"key"`
	Timestamp int64  `json:"t"`
	UserID    string `json:"userid"`
}

// UKeys represents a keyset for a single user.
type UKeys struct {
	Timestamp int64             `json:"t"`
	UserID    string            `json:"userid"`
	Keys      map[string]string `json:"keys"`
}
