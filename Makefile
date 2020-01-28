GO ?= go
TESTFOLDER := $(shell $(GO) list ./...)
TESTTAGS ?= "integration"

.PHONY: test
test:
	$(GO) test ./... -race

test-integration:
	docker-compose build && docker-compose run web $(GO) test ./... -race -tags $(TESTTAGS) -cover

test-ci:
	echo "mode: count" > coverage.out
	for d in $(TESTFOLDER); do \
		$(GO) test -tags $(TESTTAGS) -v -covermode=count -coverprofile=profile.out $$d > tmp.out; \
		cat tmp.out; \
		if grep -q "^--- FAIL" tmp.out; then \
			rm tmp.out; \
			exit 1; \
		elif grep -q "build failed" tmp.out; then \
			rm tmp.out; \
			exit 1; \
		elif grep -q "setup failed" tmp.out; then \
			rm tmp.out; \
			exit 1; \
		fi; \
		if [ -f profile.out ]; then \
			cat profile.out | grep -v "mode:" >> coverage.out; \
			rm profile.out; \
		fi; \
	done

.PHONY: vet
vet:
	$(GO) vet $(VETPACKAGES)
