#!/bin/sh

# Two separate files, "clonebeegowebapp.exe" for windows and "clonebeegowebapp" for linux
go build -o $GOPATH/bin/clonebeegowebapp.exe github.com/francoishill/beegowebapp/_clonecmd
go build -o $GOPATH/bin/clonebeegowebapp github.com/francoishill/beegowebapp/_clonecmd
