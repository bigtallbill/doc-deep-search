#!/usr/bin/env bash
GOOS=windows GOARCH=386 go build -o ./bin/docsearch-win.exe ./cmd/search
GOOS=linux go build -o ./bin/docsearch-linux ./cmd/search
