FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY . .

RUN go mod download

RUN go build -o roomsbooking cmd/roomsbooking/main.go

FROM alpine:3.23 AS final

WORKDIR /app

COPY --from=builder /build/roomsbooking .

COPY --from=builder /build/.env .

CMD [ "/app/roomsbooking" ]