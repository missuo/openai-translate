FROM golang:1.21 AS builder
WORKDIR /go/src/github.com/missuo/openai-translate
COPY main.go ./
COPY go.mod ./
COPY go.sum ./
RUN go get -d -v ./
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o openai-translate .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/src/github.com/missuo/openai-translate/openai-translate /app/openai-translate
CMD ["/app/openai-translate"]