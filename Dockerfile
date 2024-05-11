FROM ubuntu:latest
LABEL authors="christophebuffard"
WORKDIR /UDEAlarms
COPY cert/server-ca.pem ./
COPY cert/client-cert.pem ./
COPY cert/client-key.pem ./
COPY config.yaml ./config.yaml
COPY udealarms-linux-x86 ./
CMD ["/UDEAlarms/udealarms-linux-x86"]

