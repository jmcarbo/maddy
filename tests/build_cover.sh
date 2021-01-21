#!/bin/sh
exec go test -tags 'cover_main debugflags' -coverpkg 'github.com/jmcarbo/maddy,github.com/jmcarbo/maddy/pkg/...,github.com/jmcarbo/maddy/internal/...' -cover -covermode atomic -c cover_test.go -o maddy.cover
