build:
    templ generate 
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -buildvcs=false -trimpath -ldflags="-w -s" -o ./bin/portfolio ./main.go
run: build
    mkdir -p data
    ./bin/portfolio

generate:
	templ generate

test:
    go test -v ./...

benchmark:
    go test -bench=. -benchmem -benchtime=1s -count 2 -cpu 1 ./...
