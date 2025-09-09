FROM golang:1.22.5-alpine3.19 AS build

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV WORKSPACE=${GOPATH}/src/app

WORKDIR ${WORKSPACE}

RUN apk add --no-cache build-base gcc git pkgconf musl-dev openssh --virtual .build-deps

RUN mkdir /root/.ssh/ && \
    chmod 0700 /root/.ssh && \
    ssh-keyscan -p 2222 gitlab.servicenet.com.br > /root/.ssh/known_hosts

COPY id_rsa /root/.ssh/
RUN chmod 600 /root/.ssh/id_rsa

RUN git config --global url."ssh://git@gitlab.servicenet.com.br:2222".insteadOf "http://gitlab.servicenet.com.br"

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -ldflags '-extldflags "-static"' -o bin/app cmd/app/main.go

RUN rm -R /root/.ssh/

FROM alpine:3.12.1

ENV TZ=America/Recife
ENV WORKSPACE=/workspace

RUN apk add --no-cache --virtual .build-deps ca-certificates tzdata bash \
    && cp /usr/share/zoneinfo/${TZ} /etc/localtime && echo ${TZ} > /etc/timezone

WORKDIR ${WORKSPACE}

RUN mkdir -p ${WORKSPACE}/temp/arquivos/cbf801 \
    && chmod 777 ${WORKSPACE}/temp/arquivos/cbf801

COPY --from=build /go/src/app/api/ ${WORKSPACE}/api/
COPY --from=build /go/src/app/bin ${WORKSPACE}

CMD ["./app"]
