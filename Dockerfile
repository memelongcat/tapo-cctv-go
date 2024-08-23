FROM golang:1.23.0-alpine3.20

WORKDIR /app

RUN mkdir -p /app/records

RUN apk add --no-cache ffmpeg

COPY . .

RUN go build main.go

CMD ["./main"]