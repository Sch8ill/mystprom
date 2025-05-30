FROM golang:1.24.3-alpine AS builder

WORKDIR /go/src/mystprom

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s" -trimpath -o mystprom /go/src/mystprom/cmd/main.go

FROM alpine:3.21

COPY --from=builder /go/src/mystprom/mystprom /usr/bin/mystprom

EXPOSE 9300

ENTRYPOINT ["/usr/bin/mystprom"]