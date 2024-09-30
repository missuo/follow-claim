FROM golang:1.23.1 AS builder
WORKDIR /go/src/github.com/missuo/follow-claim
COPY . .
RUN go get -d -v ./
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o follow-claim .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/src/github.com/missuo/follow-claim/follow-claim /app/follow-claim
CMD ["/app/follow-claim"]