# REQUIRED SECTION
include ./.golang.mk
# END OF REQUIRED SECTION

# Run 'make help' for the list of default targets

# Example of overriding of default target
#test: ## run test with coverage using the vendor directory
#	go test -mod vendor -v -cover ./... -coverprofile cover.out

swagger-generate: ## generate swagger server and client code
	$(GOBIN)/swagger generate server --skip-main -t ./gen -f ./api/openapi.yml
	$(GOBIN)/swagger generate client -t ./gen -P abc -f ./api/openapi.yml

generate: ## run swagger-generate before running the default generate target
generate: swagger-generate generate-default


coverage: ## report on test coverage, default value reduced from 85 to 60%
coverage: test 
	goverreport -coverprofile=cover.out -sort=block -order=desc -threshold=60
