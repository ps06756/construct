FROM ubuntu:jammy

ENV TRIGGER_REBUILD=1

RUN apt-get update && apt-get install -y \
    zip unzip upx \
    bash-completion \
    build-essential \
    htop \
    iputils-ping \
    jq \
    less \
    locales \
    nano \
    netcat-traditional \
    ripgrep \
    sudo \
    time \
    multitail \
    lsof \
    ssl-cert \
    python3-pip \
    shellcheck \
    curl \
    zsh \
    gnupg2 \
    openssh-client \
    file zstd \
    dumb-init \
    procps \
    pigz \
    && rm -rf /var/lib/apt/lists/* \
    && locale-gen en_US.UTF-8

ENV LANG=en_US.UTF-8

### Git ###
ADD https://raw.githubusercontent.com/gitpod-io/workspace-images/main/base/default.gitconfig /etc/gitconfig


USER root
ENV HOME=/root
WORKDIR $HOME

RUN install -m 0755 -d /etc/apt/keyrings && \
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg && \
    chmod a+r /etc/apt/keyrings/docker.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu jammy stable" >/etc/apt/sources.list.d/docker.list

RUN apt-get update && apt-get install -y \
    docker-ce=5:27.5.1-1~ubuntu.22.04~jammy \
    docker-ce-cli=5:27.5.1-1~ubuntu.22.04~jammy \
    docker-buildx-plugin=0.20.0-1~ubuntu.22.04~jammy \
    docker-compose-plugin=2.32.4-1~ubuntu.22.04~jammy \
    && rm -rf /var/lib/apt/lists/*

ENV GO_VERSION=1.24.1
ENV GOPATH=$HOME/go-packages
ENV GOROOT=$HOME/go
ENV PATH=$GOROOT/bin:$GOPATH/bin:$PATH

ARG TARGETPLATFORM
# Determine GO_PLATFORM based on TARGETPLATFORM and write it to a file
RUN case "${TARGETPLATFORM}" in \
    "linux/arm64") \
    echo "GO_PLATFORM=linux-arm64" > /tmp/go_platform.env \
    ;; \
    *) \
    echo "GO_PLATFORM=linux-amd64" > /tmp/go_platform.env \
    ;; \
    esac

RUN . /tmp/go_platform.env && \
    curl -fsSL https://dl.google.com/go/go$GO_VERSION.$GO_PLATFORM.tar.gz | tar xzs
RUN go install -v github.com/uudashr/gopkgs/cmd/gopkgs@v2 \
    && go install -v github.com/ramya-rao-a/go-outline@latest \
    && go install -v github.com/cweill/gotests/gotests@latest \
    && go install -v github.com/fatih/gomodifytags@latest \
    && go install -v github.com/josharian/impl@latest \
    && go install -v github.com/haya14busa/goplay/cmd/goplay@latest \
    && go install -v github.com/go-delve/delve/cmd/dlv@latest \
    && go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@latest \
    && go install -v golang.org/x/tools/gopls@latest \
    && go install -v honnef.co/go/tools/cmd/staticcheck@latest \
    && rm -rf $HOME/.cache $HOME/go-packages/pkg

