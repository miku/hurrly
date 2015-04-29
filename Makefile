SHELL := /bin/bash
TARGETS = hurrly

# http://docs.travis-ci.com/user/languages/go/#Default-Test-Script
test:
	go test -v ./...

bench:
	go test -bench=.

imports:
	goimports -w .

fmt:
	go fmt ./...

vet:
	go vet ./...

all: fmt test
	go build ./...

install:
	go install ./...

clean:
	go clean
	rm -f coverage.out
	rm -f $(TARGETS)
	rm -f hurrly-*.x86_64.rpm
	rm -f hurrly_*.deb
	rm -rf packaging/deb/hurrly/usr

cover:
	go get -d && go test -v	-coverprofile=coverage.out
	go tool cover -html=coverage.out

hurrly:
	go build -o hurrly cmd/hurrly/main.go

# ==== packaging

deb: $(TARGETS)
	mkdir -p packaging/deb/hurrly/usr/sbin
	cp $(TARGETS) packaging/deb/hurrly/usr/sbin
	cd packaging/deb && fakeroot dpkg-deb --build hurrly .
	mv packaging/deb/hurrly_*deb .
