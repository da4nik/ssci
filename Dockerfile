FROM golang:1.8.3 as build

COPY . /usr/local/go/src/github.com/da4nik/ssci
WORKDIR /usr/local/go/src/github.com/da4nik/ssci

RUN mkdir /app && \
    curl https://glide.sh/get | sh && \
    glide install && \
    make build




FROM alpine:3.5

# Uncomment if project uses ssl
RUN apk add --no-cache git make ca-certificates

COPY --from=build /usr/local/go/src/github.com/da4nik/ssci/ssci /app/
WORKDIR /app
CMD ["/app/ssci"]
