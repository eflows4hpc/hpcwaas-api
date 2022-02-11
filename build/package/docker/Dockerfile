FROM golang:1 as build-env
ARG BIN=hpcwaas-api

WORKDIR /go/src/app
COPY  . ./

RUN make ${BIN}

FROM gcr.io/distroless/static
ARG BIN=hpcwaas-api

COPY --from=build-env /go/src/app/bin/${BIN} /entrypoint

ENTRYPOINT [ "/entrypoint" ]
CMD ["help"]
