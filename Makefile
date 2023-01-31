install-protoc-gen-ts:
	cd scripts && npm install
	mkdir -p scripts/protoc
	curl -L https://github.com/protocolbuffers/protobuf/releases/download/v21.5/protoc-21.5-linux-x86_64.zip -o scripts/protoc/protoc.zip
	cd scripts/protoc && unzip -o protoc.zip
	rm scripts/protoc/protoc.zip

cosmos-version = v0.45.4

download-cosmos-proto:
	mkdir -p proto/cosmos/base/query/v1beta1
	curl https://raw.githubusercontent.com/cosmos/cosmos-sdk/${cosmos-version}/proto/cosmos/base/query/v1beta1/pagination.proto -o proto/cosmos/base/query/v1beta1/pagination.proto
	mkdir -p proto/google/api
	curl https://raw.githubusercontent.com/cosmos/cosmos-sdk/${cosmos-version}/third_party/proto/google/api/annotations.proto -o proto/google/api/annotations.proto
	curl https://raw.githubusercontent.com/cosmos/cosmos-sdk/${cosmos-version}/third_party/proto/google/api/http.proto -o proto/google/api/http.proto
	mkdir -p proto/gogoproto
	curl https://raw.githubusercontent.com/cosmos/cosmos-sdk/${cosmos-version}/third_party/proto/gogoproto/gogo.proto -o proto/gogoproto/gogo.proto

gen-protoc-ts:
	download-cosmos-proto install-protoc-gen-ts
	mkdir -p ./client/src/types/generated/
	ls proto/checkers | xargs -I {} ./scripts/protoc/bin/protoc \
		--plugin="./scripts/node_modules/.bin/protoc-gen-ts_proto" \
		--ts_proto_out="./client/src/types/generated" \
		--proto_path="./proto" \
		--ts_proto_opt="esModuleInterop=true,forceLong=long,useOptionals=messages" \
		checkers/{}

build-all:
	GOOS=linux GOARCH=amd64 go build -o ./build/checkers-linux-amd64 ./cmd/checkersd/main.go
	GOOS=linux GOARCH=arm64 go build -o ./build/checkers-linux-arm64 ./cmd/checkersd/main.go
	GOOS=darwin GOARCH=amd64 go build -o ./build/checkers-darwin-amd64 ./cmd/checkersd/main.go

do-checksum:
	cd build && sha256sum checkers-linux-amd64 checkers-linux-arm64 checkers-darwin-amd64 > checkers_checksum

build-with-checksum:
	make build-all && make do-checksum