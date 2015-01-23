all: fmt
	go build -o ottily cmd/ottily/main.go
	go build -o eachline cmd/eachline/eachline.go
	go build -o eachline-parse cmd/eachline/eachline-parse.go

fmt:
	goimports -w .

clean:
	rm -f ottily eachline eachline-parse

