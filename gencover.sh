#!/bin/bash

go test -v -coverpkg=local/fdc/... -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
