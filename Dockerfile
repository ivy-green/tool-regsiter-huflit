FROM golang:1.19-alpine3.18 as builder
WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./
RUN go build -o tool cmd/main.go

FROM alpine:3.18
WORKDIR /app

COPY --from=builder app/tool .

CMD ["/app/tool"]

