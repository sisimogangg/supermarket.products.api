FROM golang:latest AS build

WORKDIR  $GOPATH/src/github.com/sisimogangg/supermarket.products.api

RUN go get -u github.com/golang/dep/cmd/dep

COPY . .

RUN dep init

RUN dep ensure 

COPY config.json /usr/bin/config.json 

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix nocgo -o /usr/bin/products.api 

FROM alpine:3.9.4

COPY --from=build /usr/bin/products.api /root/
COPY --from=build /usr/bin/config.json /root/

WORKDIR /root/

CMD ["./products.api"]