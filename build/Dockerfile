FROM golang:1.23.5-alpine AS builder

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOPROXY=direct

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main cmd/app/main.go

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata bash wget

ENV TZ=America/Recife
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /workspace

RUN mkdir -p /workspace/logs /workspace/temp/arquivos/cbf801 && \
    chmod 777 /workspace/temp/arquivos/cbf801 && \
    chmod 777 /workspace/logs

COPY --from=builder /app/main /workspace/

COPY --from=builder /app/docs /workspace/docs/

RUN chown -R appuser:appgroup /workspace

USER appuser

EXPOSE 1450

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:1450/health || exit 1

CMD ["./main"]
