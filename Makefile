local:
	air

build:
	go build -o rsvp cmd/main.go

run: build
	./rsvp