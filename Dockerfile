FROM golang:alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/main.go

FROM alpine:latest
COPY --from=builder /app/server /app/server

WORKDIR /app
ENV SERVER_PORT=8001
EXPOSE 8001
CMD [ "/app/server" ]
