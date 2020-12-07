#golang:1.15.3-alpine3.12
FROM golang@sha256:6bdf3cc6cb0509e860e626a339c19b7e842efb8444960bb72d654e33607c34f5

WORKDIR /src

ADD go.mod go.sum ./

RUN go mod download

ADD ./file.go ./
ADD ./tcp ./tcp/
ADD ./xpath ./xpath/
ADD ./cli ./cli/
ADD ./cmd ./cmd/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-s -w" \
    -o ./out/afi-html-parser \
    ./cmd/afi-html-parser/main.go


#gcr.io/distroless/static:latest
FROM gcr.io/distroless/static@sha256:04c5f0473b7ebba65bfdd4125fe81fac1701783549f9b98fd96e3566f6035fa7

COPY --from=0 /src/out/afi-html-parser /usr/bin/

CMD ["afi-html-parser"]