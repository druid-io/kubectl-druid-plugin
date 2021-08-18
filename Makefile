.PHONY: build
build:
	@docker build -t kubectl-druid .
	@docker run -v ${PWD}/build:/build kubectl-druid cp -r /go/src/pkg/ /build
