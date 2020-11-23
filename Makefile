
test:
	go test -v ./...

build:
	go build ./...

bench:
	go test -benchmem -run=^$ github.com/BenJoyenConseil/rmi -bench Benchmark*