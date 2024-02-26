BIN_AMD64 = bin/amd64
PREFIX_AMD64 = GOOS=linux GOARCH=amd64

amd64: main-ethInterfacing-amd64

main-ethInterfacing-amd64:
	@$(PREFIX_AMD64) $(GO) build -o $(BIN_AMD64)/main-ethInterfacing openpixelsystems.org/go-networkmanager-lib/cmd/mainethInterfacing/
	@echo Compiled $(BIN_AMD64)/main-ethInterfacing with version \'$(COMMIT_ID)\'