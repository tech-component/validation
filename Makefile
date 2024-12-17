.PHONY: golang-imports test test-html-output

golang-imports:
	goimports -w .
test:
	go test -cover ./...

test-html-output:
	go test -coverprofile=c.out ./... && go tool cover -html=c.out && rm -f c.out
