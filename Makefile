HOSTNAME=edgedelta.com
NAMESPACE=local
NAME=edgedelta
BINARY=terraform-provider-${NAME}
VERSION=${TERRAFORM_PROVIDER_ED_VERSION}
OS_ARCH=${TERRAFORM_PROVIDER_ED_OS_ARCH}

default: install

build:
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${VERSION}/${OS_ARCH}
	cp ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ${GOPATH}/bin

release:
	goreleaser release --rm-dist