#!/bin/bash

cd flow/
cocaine-tool app upload --name flow-tools
cocaine-tool app restart --name flow-tools --profile TEST --timeout=4

go test -coverprofile=coverage.out  github.com/cocaine/cocaine-flow/backend || cocaine-tool app remove --name bullet_first
go test -coverprofile=coverage.out  github.com/cocaine/cocaine-flow/frontHTTP || cocaine-tool app remove --name bullet_first
