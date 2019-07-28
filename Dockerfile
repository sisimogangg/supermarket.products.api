FROM golang:latest AS build
ARG stage

WORKDIR  $GOPATH/src/github.com/sisimogangg/supermarket.products.api

RUN go get -u github.com/golang/dep/cmd/dep

COPY Gopkg.toml Gopkg.lock ./

COPY . .

RUN dep ensure 

COPY config.json /usr/bin/config.json 

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix nocgo -o /usr/bin/server 

FROM alpine:3.9.4

COPY --from=build /usr/bin/server /root/
COPY --from=build /usr/bin/config.json /root/

EXPOSE 9090
WORKDIR /root/

CMD ["./server"]