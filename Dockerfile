FROM golang:1.18 as build
WORKDIR /src
COPY go.mod go.sum /src/
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 make immugw-static
RUN mkdir /empty

FROM gcr.io/distroless/base:nonroot
LABEL org.opencontainers.image.authors="Codenotary Inc. <info@codenotary.com>"

WORKDIR /usr/sbin
COPY --from=build /src/immugw /usr/sbin/immugw

ENV IMMUGW_DIR="/var/lib/immudb" \
    IMMUGW_ADDRESS="0.0.0.0" \
    IMMUGW_PORT="3323" \
    IMMUGW_IMMUDB_ADDRESS="127.0.0.1" \
    IMMUGW_IMMUDB_PORT="3322" \
    IMMUGW_MTLS="false" \
    IMMUGW_DETACHED="false" \
    IMMUGW_PKEY="" \
    IMMUGW_CERTIFICATE="" \
    IMMUGW_CLIENTCAS="" \
    IMMUGW_AUDIT="false" \
    IMMUGW_AUDIT_USERNAME="" \
    IMMUGW_AUDIT_PASSWORD=""

COPY --from=build --chown=nonroot:nonroot /empty "$IMMUGW_DIR"

EXPOSE 3323
EXPOSE 9476

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "/usr/sbin/immugw", "version" ]
ENTRYPOINT ["/usr/sbin/immugw"]
