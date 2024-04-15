# syntax=docker/dockerfile:1.4
FROM golang:1.22.1-bookworm

# install mise with apt
## xz-utils (for tar -xJv in mise install)
RUN <<EOF
apt update -y && apt install -y gpg sudo wget curl xz-utils
sudo install -dm 755 /etc/apt/keyrings
wget -qO - https://mise.jdx.dev/gpg-key.pub | gpg --dearmor | sudo tee /etc/apt/keyrings/mise-archive-keyring.gpg 1> /dev/null
echo "deb [signed-by=/etc/apt/keyrings/mise-archive-keyring.gpg arch=amd64] https://mise.jdx.dev/deb stable main" | sudo tee /etc/apt/sources.list.d/mise.list
sudo apt update
sudo apt install -y mise
EOF

# create container user
ARG CONTAINER_USER=user
ARG UID=1000
ARG GID=1000
RUN <<EOF
groupadd -g ${GID} -o ${CONTAINER_USER}
useradd -m -s /bin/bash -u ${UID} -g ${GID} -o ${CONTAINER_USER}
EOF

# grant access to container user
COPY start.sh /usr/local/bin/start.sh
RUN <<EOF
chmod 755 /usr/local/bin/start.sh
chown -R ${CONTAINER_USER} /go
EOF

# configurations for container user
USER ${CONTAINER_USER}

ARG WORKING_DIR=/workspaces/app
WORKDIR ${WORKING_DIR}

# copy setting to home
COPY mise-config.toml /home/${CONTAINER_USER}/.config/mise/config.toml
COPY .cobra.yaml /home/${CONTAINER_USER}/.cobra.yaml

# go tools setting args
ARG DELVE_VERSION=latest

RUN <<EOF
: 'mise setup for bash'
echo 'eval "$(mise activate bash)"' >> ~/.bashrc
: 'allow the execution of git commands in a container workspace'
git config --global safe.directory ${WORKING_DIR}
: 'install go tools'
go install github.com/go-delve/delve/cmd/dlv@${DELVE_VERSION}
EOF

# delve server port
ARG DEBUG_PORT=2345
EXPOSE ${DEBUG_PORT}

# start delve dap server
CMD ["start.sh"]
