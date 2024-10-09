#!/bin/bash

go test -v -coverpkg=gitlab.com/flarenetwork/fdc/fdc-client/... -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
