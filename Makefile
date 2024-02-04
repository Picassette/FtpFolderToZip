PROJECT_VERSION=0.1.0

ifeq ($(OS),Windows_NT)
    # Windows
	LOCAL_PATH=$(shell echo %cd%)
	SRC_PATH=$(LOCAL_PATH)\src
	BINARY_PATH=\dist\win
	BINARY_NAME=ftp2zip-$(PROJECT_VERSION).exe
	MKDIR=mkdir
	RMDIR=rmdir /s /q
	RM=del /q
	MOVE=move
else
    # Linux
	LOCAL_PATH=$(shell pwd%)
	SRC_PATH=$(LOCAL_PATH)/src/
	BINARY_PATH=./dist/linux/
	BINARY_NAME=ftp2zip-$(PROJECT_VERSION)
	MKDIR=mkdir -p
	RMDIR=rm -rf
	RM=rm -f
	MOVE=mv
endif

all: build

build:
	go build -C "$(SRC_PATH)" -o "$(LOCAL_PATH)$(BINARY_PATH)\$(BINARY_NAME)"

clean:
	go clean
	$(RM) $(BINARY_NAME)

rebuild: clean build

.PHONY: all build clean rebuild