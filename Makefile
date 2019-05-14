TARGET=necrobrowser
PACKAGES=action core log server zombie

all: deps build

deps: godep golint gofmt gomegacheck updatedeps

build:
	@go build -a -o $(TARGET) .

clean:
	@rm -rf $(TARGET)
	@rm -rf build
	@dep prune

updatedeps:
	@dep ensure -update -v
	@dep prune
	@git add "Gopkg.*" "vendor"
	@git commit -m "Updated deps :star2: (via Makefile)"

# tools
godep:
	@go get -u github.com/golang/dep/...

golint:
	@go get -u golang.org/x/lint/golint

gomegacheck:
	@go get honnef.co/go/tools/cmd/megacheck

gofmt:
	gofmt -s -w $(PACKAGES)
