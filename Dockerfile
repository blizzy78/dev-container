FROM golang:1.16 AS builder

RUN git clone https://github.com/magefile/mage && cd mage && go run bootstrap.go
WORKDIR /
ADD go.mod go.sum install/*.go /
RUN mage -compile mage-static


FROM ubuntu

ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN apt update && apt upgrade -y

RUN groupadd --gid $USER_GID $USERNAME && \
    useradd --uid $USER_UID --gid $USER_GID -m $USERNAME && \
    apt install -y sudo && \
    echo "$USERNAME ALL=(root) NOPASSWD:ALL" >/etc/sudoers.d/$USERNAME && \
    chmod 0440 /etc/sudoers.d/$USERNAME

ADD install/tools.sh /home/$USERNAME
RUN chmod a+rx /home/$USERNAME/tools.sh && sudo -i -u $USERNAME /home/$USERNAME/tools.sh && rm /home/$USERNAME/tools.sh

ADD install/go.sh /home/$USERNAME
RUN chmod a+rx /home/$USERNAME/go.sh && sudo -i -u $USERNAME /home/$USERNAME/go.sh && rm /home/$USERNAME/go.sh

COPY --from=builder /mage-static /home/$USERNAME
RUN chmod 755 /home/$USERNAME/mage-static && sudo -i -u $USERNAME /home/$USERNAME/mage-static -v && rm /home/$USERNAME/mage-static

USER $USERNAME
