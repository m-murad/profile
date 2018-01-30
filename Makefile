.PHONY: proto

proto:
	protoc -I. -Ivendor  --go_out=plugins=grpc:. --govalidators_out=. proto/geoip/*.proto
	protoc --go_out=plugins=grpc,import_path=profile:. -I. -I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway \
		proto/profile/*.proto

protogw:
	protoc --go_out=plugins=grpc,import_path=profile:. -I. -I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway \
		--swagger_out=logtostderr=true:. \
		--grpc-gateway_out=logtostderr=true:. \
		proto/profile/*.proto



build:
	docker build -t profile:latest .
	docker tag profile:latest gcr.io/fish2-14312551/test_profile:v1
	gcloud docker -- push gcr.io/fish2-14312551/test_profile:v1

