FROM golang:1.23.6

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o /file_analysis

CMD ["/file_analysis"]