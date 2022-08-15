//go:build tools
// +build tools

package main

import (
	_ "github.com/git-chglog/git-chglog/cmd/git-chglog"
	_ "github.com/go-swagger/go-swagger/cmd/swagger"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/mcubik/goverreport"
	_ "github.com/securego/gosec/cmd/gosec"
	_ "github.com/segmentio/golines"
	_ "golang.org/x/tools/cmd/goimports"
	_ "mvdan.cc/gofumpt"
)
