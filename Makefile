.PHONY: build
build:
	@docker build -t gcr.io/rilldata/kubectl-druid .
	@docker run -v ${PWD}/build:/build -it gcr.io/rilldata/kubectl-druid cp -r /go/src/pkg/ /build
