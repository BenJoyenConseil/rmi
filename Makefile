
test:
	go test ./...

build:
	go build ./...

bench:
	go test -benchmem -run=^$ github.com/BenJoyenConseil/rmi -bench Benchmark*

coverage:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...