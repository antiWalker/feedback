FROM golang:1.15.7-alpine3.13
WORKDIR /app
COPY . .
RUN apk update && apk add git
RUN go mod download
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]