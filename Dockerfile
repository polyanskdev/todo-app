FROM golang:1.17.1-buster

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-postgres.sh executable
RUN chmod +x wait-postgres.sh

# build app
RUN go mod download
RUN go build -o todo-app ./cmd/main.go

CMD ["./todo-app"]