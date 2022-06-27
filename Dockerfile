FROM docker.io/library/golang:1.18
WORKDIR /src
COPY . .
RUN go build
ENTRYPOINT ["/src/save-something"]
