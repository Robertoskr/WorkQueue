# syntax=docker/dockerfile:1

# Alpine is chosen for its small footprint
# compared to Ubuntu
FROM golang:1.16-alpine
WORKDIR /app

COPY main.go ./

RUN go build main.go

EXPOSE 8000

CMD ["./main"]