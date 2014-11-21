ottily: fmt
	go build -o ottily cmd/ottily/main.go

fmt:
	goimports -w .

clean:
	rm -f ottily