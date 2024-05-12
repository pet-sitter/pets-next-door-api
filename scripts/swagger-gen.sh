#!/bin/bash
export GOROOT=$(go env GOROOT)
swag init -d ./cmd/server -o ./pkg/docs --pd
