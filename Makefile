build:
    go build -o chatroom-app ./cmd/main.go

run:
    ./chatroom-app

test:
    go test ./...

clean:
    rm -f chatroom-app