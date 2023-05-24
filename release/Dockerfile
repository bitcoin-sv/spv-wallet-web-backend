# syntax=docker/dockerfile:1

FROM --platform=$TARGETPLATFORM alpine:3.18

WORKDIR /app

ENV DB_SCHEMA_PATH=./migrations
COPY ./data/sql/migrations $DB_SCHEMA_PATH
COPY ./release/entrypoint.sh .
ARG project_name
COPY ./$project_name .

ARG version
ENV APP_BINARY=${project_name}
ENV VERSION=${version:-develop}
ENV TAG=${tag:-master}

ENTRYPOINT ["./entrypoint.sh"]
