FROM golang:latest AS build

WORKDIR /go/src/github.com/sisimogangg/supermarket.products.api

COPY . .

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep  ensure -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o products-api


FROM alpine:3.9.4

RUN apk add ca-certificates

RUN mkdir /app
WORKDIR /app

COPY config.json .
COPY firebaseServiceAccount.json .
COPY --from=build /go/src/github.com/sisimogangg/supermarket.products.api/products-api .


CMD ["./products-api"]