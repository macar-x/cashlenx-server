FROM mcr.microsoft.com/devcontainers/base:ubuntu-24.04

# Install Go
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends wget git curl \
    && wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz \
    && rm go1.21.5.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/home/vscode/go"
ENV PATH="${GOPATH}/bin:${PATH}"

# Install Flutter
RUN git clone https://github.com/flutter/flutter.git -b stable /opt/flutter \
    && /opt/flutter/bin/flutter precache \
    && /opt/flutter/bin/flutter config --enable-web

ENV PATH="/opt/flutter/bin:${PATH}"

# Install Android SDK dependencies (for future mobile builds)
RUN apt-get -y install --no-install-recommends \
    openjdk-17-jdk \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
