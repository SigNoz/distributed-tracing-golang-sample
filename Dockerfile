# use a minimal alpine image
FROM golang:alpine3.17

ARG TARGETOS=linux TARGETARCH=amd64

# add ca-certificates in case you need them
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# set working directory and copy source code
WORKDIR /go/src/github.com/NamanJain8/distributed-tracing-golang-sample
COPY . .

# build the binary and remove the source code
RUN apk add --no-cache bash && \
    go mod tidy && \
    go mod vendor && \
    go build -o /usr/local/bin/order ./order && \
    go build -o /usr/local/bin/users ./users && \
    go build -o /usr/local/bin/payment ./payment && \
    rm -rf /go/src/github.com/NamanJain8/distributed-tracing-golang-sample

# copy run.sh to execute the binaries in background
COPY run.sh /run.sh

# change working directory to root
WORKDIR /

# expose ports for distributed tracing golang sample
EXPOSE 8080 8081 8082

# run the binary as the entrypoint
ENTRYPOINT [ "/run.sh" ]
