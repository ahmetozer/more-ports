FROM golang:1.16
WORKDIR $GOPATH/src/github.com/ahmetozer/more-ports/
COPY go.mod go.sum ./
RUN go mod download

COPY client ./client
COPY server ./server
COPY config ./config
COPY pkg ./pkg
COPY .git ./.git
COPY *.go ./
RUN export GIT_COMMIT=$(git rev-list -1 HEAD) && \
    export GIT_TAG=$(git tag | tail -1) && \
    export GIT_URL=$(git config --get remote.origin.url) && \
    CGO_ENABLED=0 go build -v -ldflags="-X 'main.GitUrl=$GIT_URL' -X 'main.GitTag=$GIT_TAG' -X 'main.GitCommit=$GIT_COMMIT' -X 'main.BuildTime=$(date -Isecond)' -X 'main.RunningEnv=container'" -o /app/more-ports

RUN export DEBIAN_FRONTEND=noninteractive && apt update && apt install -y libcap2-bin
RUN setcap CAP_NET_BIND_SERVICE=+eip /app/more-ports

RUN echo "nobody:x:65534:65534:Nobody:/:" > /app/passwd.minimal && \
mkdir -p /app/tmp/cert &&  chown nobody:nogroup -R /app/tmp


FROM scratch
USER nobody
COPY config /config
COPY --from=0  /app/more-ports /bin/more-ports
COPY --from=0  /app/passwd.minimal /etc/passwd
COPY --from=0  /app/tmp /tmp
LABEL org.opencontainers.image.source="https://github.com/ahmetozer/more-ports"
ENTRYPOINT [ "/bin/more-ports" ]