# syntax=docker/dockerfile:1.4
FROM golang:1.25.3-trixie

# shell settings
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
ENV SHELL=bash

# workspace settings
ARG WORKING_DIR=/workspaces/app

# set `DEVCONTAINER` environment variable to help with orientation
ENV DEVCONTAINER=true

# container user settings
ARG CONTAINER_USER=user
ARG UID=1000
ARG GID=1000
ENV XDG_LOCAL="/home/${CONTAINER_USER}/.local"

# mise settings
ENV MISE_DATA_DIR="/mise"
ENV MISE_CONFIG_DIR="/mise"
ENV MISE_CACHE_DIR="/mise/cache"
ENV MISE_INSTALL_PATH="/usr/local/bin/mise"
ENV PATH="/mise/shims:$PATH"
# ENV MISE_VERSION="..."

# pnpm settings
ENV PNPM_HOME="/home/${CONTAINER_USER}/.local/share/pnpm"

# claude code settings
ARG TZ=
ENV TZ="${TZ}"
ENV CLAUDE_CODE_CONFIG="/home/${CONTAINER_USER}/.claude"

# Copy scripts
COPY start.sh /usr/local/bin/
COPY post-create.sh /usr/local/bin/
COPY init-firewall.sh /usr/local/bin/

RUN <<EOF
: 'add debian sid for latest git'
echo 'deb http://deb.debian.org/debian sid main' > /etc/apt/sources.list.d/sid.list; \
printf 'Package: *\nPin: release a=sid\nPin-Priority: 100\n\n' \
  >  /etc/apt/preferences.d/git-from-sid
printf 'Package: git*\nPin: release a=sid\nPin-Priority: 990\n' \
  >> /etc/apt/preferences.d/git-from-sid
: 'install mise deendencies and utils for init-firewall'
apt-get update \
  && apt-get -y --no-install-recommends install \
    sudo curl git ca-certificates build-essential \
    less procps fzf man-db unzip gnupg2 \
    iptables ipset iproute2 dnsutils aggregate jq \
  && rm -rf /var/lib/apt/lists/*
: 'install mise'
curl https://mise.run | sh
: 'create container user'
groupadd -g ${GID} -o ${CONTAINER_USER}
useradd -m -s /bin/bash -u ${UID} -g ${GID} -o ${CONTAINER_USER}
: 'create workspace and config directories'
mkdir -p /go /mise ${XDG_LOCAL} ${PNPM_HOME} ${WORKING_DIR} ${WORKING_DIR}/node_modules ${CLAUDE_CODE_CONFIG}
: 'grant access to container user'
chown -R ${CONTAINER_USER} /go /mise ${XDG_LOCAL} ${PNPM_HOME} ${WORKING_DIR} ${WORKING_DIR}/node_modules ${CLAUDE_CODE_CONFIG}
: 'setup start script'
chmod +x /usr/local/bin/start.sh
chmod +x /usr/local/bin/post-create.sh
: 'set up init-firewall script'
chmod +x /usr/local/bin/init-firewall.sh \
  && echo "%${CONTAINER_USER} ALL=(root) NOPASSWD: /usr/local/bin/init-firewall.sh" > /etc/sudoers.d/${CONTAINER_USER}-init-firewall \
  && chmod 0440 /etc/sudoers.d/${CONTAINER_USER}-init-firewall
: 'persist bash history'
SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  && mkdir /commandhistory \
  && touch /commandhistory/.bash_history \
  && chown -R ${CONTAINER_USER} /commandhistory
EOF

# configurations for container user
USER ${CONTAINER_USER}
WORKDIR ${WORKING_DIR}

# copy setting to home
COPY --chown=${UID}:${GID} mise-config.toml ${MISE_CONFIG_DIR}/config.toml
COPY --chown=${UID}:${GID} .cobra.yaml /home/${CONTAINER_USER}/.cobra.yaml

# go tools setting args
ARG DELVE_VERSION=latest
ARG GOPLS_VERSION=latest

RUN <<EOF
: 'allow the execution of git commands in a container workspace'
git config --global safe.directory ${WORKING_DIR}
: 'install go tools'
go install github.com/go-delve/delve/cmd/dlv@${DELVE_VERSION}
go install golang.org/x/tools/gopls@${GOPLS_VERSION}
EOF

# delve server port
ARG DEBUG_PORT=2345
EXPOSE ${DEBUG_PORT}
# start delve dap server
CMD ["start.sh"]
