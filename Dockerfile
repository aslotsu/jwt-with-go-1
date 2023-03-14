FROM golang:1.20-alpine3.16

WORKDIR /app

COPY go.sum .

COPY go.mod .

COPY . .

RUN go mod download

EXPOSE 8000

RUN go build -o /myapp

CMD ["/myapp"]

