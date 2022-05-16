test_all:
	go test -v ./...

gen_mocks:
	mockgen -destination=pkg/core/sender_mock_test.go -package=core \
		github.com/meddion/pkg/core Sender
