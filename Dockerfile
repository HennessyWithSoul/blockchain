FROM golang:1.18
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /build
COPY go.* ./
COPY . .
RUN go build -o /main .
RUN chmod +x main
EXPOSE 8081
ENTRYPOINT ["/main"]
