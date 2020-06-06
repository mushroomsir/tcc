.PHONY: test
test: 
	go test -p 1 ./...


.PHONY: coverhtml
coverhtml:
	@mkdir -p coverage
	@go test -p 1 -coverprofile=coverage/cover.out ./...
	@go tool cover -html=coverage/cover.out -o coverage/coverage.html
	@go tool cover -func=coverage/cover.out | tail -n 1

.PHONY: lint
lint:
	@hash golint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u golang.org/x/lint/golint; \
	fi
	@golint -set_exit_status ${PKG_LIST}