SVC_NAME=vids-server
CLIENT_NAME=vids-client

CUR_DIR=$(shell pwd)
OUT_DIR=$(CUR_DIR)/bin

CURRENT_BRANCH := $(shell git symbolic-ref --short HEAD)
GITCOMMIT := $(shell git rev-parse --short HEAD)

DATE := $(shell date +%Y%m%d)

vids-server:
	@mkdir -p $(OUT_DIR)
	go build -ldflags="-w -s" -o $(OUT_DIR)/$(SVC_NAME) ./cmd/$(SVC_NAME)/
	@echo " --->  build done!\n"

vids-client:
	@mkdir -p $(OUT_DIR)
	go build -ldflags="-w -s" -o $(OUT_DIR)/$(CLIENT_NAME) ./cmd/$(CLIENT_NAME)/
	@echo " --->  build done!\n"

# nohup $(OUT_DIR)/$(SVC_NAME) run --repo $(CUR_DIR)/.vids/ --config $(CUR_DIR)/config.toml --debug > $(CUR_DIR)/.vids/vids.log &
run-server:
	$(OUT_DIR)/$(SVC_NAME) run --repo $(CUR_DIR)/.vids/ --config $(CUR_DIR)/config.toml --debug

run-client:
	$(OUT_DIR)/$(CLIENT_NAME) msql-import --host 127.0.0.1 --user root --pass root --dbname test_db --tbnames test_tb1,test_tb2

stop-server:
	pkill -9 $(SVC_NAME)