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

## TODO for open-source version

Well, despite the disclaimer above, I probably will:

- Add an API spec in some format that can be pretty-printed.
- Add sanitized test data.
- Clean up and release the API conformance-test driver (which would
  then, effectively, be an implementation of a client in Python)
- Add directions on how to build for folks who aren't hip enough to
  know how to go get things.
