FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY *.go ./
RUN go build -ldflags="-s -w" -o /pwgen .

FROM alpine:3.20
COPY --from=build /pwgen /usr/local/bin/pwgen
ENTRYPOINT ["pwgen"]
