FROM golang:1.12 AS build
WORKDIR /go/src/github.com/alsx/wallet/
RUN go get -d -v golang.org/x/net/html
COPY . .
ENV GO111MODULE=on
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o wallet .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build /go/src/github.com/alsx/wallet/wallet .
EXPOSE 80:80
CMD ["./wallet"]
