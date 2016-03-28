target=release

$(target)/itpkg:
	go build -ldflags "-s" -o $@ main.go

clean:
	-rm -r $(target)

format:
	for f in `find . -type f -iname '*.go'`; do gofmt -w $$f; done
