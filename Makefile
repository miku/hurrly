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
	rm -f debian/hurrly*.deb
	rm -rf debian/hurrly/usr

cover:
	go get -d && go test -v	-coverprofile=coverage.out
	go tool cover -html=coverage.out

hurrly:
	go build cmd/hurrly/main.go

# ==== packaging

deb: $(TARGETS)
	mkdir -p debian/hurrly/usr/sbin
	cp $(TARGETS) debian/hurrly/usr/sbin
	cd debian && fakeroot dpkg-deb --build hurrly .
