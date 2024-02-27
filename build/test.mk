run-tests:
	$(GO) test ./... -coverprofile=./coverage.out
	$(GO) tool cover -html=./coverage.out -o ./coverage.html