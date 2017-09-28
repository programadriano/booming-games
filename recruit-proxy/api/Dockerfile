FROM golang:latest

# Copy the local package files to the container’s workspace.
ADD . /go/src/github.com/brunomilani/booming-games/recruit-proxy/api

# Install our dependencies       
RUN go get github.com/satori/go.uuid
RUN go get github.com/streadway/amqp 

# Install api binary globally within container
RUN go install github.com/brunomilani/booming-games/recruit-proxy/api

# Set binary as entrypoint
ENTRYPOINT /go/bin/api

# Expose default port (8020)
EXPOSE 8020 
