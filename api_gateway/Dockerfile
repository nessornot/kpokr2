FROM golang:1.23.6

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o /api_gateway

COPY docs ./docs

CMD ["/api_gateway"]