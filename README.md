# go-compare
A simple utility to compare dependencies of bunch of go projects and report the common deps. This repo is go version of [goanalysis](github.com/deep2essence/goanalysis) written in python.
### Usage
```
[install]
$ go get -u github.com/deep2essence/go-compare

[run]
$ go-compare repo.lst --ignore-version

[debug]
$ go run main.go repo.lst
```


