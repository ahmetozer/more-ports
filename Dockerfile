From golang
WORKDIR $GOPATH/src/github.com/ahmetozer/more-ports/
COPY . .
RUN go get -d -v .
RUN export GIT_COMMIT=$(git rev-list -1 HEAD) && \
    export GIT_TAG=$(git tag | tail -1) && \
    export GIT_URL=$(git config --get remote.origin.url) && \
    CGO_ENABLED=0 go build -v -ldflags="-X 'main.GitUrl=$GIT_URL' -X 'main.GitTag=$GIT_TAG' -X 'main.GitCommit=$GIT_COMMIT' -X 'main.BuildTime=$(date -Isecond)' -X 'main.RunningEnv=container'" -o /app/more-ports
FROM scratch
COPY --from=0  /app/more-ports /bin/more-ports
LABEL org.opencontainers.image.source="https://github.com/ahmetozer/more-ports"
ENTRYPOINT [ "/bin/more-ports" ]