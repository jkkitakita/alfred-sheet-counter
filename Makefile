SHELL := /bin/bash
.PHONY: all build archive
PLIST=info.plist
EXEC_BIN=alfred-sheet-counter
DIST_FILE=alfred-sheet-counter.alfredworkflow

all: build archive

build:
	go build -v -o $(EXEC_BIN) .

archive: build $(PLIST)
	zip -r $(DIST_FILE) $(PLIST) $(EXEC_BIN)
