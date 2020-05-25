FROM golang:1.14 AS build

WORKDIR /stonks

COPY . .

RUN go get -v && go build -v -o stonks

FROM alpine

WORKDIR /bin

COPY --from=build /stonks/stonks stonks

RUN chmod +x stonks

CMD [ "stonks" ]