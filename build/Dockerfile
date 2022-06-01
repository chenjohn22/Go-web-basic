ARG RepositoryUri
FROM ${RepositoryUri}golang:1.14.2-alpine3.11 as builder
WORKDIR /app
COPY . .
RUN go build -o main

FROM ${RepositoryUri}alpine:3.12.0
ARG version
ENV version=$version
MAINTAINER justin
RUN apk add --no-cache tzdata
COPY --from=builder /app/build/entrypoint.sh /opt/entrypoint.sh
COPY --from=builder /app/config/env.json.example /opt/config/env.json.example
COPY --from=builder /app/build/gene-config.sh /opt/gene-config.sh
COPY --from=builder /app/main /opt
COPY --from=builder /app/sql /opt/sql
COPY --from=builder /app/web /opt/web
#COPY --from=builder /app/locales /opt/locales

WORKDIR /opt
ENTRYPOINT ["/opt/entrypoint.sh"]
CMD sh -c './gene-config.sh && ./main'