FROM golang

WORKDIR /app

COPY . .

RUN go build main.go

CMD ["./main"]