FROM golang:1.26-bookworm AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=0.1.0" -o tuipr ./cmd/tuipr/
RUN tar -czf tuipr_0.1.0_Linux_x86_64.tar.gz tuipr

FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive
ENV HOMEBREW_PREFIX=/home/brewuser/.linuxbrew
ENV HOMEBREW_CELLAR=/home/brewuser/.linuxbrew/Cellar
ENV HOMEBREW_REPOSITORY=/home/brewuser/.linuxbrew/Homebrew
ENV PATH="/home/brewuser/.linuxbrew/bin:/home/brewuser/.linuxbrew/sbin:${PATH}"

RUN apt-get update && apt-get install -y \
    curl \
    git \
    build-essential \
    procps \
    file \
    python3 \
    && rm -rf /var/lib/apt/lists/*

RUN useradd -m -s /bin/bash brewuser

USER brewuser
WORKDIR /home/brewuser

RUN git clone --depth=1 https://github.com/Homebrew/brew.git /home/brewuser/.linuxbrew/Homebrew \
    && mkdir -p /home/brewuser/.linuxbrew/bin \
    && ln -sf /home/brewuser/.linuxbrew/Homebrew/bin/brew /home/brewuser/.linuxbrew/bin/brew

COPY --from=builder /src/tuipr_0.1.0_Linux_x86_64.tar.gz /tmp/tuipr_0.1.0_Linux_x86_64.tar.gz

COPY --chown=brewuser:brewuser docker/tests/test-brew-install.sh /tmp/test-brew-install.sh
RUN chmod +x /tmp/test-brew-install.sh

ENTRYPOINT ["/bin/bash", "/tmp/test-brew-install.sh"]
