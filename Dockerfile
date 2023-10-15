FROM alpine:latest
WORKDIR /zk

# base name of the executable
ENV exeBaseName="zk-axon"

# full path to the all the executables
ENV exeAMD64="bin/${exeBaseName}-amd64"
ENV exeARM64="bin/${exeBaseName}-arm64"

# copy the executables
COPY "$exeAMD64" .
COPY "$exeARM64" .

# call the start script
CMD ["./app-start.sh","--amd64","$exeAMD64","--arm64","$exeARM64", "-c", "config/config.yaml"]
