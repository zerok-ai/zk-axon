FROM golang:1.18-alpine
WORKDIR /zk
COPY zk-axon .
CMD ["/zk/zk-axon", "-c", "/zk/config/config.yaml"]

