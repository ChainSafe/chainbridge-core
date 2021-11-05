#!/usr/bin/env bash
go test -p=1 $(go list ./... | grep 'e2e')