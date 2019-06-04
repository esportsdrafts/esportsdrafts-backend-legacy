
PACKAGE_NAME ?= auth

.PHONY: codegen

codegen:
	oapi-codegen --package $(PACKAGE_NAME) -generate types,server,spec schemas/auth.yaml > auth/auth.gen.go
