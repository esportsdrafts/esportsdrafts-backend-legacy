
SERVICE_DIR 	:= $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
SERVICE_NAME 	:= $(shell basename $(SERVICE_DIR))
OUTPUT_FILE 	?= api/$(SERVICE_NAME).gen.go
SPEC_FILE 	?= schemas/$(SERVICE_NAME).yaml

.PHONY: codegen

codegen:
	oapi-codegen --package $(SERVICE_NAME) -generate types,server,spec $(SPEC_FILE) > $(OUTPUT_FILE)

