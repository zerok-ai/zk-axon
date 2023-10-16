#R: Move this to Makefile. We can keep all build and deploy scripts at one place.

GOOS=linux GOARCH=amd64 go build -o zk-axon cmd/main.go
docker build -t zk-axon:dev .
sh ./gcp-artifact-deploy-go.sh