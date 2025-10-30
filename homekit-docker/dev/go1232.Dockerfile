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
