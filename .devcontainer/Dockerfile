ARG VARIANT="1.20"
FROM --platform=linux/amd64 mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

ENV MONGO_BIN=mongod

ARG MONGO_VERSION="4.2"
RUN curl -sSL "https://www.mongodb.org/static/pgp/server-${MONGO_VERSION}.asc" | (OUT=$(apt-key add - 2>&1) || echo $OUT) && \
    echo "deb http://repo.mongodb.org/apt/debian buster/mongodb-org/${MONGO_VERSION} main" | tee /etc/apt/sources.list.d/mongodb-org-${MONGO_VERSION}.list && \
    apt-get update && export DEBIAN_FRONTEND=noninteractive && \
    apt-get install -y \
    mongodb-org \
    mongodb-org-tools \
    mongodb-org-shell && \
    apt-get clean -y && rm -rf /var/lib/apt/lists/*
