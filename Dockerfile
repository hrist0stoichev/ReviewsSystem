FROM golang:1.14.3-alpine3.11 AS builder
WORKDIR /go/src/reviewssystem
COPY . .
RUN go get -d ./...
RUN GOOS=linux GO111MODULE=on CGO_ENABLED=0 go build -o /go/bin/reviewssystem


FROM node:14.4.0 AS frontend-builder
WORKDIR /go/src/reviewssystem
COPY front-end/ .
RUN npm install && npm run build


FROM alpine:3.11
EXPOSE 8001
COPY --from=builder /go/bin/reviewssystem /reviewssystem
COPY --from=builder /go/src/reviewssystem/db/migrations /db-migrations/
COPY --from=frontend-builder /go/src/reviewssystem/dist /static/
ENTRYPOINT ["/reviewssystem"]