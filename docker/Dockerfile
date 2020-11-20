FROM gesquive/go-builder:latest AS builder

ENV APP=reflect

COPY dist/ /dist/
RUN copy-release
RUN chmod +x /app/reflect

RUN mkdir -p /etc/reflect
COPY docker/config.yml /etc/reflect

# =============================================================================
FROM gesquive/docker-base:busybox
LABEL maintainer="Gus Esquivel <gesquive@gmail.com>"

# Import from builder
COPY --from=builder /app/reflect /app/
COPY --from=builder /etc/reflect/ /etc/reflect/

WORKDIR /config
VOLUME /config
EXPOSE 2626/tcp

ENTRYPOINT ["run"]
CMD ["/app/reflect"]
