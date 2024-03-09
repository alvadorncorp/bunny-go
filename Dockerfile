FROM golang:1.21-alpine3.19 as builder

RUN mkdir -p /builder
RUN apk add make bash

WORKDIR /builder

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN make build

FROM alpine:3.19

WORKDIR /bin
COPY --from=builder /builder/build/cli /bin/bunny-cli

CMD ["/bin/bunny-cli"]