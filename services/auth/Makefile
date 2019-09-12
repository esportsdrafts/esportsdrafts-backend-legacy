
PACKAGE_NAME 	?= auth
OUTPUT_FILE 	?= api/auth.gen.go
SPEC_FILE 	?= schemas/auth.yaml

.PHONY: codegen

codegen:
	oapi-codegen --package $(PACKAGE_NAME) -generate types,server,spec $(SPEC_FILE) > $(OUTPUT_FILE)
