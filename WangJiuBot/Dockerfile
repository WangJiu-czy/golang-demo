FROM golang:1.20 as builder
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o WangJiu .

FROM golang:1.20

WORKDIR /app
COPY --from=builder /app/WangJiu .
COPY --from=builder /app/develop .
RUN chmod +x /app/WangJiu

ENTRYPOINT ["/app/WangJiu"]