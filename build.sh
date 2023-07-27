pwd
protoc --version
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u google.golang.org/grpc
go mod vendor
protoc --go_out=./ --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./replica/src/consensus.proto
protoc --go_out=./ ./proto/client/client.proto
go build -v -o ./replica/bin/replica ./replica/
go build -v -o ./client/bin/client ./client/