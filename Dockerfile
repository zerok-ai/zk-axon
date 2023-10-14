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








+# The preceding star will make the file copy only if it exists
+COPY *bin/zk-wsp-client-arm64 /app/myapp-arm64

-# arm64 is optional architecture for the image.
-RUN if [ -f "bin/zk-wsp-client-arm64" ]; then \
-    COPY bin/zk-wsp-client-arm64 /app/myapp-arm64; \
-    RUN chmod +x /app/myapp-arm64; \
-fi
+RUN chmod +x /app/*


-# Run the Go executable
-CMD ["./wsp-client-start.sh"]
\ No newline at end of file
+# Run the start script
+CMD ["./app-start.sh","-amd64","myapp-amd64","-arm64","myapp-arm64","-c","/opt/wsp-config.yaml"]
\ No newline at end of file
