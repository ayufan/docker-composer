ARG ARCH=
FROM ${ARCH}golang:alpine as build
COPY . $GOPATH/src/github.com/ayufan/docker-composer/
RUN cd $GOPATH/src/github.com/ayufan/docker-composer/ && \
  go mod init && \
  go mod tidy && \
  go install -v ./...

FROM ${ARCH}alpine:latest

RUN apk add -U git docker bash docker-compose && \
  git config --global receive.denyCurrentBranch updateInstead && \
  git config --global user.name Composer && \
  git config --global user.email you@example.com

ENV APPS_DIR=/srv/apps \
  GOPATH=/go

COPY --from=0 /go/bin/docker-composer /usr/bin/composer
VOLUME ["/srv/apps"]
ADD examples/ /srv/apps/

ENTRYPOINT ["/usr/bin/composer"]
