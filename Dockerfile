# Minimal runtime with certs
FROM gcr.io/distroless/static:nonroot
COPY --from=alpine:3.20 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY mcdowell /mcdowell
USER nonroot:nonroot

ENTRYPOINT ["/mcdowell"]
