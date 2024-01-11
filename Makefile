.SILENT:

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o ./bin/fs_syncer ./cmd/fs_syncer
clean:
	rm -f /bin/fs_syncer
run:
	./bin/fs_syncer --src-folder-path=$(src-folder-path) --output-folder-path=$(output-folder-path) --stdout-log-enable=$(stdout-log-enable) --log-lvl=$(log-lvl)
prod:
	make build && make run src-folder-path=./src-folder output-folder-path=./output-folder stdout-log-enable=true log-lvl=info
test:
	go test ./...
bench: 
	go test ./... -bench=. -benchmem -benchtime 10s 
