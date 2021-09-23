FROM golang AS mage-builder

RUN git clone https://github.com/magefile/mage && cd mage && go run bootstrap.go
WORKDIR /
ADD go.mod go.sum install/*.go /
RUN mage -compile mage-static && chmod 755 mage-static


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

ADD install/tools.sh install/go.sh /home/$USERNAME
RUN sudo -i -u $USERNAME bash </home/$USERNAME/tools.sh
RUN sudo -i -u $USERNAME bash </home/$USERNAME/go.sh
RUN rm /home/$USERNAME/tools.sh /home/$USERNAME/go.sh

COPY --from=mage-builder /mage-static /home/$USERNAME
RUN sudo -i -u $USERNAME /home/$USERNAME/mage-static -v && rm /home/$USERNAME/mage-static

RUN apt clean && rm -rf /var/lib/apt/lists/*

USER $USERNAME
