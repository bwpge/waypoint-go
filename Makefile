APP_NAME := waypoint
BUILD_DIR := build
BIN_PATH := $(BUILD_DIR)/$(APP_NAME)
INSTALL_PATH := /usr/local/bin/$(APP_NAME)
SERVICE_NAME := $(APP_NAME).service
SERVICE_PATH := /etc/systemd/system/$(SERVICE_NAME)

.PHONY: all build install service clean

all: build

build:
	go build -o $(BIN_PATH) .

install: build
	sudo cp $(BIN_PATH) $(INSTALL_PATH)
	sudo chown root:root $(INSTALL_PATH)
	sudo chmod 755 $(INSTALL_PATH)

service:
	sudo cp $(SERVICE_NAME) $(SERVICE_PATH)
	sudo chown root:root $(SERVICE_PATH)
	sudo chmod 755 $(SERVICE_PATH)
	sudo systemctl enable $(APP_NAME)
	sudo systemctl restart $(APP_NAME)

clean:
	rm -rf $(BUILD_DIR)
