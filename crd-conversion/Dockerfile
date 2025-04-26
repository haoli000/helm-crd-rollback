FROM golang:1.24.1 as builder

WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o webhook ./cmd/webhook/main.go

FROM alpine:3.17
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /workspace/webhook /webhook
USER 65532:65532

ENTRYPOINT ["/webhook"]