# Keyshop: a stub keyserver for E2E

This is a tiny Golang keyserver stub provided for testing the
E2E extension.

I don't anticipate pushing much new code to this branch myself;
but I'd happily accept pull requests to make this actually useful
(and provide, like, actual security properties).

**Before you even think about deploying this accessible on anything
but localhost, you MUST do something like**

    ag 'FIXME|OSS' --go -C4

**read the comments, and fix the FIXMEs.**

If you want to deploy this in a realistic environment, you may
want to privilege-separate the key authority. PRs to support using
cloud HSMs would be accepted.

**Note: DO NOT deploy this in production using Go 1.4, unless you
incorporate the broken-random-safe ECDSA patch [here](https://go-review.googlesource.com/#/c/3340/)**

(It's called keyshop because I have a low opinion of certificate
authorities, which this essentially is.)

## Building

Install an official [tarball](http://golang.org/doc/install#install)
or use your platform-of-choice's package manager. Or do the right
thing, and build [Go from source](https://go.googlesource.com).

Just run:

    go get github.com/yahoo/keyshop/ks/cmd/...

and you'll have three new Go binaries in your `$GOPATH/bin`. If
you are Yahoo-internal, you probably want to clone this repo to
its import path. E.g.:

    git clone $PARANOIDS_GIT/keyshop-oss.git \
      ${GOPATH}/src/github.com/yahoo/keyshop

## Using

Add `${GOPATH}/bin` to your path and

    cd ${GOPATH}/src/github.com/yahoo/keyshop
    genkauth
    ./scripts/mktls.sh
    ks -alsologtostderr -v 4 -log_dir ./data/logs
    

## TODO for open-source version

Well, despite the disclaimer above, I probably will:

- Add an API spec in some format that can be pretty-printed.
- Add sanitized test data.
- Clean up and release the API conformance-test driver (which would
  then, effectively, be an implementation of a client in Python).
