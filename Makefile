vet:
	./deployment/vet.sh
lint:
	./deployment/lint.sh
test:
	./deployment/test.sh
sure: lint vet test
