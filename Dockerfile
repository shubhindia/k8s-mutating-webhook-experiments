FROM --platform=linux/amd64 golang:1.20
WORKDIR /opt
COPY ./bin/ .
CMD ["./manager", "--tls-cert", "/etc/opt/tls.crt", "--tls-key", "/etc/opt/tls.key"]