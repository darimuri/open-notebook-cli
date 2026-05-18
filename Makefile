.PHONY: build install-local

build:
	go build -ldflags "-X github.com/darimuri/open-notebook-cli/cmd.version=$$(git describe --tags 2>/dev/null || echo dev)" -o open-notebook ./main.go

install-local: build
	cp open-notebook ~/.local/bin/open-notebook
	chmod +x ~/.local/bin/open-notebook