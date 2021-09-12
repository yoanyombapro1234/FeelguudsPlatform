#!/bin/bash

protoc -I. \
	-I$(GOPATH)/src \
	-I=$(GOPATH)/src/github.com/infobloxopen/protoc-gen-gorm \
	-I=$(GOPATH)/src/github.com/infobloxopen/atlas-app-toolkit \
	-I=$(GOPATH)/src/github.com/lyft/protoc-gen-validate/validate/validate.proto \
	-I=$(GOPATH)/src/github.com/infobloxopen/protoc-gen-gorm/options \
	-I=$(GOPATH)/src/github.com/protobuf/src/google/protobuf/timestamp.proto \
	--gogoopsee_out=plugins=grpc+graphql,Mopsee/protobuf/opsee.proto=github.com/opsee/protobuf/opseeproto,Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:./internal/merchant/ --proto_path=$(GOPATH)/src:./internal/merchant/*.proto

protoc -I. \
		-I$(GOPATH)/src \
		-I=$(GOPATH)/src/github.com/infobloxopen/protoc-gen-gorm \
		-I=$(GOPATH)/src/github.com/infobloxopen/atlas-app-toolkit \
		-I=$(GOPATH)/src/github.com/lyft/protoc-gen-validate/validate/validate.proto \
		-I=$(GOPATH)/src/github.com/infobloxopen/protoc-gen-gorm/options \
		-I=$(GOPATH)/src/github.com/protobuf/src/google/protobuf/timestamp.proto \
		--proto_path=${GOPATH}/src/github.com/gogo/protobuf/protobuf \
		--govalidators_out=./internal/merchant/ \
		--go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:./internal/merchant \
		--gorm_out="engine=postgres:./internal/merchant/" ./internal/merchant/merchant.proto
