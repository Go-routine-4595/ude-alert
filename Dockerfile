FROM ubuntu:latest
LABEL authors="christophebuffard"
WORKDIR /UDEAlarms
COPY certs/server-ca.pem ./
COPY certs/client-cert.pem ./
COPY certs/cilent-key.pem ./
COPY conf.yml ./conf.yml
COPY udealarms-linux-x86 ./
CMD ["/UDEAlarms/udealarms-linux-x86"]

ENTRYPOINT ["top", "-b"]