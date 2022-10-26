.PHONY: proto clean

install-proto-plugin:
	go get \
           		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
           		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
           		google.golang.org/protobuf/cmd/protoc-gen-go \
           		google.golang.org/grpc/cmd/protoc-gen-go-grpc && \
	go install \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest \
		google.golang.org/protobuf/cmd/protoc-gen-go@latest \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


install-proto-include:
	git clone https://github.com/googleapis/googleapis.git ${GOPATH}/pkg/mod/github.com/googleapis/googleapis && \
	git clone https://github.com/grpc-ecosystem/grpc-gateway.git ${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway

clean:
	@rm -f api/*/*/*.go

proto: clean
	@protoc \
     	--proto_path=$(GOPATH)/pkg/mod/github.com/googleapis/googleapis \
     	--proto_path=./api/im_balance/v1/proto \
     	--go_out=. \
     	--go-grpc_out=. \
     	--grpc-gateway_out=. \
	api/*/*/*/*.proto
