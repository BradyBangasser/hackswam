FROM golang:1.25 as builder

WORKDIR /app

COPY . ./
RUN go mod download
RUN python3 go-builder.py

RUN CGO_ENABLED=1 GOOS=linux go build -o server main.go

# runtime
FROM registry.fedoraproject.org/fedora-minimal:latest

WORKDIR /app
COPY --from=builder /app/server /app/server

# create uploads dir
RUN mkdir -p /app/uploads

EXPOSE 8080

CMD ["./server"]
