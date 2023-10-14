.PHONY: all
all: tidy format vet openapi test

.PHONY: tidy
tidy:
	@echo "=> Executing go mod tidy"
	@go mod tidy

.PHONY: format
format:
	@echo "=> Formatting code and organizing imports"
	@goimports -w ./

.PHONY: vet
vet:
	@echo "=> Executing go vet"
	@go vet ./...

COMMON_FLAGS := -covermode=atomic -coverprofile=/tmp/coverage.out -coverpkg=./... -count=1 -race -shuffle=on

.PHONY: test
test:
	@go test ./... $(COMMON_FLAGS)
	@cat .covignore | sed '/^[[:space:]]*$$/d' >/tmp/covignore
	@grep -Fvf /tmp/covignore /tmp/coverage.out  > /tmp/coverage.out-filtered
	@go tool cover -func /tmp/coverage.out-filtered | grep total: | sed -e 's/\t//g' | sed -e 's/(statements)/ /'

.PHONY: test-cover
test-cover: test
	@echo "=> Running tests and generating report"
	@go tool cover -html=/tmp/coverage.out-filtered

.PHONY: run
run:
	@echo "=> Running application"
	@go run cmd/main.go

.PHONY: openapi
openapi:
	@oapi-codegen -generate types --exclude-schemas=Error -o internal/handler/openapi_types.gen.go -package handler internal/handler/openapi/spec.yaml

.PHONY: openapi-ui
openapi-ui:
	@echo "=> Running SwaggerUI"
	@echo "Open http://localhost:9999 after the message 'Configuration complete; ready for start up'"
	@echo "Press Ctrl+C to stop"
	@echo ""
	@docker run --rm -p 9999:8080 \
		-v $(PWD)/internal/handler/openapi/spec.yaml:/spec.yaml \
		-e SWAGGER_JSON=/spec.yaml \
		swaggerapi/swagger-ui 2> /dev/null | grep "Configuration complete; ready for start up"