FROM golang:1.17-buster AS build

ENV GOPATH=/
WORKDIR /src/
COPY ./ /src/

RUN go mod download; CGO_ENABLED=0 go build -o /exam-app ./cmd/main.go

FROM alpine:latest

COPY --from=build /exam-app /exam-app
COPY ./config/config.json /config/
COPY ./static/ /static/
COPY ./wait-for-postgres.sh ./

RUN apk --no-cache add postgresql-client util-linux && chmod +x ./wait-for-postgres.sh && mkdir /log

CMD ["./exam-app"]