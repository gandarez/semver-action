FROM golang:1.21-alpine

RUN apk add --update --no-cache \
    make \
    git \
    curl

WORKDIR /go/src/github.com/gandarez/semver-action

COPY . .

# build
RUN make build-linux

# apply permissions
RUN chmod a+x ./build/linux/amd64/semver

# symbolic link
RUN ln -s /go/src/github.com/gandarez/semver-action/build/linux/amd64/semver /bin/

# Specify the container's entrypoint as the action
ENTRYPOINT ["/bin/semver"]
