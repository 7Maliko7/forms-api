FROM golang:1-alpine AS build

RUN apk add --no-cache git

WORKDIR /go/src/github.com/7Maliko7/forms-api

COPY ./ /go/src/github.com/7Maliko7/forms-api
RUN go build -o /bin/app /go/src/github.com/7Maliko7/forms-api/cmd/server

FROM alpine:3.17.3

EXPOSE 8080

COPY --from=build /bin/app /bin/app
COPY ./docker/config.yml /bin/config.yml

ENTRYPOINT ["/bin/app", "-c", "/bin/config.yml"]