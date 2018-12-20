FROM golang:1.9
# create a working directory
WORKDIR /go/src/github.com/Benjar12/knock_challenge
COPY . .

RUN mkdir -p /knock/bin

# This is just here because compose is creating problems.
ENV PATH "$PATH:/knock/bin"

# Get all dependancies
RUN go get ./

# Build the gateway
RUN go build -o /knock/bin/main .