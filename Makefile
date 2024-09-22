default: \
	generate_mocks

install_mock_dependenies:
	go install github.com/golang/mock/mockgen@latest
	go get github.com/golang/mock/mockgen/model

generate_mocks_repository:
	mkdir -p ./pkg/authserv/repository/mocks
	mockgen -destination=pkg/authserv/repository/mocks/mock_repository.go -package=mocks github.com/SergeyIvanovDevelop/tss-tools/pkg/authserv/repository AuthRepository

generate_mocks_authenticator:
	mkdir -p ./pkg/authserv/auth/mocks
	mockgen -destination=pkg/authserv/auth/mocks/mock_auth.go -package=mocks github.com/SergeyIvanovDevelop/tss-tools/pkg/authserv/auth Authenticator

generate_mocks:
	$(MAKE) install_mock_dependenies
	$(MAKE) generate_mocks_repository
	$(MAKE) generate_mocks_authenticator

.PHONY: \
	default \
	install_mock_dependenies \
	generate_mocks_repository \
	generate_mocks_authenticator \
	generate_mocks