
BIN=fujlex
BIN_DIR=bin
SRC_DIR=src

build:
	@echo '---> Building Linux Binary...'
	go build -o $(BIN_DIR)/$(BIN) $(SRC_DIR)/*.go

run: build
	@echo '---> Running Binary...'
	@echo ''
	$(BIN_DIR)/$(BIN)
