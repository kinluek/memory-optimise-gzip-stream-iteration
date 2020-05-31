PHONY: build run profile

build:
	go build -o bin/profile ./cmd/profile/main.go

run:
	./bin/profile -memprofile ./bin/mem.prof

profile:
	go tool pprof --alloc_space -trim=false ./bin/profile ./bin/mem.prof