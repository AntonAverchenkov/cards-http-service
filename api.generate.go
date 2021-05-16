package main

//go:generate $GOPATH/bin/oapi-codegen --package api --generate types  -o internal/api/types.go  api.yaml
//go:generate $GOPATH/bin/oapi-codegen --package api --generate server -o internal/api/server.go api.yaml
//go:generate $GOPATH/bin/oapi-codegen --package api --generate spec   -o internal/api/spec.go   api.yaml
