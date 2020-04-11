.PHONY: gen run dbtest mocks test

gen:
	go generate ./...

MOCKS=DB Mailer

mocks:
	mockery -recursive -name "$(shell echo "$(MOCKS)" | tr ' ' '|')"

test:
	go generate ./...
	go test ./...
