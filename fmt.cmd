go install mvdan.cc/gofumpt@latest
gofumpt -l -w .
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
pause