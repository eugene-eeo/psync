test_all:
	go build
	go test github.com/eugene-eeo/psync/blockfs
	bats test
