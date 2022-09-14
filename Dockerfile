FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /proxer

EXPOSE 8080

CMD [ "/proxer" ]