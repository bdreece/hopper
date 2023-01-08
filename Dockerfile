FROM golang:alpine

WORKDIR /usr/src/hopper

ENV HOPPER_HOSTNAME="127.0.0.1"
ENV HOPPER_USERNAME="root"
ENV HOPPER_PASSWORD="password123"
ENV HOPPER_SECRET="mysecret12345"
ENV PROJECT_ROOT="/usr/src/hopper"

RUN apk --no-cache --update add gcc libc-dev git protobuf protobuf-dev
RUN go install github.com/golang/protobuf/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY go.mod go.sum ./
RUN go mod download -x && \
    go mod verify

COPY . .
RUN go generate -v -x ./...
RUN go build -v -o /usr/local/bin ./...

CMD [ "hopper" ]