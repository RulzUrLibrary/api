RulzUrLibrary API
=================

Golang implementation of the Rulz API


Utils
-----

Parsing log: `go run main.go | while read line; do grep '^{' <<< "$line" | jq '.'; done`
