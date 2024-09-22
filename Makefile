default: \
	generate_mocks

install_mock_dependenies:
	go install github.com/golang/mock/mockgen@latest
	go get github.com/golang/mock/mockgen/model

generate_mocks_repository:
	mkdir -p ./pkg/authserv/repository/mocks
	mockgen -destination=pkg/authserv/repository/mocks/mock_repository.go -package=mocks github.com/SergeyIvanovDevelop/tss-tools/pkg/authserv/repository AuthRepository

generate_mocks:
	$(MAKE) install_mock_dependenies
	$(MAKE) generate_mocks_repository

.PHONY: \
	default \
	install_mock_dependenies \
	generate_mocks_repository \
	generate_mocks