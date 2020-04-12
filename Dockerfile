FROM golang:1.13-alpine as build-env

ADD . /kbcdatastudioproc

WORKDIR /kbcdatastudioproc

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/kbcdatastudioproc cmd/main.go

FROM scratch

COPY --from=build-env /go/bin/kbcdatastudioproc /kbcdatastudioproc

ENTRYPOINT ["/kbcdatastudioproc"]
