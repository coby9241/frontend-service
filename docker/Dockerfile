FROM golang:1.12 AS build-env

ENV GO111MODULE=on
ENV TZ=Asia/Singapore
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR $GOPATH/src/github.com/coby9241/frontend-service

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go run cmd/compile/main.go

RUN chmod +x scripts/build_binary.sh
RUN ./scripts/build_binary.sh

# thin Production image
FROM debian:stretch

ENV GOPATH /go
RUN mkdir -p "$GOPATH/src/github.com" && chmod -R 777 "$GOPATH"
WORKDIR /frontend-service

ENV TZ=Asia/Singapore
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Update root certicates, neccessary for communicating with s3
RUN apt-get update && apt-get install -y --no-install-recommends \
  ca-certificates

COPY --from=build-env /go/src/github.com/coby9241/frontend-service/templates ./templates
COPY --from=build-env /go/src/github.com/coby9241/frontend-service/frontend-service ./

RUN chmod +x ./frontend-service
CMD ["./frontend-service"]
