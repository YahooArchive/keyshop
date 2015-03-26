package ks

// DKeys represent the key for a single device.
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
