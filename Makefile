build: build-app
run:
	make build-app && ./astra-test
build-app:
	go build -o astra-test
clean:
	rm -f as && go clean --modcache && go mod tidy
