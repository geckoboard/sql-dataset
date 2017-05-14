BUILD_DIR=builds

build:
	mkdir -p $(BUILD_DIR)
	docker pull karalabe/xgo-latest
	go get github.com/karalabe/xgo
	xgo -dest=$(BUILD_DIR) -targets="darwin-10.10/* windows-8.0/amd64 windows-8.0/386 linux/amd64 linux/386" .

test:
	MYSQL_URL=root@/testdb?parseTime=true POSTGRES_URL=postgres://postgres@localhost:5432/testdb?sslmode=disable go test ./... -v
