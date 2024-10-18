#!/bin/bash

go test -v -coverpkg=github.com/flare-foundation/fdc-client/... -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
