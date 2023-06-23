
build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/manager

docker-build: build
	docker build -t shubhindia/k8s-mutating-webhook-experiments:v1 .
