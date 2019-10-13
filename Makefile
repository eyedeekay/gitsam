
build:
	cd gitsam && go build

try: build
	./gitsam/gitsam

fmt:
	find . -name '*.go' -exec gofmt -w {} \;

USER_GH=eyedeekay
packagename=gitsam
VERSION=0.0.2

tag:
	gothub release -s $(GITHUB_TOKEN) -u $(USER_GH) -r $(packagename) -t v$(VERSION) -d "I2P Tunnel Management tool for Go applications"