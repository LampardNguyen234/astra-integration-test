build: build-app
run:
	make build-app && ./astra-test
build-app:
	go build -o astra-test
clean:
	rm -f astra-test && go clean --modcache && go mod tidy
