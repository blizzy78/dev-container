FROM golang AS builder

RUN git clone https://github.com/magefile/mage && cd mage && go run bootstrap.go

COPY cmd/installer /installer
COPY cmd/containerrunner /containerrunner

WORKDIR /installer
RUN mage -compile installer && chmod 755 installer

WORKDIR /containerrunner
RUN go build -o containerrunner . && chmod 755 containerrunner


FROM ubuntu

ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN apt clean && rm -rf /var/lib/apt/lists/* && apt upgrade -y && apt update

RUN groupadd --gid $USER_GID $USERNAME && \
    useradd --uid $USER_UID --gid $USER_GID -m $USERNAME && \
    apt install -y sudo && \
    echo "$USERNAME ALL=(root) NOPASSWD:ALL" >/etc/sudoers.d/$USERNAME && \
    chmod 0440 /etc/sudoers.d/$USERNAME

COPY --from=builder /installer/installer /containerrunner/containerrunner /home/$USERNAME

RUN sudo -i -u $USERNAME /home/$USERNAME/installer -v && rm /home/$USERNAME/installer

RUN apt clean && rm -rf /var/lib/apt/lists/*

USER $USERNAME
