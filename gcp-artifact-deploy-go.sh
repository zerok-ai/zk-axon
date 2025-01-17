#R: Move this to Makefile. We can keep all build and deploy scripts at one place.

LOCATION="us-west1"
PROJECT_ID="zerok-dev"
REPOSITORY="zk-axon"
IMAGE="zk-axon"
TAG="dev"
ART_Repo_URI="$LOCATION-docker.pkg.dev/$PROJECT_ID/$REPOSITORY/$IMAGE:$TAG"

docker build -t $ART_Repo_URI .
#docker tag $IMAGE $ART_Repo_URI

gcloud auth configure-docker \
    $LOCATION-docker.pkg.dev

docker push $ART_Repo_URI