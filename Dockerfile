FROM gcr.io/distroless/base-debian12:debug

ARG TARGETOS
ARG TARGETARCH

COPY nexus-$TARGETOS-$TARGETARCH /usr/local/bin/nexus

EXPOSE 8090
EXPOSE 8091

ENTRYPOINT ["/usr/local/bin/nexus", "serve", "--http=0.0.0.0:8090"]

