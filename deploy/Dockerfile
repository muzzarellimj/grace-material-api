# syntax=docker/dockerfile:1

FROM golang:1.21.6

# set working dir
WORKDIR /dist

# copy module manifest
COPY go.mod go.sum ./

# install Go dependencies
RUN go mod download

# copy source
COPY . .

# build source to binary
RUN CGO_ENABLED=0 GOOS=linux go build -o grace-material-api /dist/cmd/grace-material-api

# expose port to listen on
EXPOSE 8080

# run grace-material-api
CMD ["/dist/grace-material-api"]
