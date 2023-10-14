FROM golang:1.18-alpine

WORKDIR /zk

# amd64 is default architecture for the image.
RUN echo "Copying Amd64 file."
COPY bin/zk-axon-amd64 /zk/zk-axon-amd64

COPY app-start.sh /zk/app-start.sh

# The preceding star will make the file copy only if it exists
COPY *bin/zk-axon-arm64 /zk/zk-axon-arm64

RUN chmod +x /zk/*

# Run the start script
#CMD ["/zk/zk-axon", "-c", "/zk/config/config.yaml"]
CMD ["./app-start.sh","-amd64","zk-axon-amd64","-arm64","zk-axon-arm64","-c","/zk/config/config.yaml"]
