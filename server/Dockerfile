FROM golang:1.22-alpine AS builder

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/gobaseproject ./cmd/gobaseproject

FROM alpine:3.20

WORKDIR /app

RUN addgroup -S app && adduser -S app -G app
COPY --from=builder /out/gobaseproject /app/gobaseproject

USER app
EXPOSE 8080

CMD ["/app/gobaseproject"]
