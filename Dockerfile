FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git protoc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
RUN export PATH="$PATH:$(go env GOPATH)/bin"
WORKDIR /go/src/app
COPY . .
RUN protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/links.proto
RUN go mod init golinkcut
RUN go mod tidy
RUN GOOS=linux go build -o ./bin/golinkcut ./cmd/

FROM alpine:3.17
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/src/app/bin /go/bin
EXPOSE 8080