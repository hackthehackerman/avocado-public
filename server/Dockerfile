# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

ENV ENV production

ADD . /app
WORKDIR /app

RUN go mod download
RUN go build -o /avocado

EXPOSE 8080
CMD ["/avocado"]