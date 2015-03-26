// Copyright 2015 Yahoo!
// Author:  David Leon Gil (dgil@yahoo-inc.com)
// License: Apache 2
package ks

import (
	"errors"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/golang/glog"
)

var (
	ks      *state
	buckets = []string{"issued"}
	errNsk  = errors.New("no such key")
	errNsu  = errors.New("no such user")
)

type Key struct {
	Userid   []byte
	Deviceid []byte
	Key      []byte
}

type KeyShop interface {
	New(userid, deviceid, key []byte) int
	Update(userid, deviceid, key []byte) int
	Get(userid []byte) (map[string][]byte, int)
}

type state struct {
	db *bolt.DB
	KeyShop
}

func (s *state) New(userid, deviceid, key []byte) (status int) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(userid)
		if err != nil {
			return err
		}
		d, err := b.CreateBucketIfNotExists(deviceid)
		if err != nil {
			return err
		}
		d.Put([]byte("key"), key)
		t, _ := time.Now().UTC().GobEncode()
		d.Put([]byte("updated"), t)
		return nil
	})
	switch err {
	case bolt.ErrBucketExists:
		return http.StatusConflict
	case nil:
		return http.StatusCreated
	default:
		return http.StatusInternalServerError
	}
}

func (s *state) Update(userid, deviceid, key []byte) (status int) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(userid)
		if err != nil {
			return err
		}
		d := b.Bucket(deviceid)
		if b == nil {
			return errNsk
		}
		d.Put([]byte("key"), key)
		t, _ := time.Now().UTC().GobEncode()
		d.Put([]byte("updated"), t)
		return nil
	})
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func (s *state) NewOrUpdate(userid, deviceid, key []byte) (status int) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(userid)
		if err != nil {
			glog.Errorf("error creating or getting %s/%s bucket: %s", userid)
			return err
		}
		b.Put(deviceid, key)
		return nil
	})
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func (s *state) Get(userid string) (keys map[string]string, status int) {
	keys = make(map[string]string)
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userid))
		if b == nil {
			glog.Infof("Userid %s not found", userid)
			return errNsu
		}
		b.ForEach(func(k, v []byte) error {
			glog.Infof("%s/%s: %s", userid, k, v)
			keys[string(k)] = string(v)
			return nil
		})
		return nil
	})
	switch err {
	case errNsu:
		glog.Infof("user %s does not have any keys", userid)
		return nil, http.StatusNotFound
	case nil:
		glog.Infof("found keys for user %s: %s", userid, keys)
		return keys, http.StatusOK
	default:
		glog.Infof("error trying to get keys for %s: %s", userid, err)
		return nil, http.StatusInternalServerError
	}
}
