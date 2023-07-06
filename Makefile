
BIN=4RCH
BIN_DIR=bin/lin
SRC_DIR=src

build:
	@echo '---> Building Linux Binary...'
	go build -o $(BIN_DIR)/$(BIN) $(SRC_DIR)/*.go

run: build
	@echo '---> Running Binary...'
	@echo ''
	$(BIN_DIR)/$(BIN)
