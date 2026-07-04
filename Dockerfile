FROM golang:1.26-alpine

RUN apk add --update --no-cache \
    make \
    git \
    curl

WORKDIR /go/src/github.com/gandarez/semver-action

COPY . .

# build for the container's native architecture (amd64 or arm64)
RUN make build-linux-native

# apply permissions
RUN chmod a+x ./build/linux/semver

# symbolic link
RUN ln -s /go/src/github.com/gandarez/semver-action/build/linux/semver /bin/

# Specify the container's entrypoint as the action
ENTRYPOINT ["/bin/semver"]
