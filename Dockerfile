FROM golang:1.19.2-alpine3.15 AS builder

WORKDIR /app

COPY . ./

RUN go mod download

RUN go build -o person-api .

FROM alpine:latest

WORKDIR /root

COPY --from=builder /app/person-api .

EXPOSE 8080
EXPOSE 5432

CMD ["./person-api"]