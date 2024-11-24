FROM golang:1.23.3 AS builder

WORKDIR /app

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o dist .

FROM scratch

COPY --from=builder /app/dist /app

CMD ["/app"]
