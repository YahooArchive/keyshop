// Copyright 2015 Yahoo!
// Author:  David Leon Gil (dgil@yahoo-inc.com)
// License: Apache 2
package ks

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/boltdb/bolt"
	"github.com/golang/glog"
	"github.com/yahoo/keyshop/ks/kauth"
)

var (
	db *bolt.DB
	ka *kauth.Kauth
)

func initStorage() {
	glog.Infof("initializing storage")
	db, err := bolt.Open(Config.DbFn, 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	ks = &state{db: db}
	if err != nil {
		glog.Fatalf("couldn't open keystore database at %s: %s", Config.DbFn, err)
	}
	if err != nil {
		glog.Fatalf("error initializing buckets")
	}
	glog.Infof("successfully initialized storage")
}

func initKauth() {
	glog.Infof("initializing stub key authority")
	b, err := ioutil.ReadFile(Config.KauthFn)
	if err != nil {
		panic(fmt.Sprintf("error reading kauth PEM file: %s", err))
	}
	ka, err = kauth.New(b)
	if err != nil {
		panic(fmt.Sprintf("error parsing kauth PEM file: %s", err))
	}
}

func init() {
	glog.Infof("starting server")
	initStorage()
	initKauth()
}
