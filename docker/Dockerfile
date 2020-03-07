FROM index.docker.io/gesquive/go-builder:latest AS builder

ENV APP=reflect

COPY dist/ /dist/
RUN copy-release

RUN mkdir -p /etc/${APP}
COPY docker/config.yml /etc/${APP}

# =============================================================================
FROM gesquive/docker-base:latest
LABEL maintainer="Gus Esquivel <gesquive@gmail.com>"

# Import from builder
COPY --from=builder /app/${APP} /app/
COPY --from=builder /etc/${APP} /etc/

WORKDIR /config
VOLUME /config
EXPOSE 2626/tcp

ENTRYPOINT ["run", "/app/reflect"]
