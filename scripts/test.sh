
protoc -I .\
       -I"$GOPATH"/src \
       -I="$GOPATH"/src/github.com/infobloxopen/protoc-gen-gorm/options/gorm.proto \
       --govalidators_out=./internal/merchant/ \
		   --go_out=plugins=./internal/merchant \
		   --proto_path="$GOPATH"/src/github.com/gogo/protobuf/protobuf \
		   --gorm_out="engine=postgres:./internal/merchant/" \
		   --go_out=./internal/merchant ./internal/merchant/merchant.proto

protoc --go_out=.  \
			 --go_opt=paths=source_relative \
       --go-grpc_out=.
       --go-grpc_opt=paths=source_relative ./internal/merchant/merchant.proto
