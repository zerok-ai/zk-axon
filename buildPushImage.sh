GOOS=linux GOARCH=amd64 go build -o zk-axon cmd/main.go
docker build -t zk-axon:dev .
sh ./gcp-artifact-deploy-go.sh