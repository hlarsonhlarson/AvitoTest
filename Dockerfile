FROM golang:1.16.0-buster

RUN go version
ENV GOPATH=/

RUN mkdir ./app

COPY ./src ./app

WORKDIR ./app

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x wait-for-postgres.sh

RUN go mod download
RUN go build -o avito main.go
CMD ["./avito"]
