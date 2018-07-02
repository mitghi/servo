all:
	go build .

clean:
	rm -f servo

test:
	go test -v .
