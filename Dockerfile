FROM golang:1.18.0

# hadolint ignore=DL3008
RUN apt-get update
RUN apt-get install -y --no-install-recommends xz-utils zip apt-transport-https curl gnupg2 ca-certificates unzip

# Install mage
ARG MAGE_VERSION=1.13.0
RUN curl -L -o /tmp/mage.tar.gz "https://github.com/magefile/mage/releases/download/v${MAGE_VERSION}/mage_${MAGE_VERSION}_Linux-64bit.tar.gz" && tar -C /tmp -zxvf /tmp/mage.tar.gz && mv /tmp/mage /usr/local/bin

WORKDIR /work
COPY go.mod /work
COPY go.sum /work
RUN go mod download
RUN go mod verify
COPY . /work
RUN go build -o main ./cmd/releaser-server/*.go
#RUN GOBUILD_MAIN_DIRECTORY=./cmd/releaser-server mage go:build
# Create appuser
ENV USER=appuser
ENV UID=10001
RUN addgroup --gid ${UID} ${USER} && adduser --uid ${UID} --ingroup ${USER} ${USER}
# Use an unprivileged user.
COPY main /releaser

# Change to the USER only after we run the build, since the build saves to root
USER appuser:appuser
EXPOSE 8080
ENTRYPOINT ["/releaser"]