test:
	go test -v

cover:
	rm -rf *.coverprofile
	go test -coverprofile=timeout.coverprofile
	gover
	go tool cover -html=timeout.coverprofile