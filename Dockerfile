FROM golang:1.22.0-alpine AS builder

WORKDIR /go/src/mystprom

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o mystprom /go/src/mystprom/cmd/main.go

FROM alpine:3.19

COPY --from=builder /go/src/mystprom/mystprom /usr/bin/mystprom

EXPOSE 9300

ENTRYPOINT ["/usr/bin/mystprom"]