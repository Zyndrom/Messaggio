FROM golang:alpine as builder

WORKDIR /build
ADD go.mod .

COPY . .

RUN go build -o messages_service cmd/main.go

FROM alpine:3

WORKDIR /build

COPY --from=builder /build/.env /build/.env
COPY --from=builder /build/migrations /build/migrations
COPY --from=builder /build/messages_service /build/messages_service

CMD ["./messages_service"]
