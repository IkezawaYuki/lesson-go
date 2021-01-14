#!/bin/zsh

go test -bench=.

echo "32 Bit version"
GOARCH=386 go test -bench=. bench_test.go intset.go mapintset.go