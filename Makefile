all: fmt
	go build -o ottily cmd/ottily/main.go
	go build -o eachline cmd/eachline/eachline.go
	go build -o eachline-parse cmd/eachline/eachline-parse.go
	go build -o eachline-parse-parallel cmd/eachline/eachline-parse-parallel.go

fmt:
	goimports -w .

clean:
	rm -f ottily eachline eachline-parse eachline-parse-parallel

