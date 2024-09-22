default: \
	generate_mocks

generate_mocks:
	go install github.com/golang/mock/mockgen@latest
	go get github.com/golang/mock/mockgen/model
	mkdir -p ./pkg/authserv/repository/mocks
	mockgen -destination=pkg/authserv/repository/mocks/mock_repository.go -package=mocks github.com/SergeyIvanovDevelop/tss-tools/pkg/authserv/repository AuthRepository

.PHONY: \
	default \
	generate_mocks