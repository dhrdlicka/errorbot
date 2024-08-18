FROM golang:alpine AS builder

WORKDIR /app

COPY ./ /app/

RUN go build .

FROM alpine

WORKDIR /app

COPY --from=builder /app/errorbot /app/errorbot
COPY ./json /app/json
COPY ./yaml /app/yaml

CMD [ "/app/errorbot" ]