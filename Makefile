override DOCKER                         = $(shell which docker)
override GIT_VERSION                    = $(shell git rev-parse --abbrev-ref HEAD)${CUSTOM} $(shell git rev-parse HEAD)
override PROJECT_NAME                   = wangjiu #编译后的文件名
override LDFLAGS                                = -ldflags "-X 'main.version=\"${GIT_VERSION}\"'"
override GOBIN                                  = ${shell pwd}/bin

GO_COMPILER_IMAGE ?= golang:1.20


PROJECT_VERSION = $(shell if [ "$$(git tag --points-at HEAD | tail -n1)" ]; then git tag --points-at HEAD | tail -n1 | sed 's/v\(.*\)/\1/'; else git rev-parse --abbrev-ref HEAD | sed 's/release-\(.*\)/\1/' | tr '-' '\n' | head -n1; fi)
#
default: install
#
install:
        go build  ${LDFLAGS} -o $(GOBIN)/$(PROJECT_NAME) ./

docker_install:
        $(DOCKER) run -e GOPROXY=https://goproxy.cn,direct -v $(shell pwd):/universe --rm $(GO_COMPILER_IMAGE) sh -c "cd /universe && make install $(MAKEFLAGS)"
