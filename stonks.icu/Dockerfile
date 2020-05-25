FROM golang:1.14-alpine AS build

WORKDIR /stonks

COPY . .
RUN go mod download

RUN go install github.com/ericm/stonks/stonks.icu

FROM alpine

WORKDIR /bin

COPY --from=build /go/bin/stonks.icu ./stonks

RUN chmod +x stonks

CMD [ "stonks" ]