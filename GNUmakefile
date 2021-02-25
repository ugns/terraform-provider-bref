default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: major
major:
	git tag $$(svu major)
	git push --tags

.PHONY: minor
minor:
	git tag $$(svu minor)
	git push --tags

.PHONY: patch
patch:
	git tag $$(svu patch)
	git push --tags
