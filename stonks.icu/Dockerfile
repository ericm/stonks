FROM golang:1.15-alpine AS build

WORKDIR /stonks

COPY . .
RUN go mod download

RUN go install github.com/ericm/stonks/stonks.icu

FROM alpine

RUN apk add tzdata

WORKDIR /bin

COPY --from=build /go/bin/stonks.icu ./stonks

RUN chmod +x stonks

CMD [ "stonks" ]