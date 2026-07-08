FROM golang:1.26-alpine

WORKDIR /app

COPY contracts/auth-service/ contracts/auth-service/

COPY services/auth-service ./

RUN go mod tidy

EXPOSE 8080

# CMD ["air"]