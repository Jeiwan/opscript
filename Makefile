default:
	cat Makefile

spec:
	go run *.go buildspec
	go-bindata -pkg internal -o internal/bindata.go spec.json

test:
	go test -v ./...