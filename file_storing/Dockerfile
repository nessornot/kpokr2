FROM golang:1.23.6

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o /file_storing

CMD ["/file_storing"]