RulzUrLibrary API
=================

[![Build Status](https://travis-ci.org/RulzUrLibrary/api.svg?branch=master)](https://travis-ci.org/RulzUrLibrary/api)

Golang implementation of the Rulz API


Utils
-----

Parsing log: `go run main.go | while read line; do grep '^{' <<< "$line" | jq '.'; done`

Running specific test: `go test ./test -v -run TestWishlistsUpdate` (beware this is a regex)
