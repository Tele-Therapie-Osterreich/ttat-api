.PHONY: gen run dbtest mocks test

gen:
	go generate ./...

# VB_TEST_DB='postgres://localhost/hh_test?sslmode=disable'
dbtest:
	go test -v ./test/db | ../../dev-tools/test-colours

mocks:
	mockery -recursive -name "$(shell echo "$(MOCKS)" | tr ' ' '|')"

test:
	go generate ./...
	go test ./...
