FROM golang:1.23 as builder
WORKDIR /app
COPY . .
RUN go build -o app ./cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd

FROM scratch
WORKDIR /app
COPY --from=builder /app/app .
CMD ["./app"]