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

(It's called keyshop because I have a low opinion of certificate
authorities, which this essentially is.)

## Building

Install Go following [golang.org's instructions](https://golang.org/doc/install)
or use your platform-of-choice's package manager.

Just run:

    go get github.com/yahoo/keyshop/ks/cmd/...

and you'll have three new Go binaries in your `$GOPATH/bin`: genkauth, localcert, and ks.

## Using

Add `${GOPATH}/bin` to your path and

    cd ${GOPATH}/src/github.com/yahoo/keyshop
    genkauth
    ./scripts/mktls.sh
    ks -alsologtostderr -v 4 -log_dir ./data/logs
