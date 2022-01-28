FROM golang:1 as build-env

WORKDIR /go/src/app
COPY  . ./

RUN make

FROM gcr.io/distroless/static

COPY --from=build-env /go/src/app/bin/hpcwaas-api /
CMD ["/hpcwaas-api"]
