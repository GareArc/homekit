# Dockerfile.dev
FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

# Base packages + repo setup helpers
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates curl git wget gnupg build-essential sudo \
    xz-utils unzip pkg-config \
 && rm -rf /var/lib/apt/lists/*

# ---------- GitHub CLI ----------
RUN (type -p wget >/dev/null || (apt update && apt install wget -y)) \
&& mkdir -p -m 755 /etc/apt/keyrings \
&& out=$(mktemp) && wget -nv -O$out https://cli.github.com/packages/githubcli-archive-keyring.gpg \
&& cat $out | tee /etc/apt/keyrings/githubcli-archive-keyring.gpg > /dev/null \
&& chmod go+r /etc/apt/keyrings/githubcli-archive-keyring.gpg \
&& mkdir -p -m 755 /etc/apt/sources.list.d \
&& echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list > /dev/null \
&& apt update \
&& apt install gh -y

# ---------- nvm + Node (LTS) ----------
ENV NVM_DIR=/usr/local/nvm
RUN mkdir -p "$NVM_DIR" \
 && curl -fsSL https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash \
 && printf 'export NVM_DIR=%s\n[ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh"\n[ -s "$NVM_DIR/bash_completion" ] && . "$NVM_DIR/bash_completion"\n' "$NVM_DIR" \
    > /etc/profile.d/nvm.sh \
 && bash -lc 'nvm install --lts && nvm alias default lts/* && nvm use default && \
     ln -sf "$NVM_DIR/versions/node/$(nvm version)/bin/node" /usr/local/bin/node && \
     ln -sf "$NVM_DIR/versions/node/$(nvm version)/bin/npm"  /usr/local/bin/npm  && \
     ln -sf "$NVM_DIR/versions/node/$(nvm version)/bin/npx"  /usr/local/bin/npx'

# ---------- uv (Python package/devenv manager) ----------
ENV UV_INSTALL_DIR=/usr/local
RUN curl -LsSf https://astral.sh/uv/install.sh | sh

# ---------- Go toolchain (configurable) ----------
ARG GO_VERSION=1.23.2
ENV GOPATH=/go
ENV PATH=/usr/local/go/bin:/go/bin:/usr/local/bin:$PATH
RUN mkdir -p "$GOPATH" \
 && curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -o /tmp/go.tgz \
 && rm -rf /usr/local/go && tar -C /usr/local -xzf /tmp/go.tgz && rm /tmp/go.tgz \
 && printf 'export GOPATH=%s\nexport PATH=/usr/local/go/bin:$GOPATH/bin:$PATH\n' "$GOPATH" > /etc/profile.d/go.sh

# ---------- workspace ----------
RUN mkdir -p /root/code
WORKDIR /root/code
