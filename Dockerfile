FROM golang:1.18.0

# hadolint ignore=DL3008
RUN apt-get update
RUN apt-get install -y --no-install-recommends xz-utils zip apt-transport-https curl gnupg2 ca-certificates unzip

WORKDIR /work
COPY go.mod /work
COPY go.sum /work
RUN go mod download
RUN go mod verify
COPY . /work
RUN go build -o main ./cmd/releaser-server/*.go
COPY main /releaser
# Create appuser
ENV USER=appuser
ENV UID=10001
RUN addgroup --gid ${UID} ${USER} && adduser --uid ${UID} --ingroup ${USER} ${USER}

# Change to the USER only after we run the build, since the build saves to root
USER appuser:appuser
EXPOSE 8080
ENTRYPOINT ["/releaser"]
