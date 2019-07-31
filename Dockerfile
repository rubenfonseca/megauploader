FROM golang:alpine AS builder

WORKDIR /app

COPY *.go ./

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o megauploader .

FROM scratch

COPY --from=builder /app/megauploader /app/megauploader

EXPOSE 9292

ENTRYPOINT ["/app/megauploader"]
