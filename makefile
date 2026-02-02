build :
	go build -o blockchain
run: build
	./blockchain
test:
	go test -v ./...