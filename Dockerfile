FROM golang:1.20
WORKDIR /opt
COPY . .
RUN go mod tidy
RUN go build main.go
CMD ["./main", "--tls-cert", "/etc/opt/tls.crt", "--tls-key", "/etc/opt/tls.key"]