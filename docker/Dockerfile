FROM alpine:3.22 AS certs
RUN apk add --no-cache ca-certificates

FROM scratch
LABEL org.opencontainers.image.source="https://github.com/traggo/server"
WORKDIR /opt/traggo
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY traggo /opt/traggo
EXPOSE 3030
ENTRYPOINT ["./traggo"]
