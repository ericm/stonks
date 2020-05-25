FROM golang:1.14 AS build

WORKDIR /stonks

COPY . .

RUN go get -v && go build -v -o stonks

FROM alpine

WORKDIR /stonks

COPY --from=build stonks stonks

ENTRYPOINT [ "/stonks/stonks" ]