FROM golang AS installer-builder

COPY cmd/installer /source

WORKDIR /source
RUN go run github.com/magefile/mage@latest -compile installer && chmod 755 installer


FROM golang AS containerrunner-builder

COPY cmd/containerrunner /source

WORKDIR /source
RUN go build -o containerrunner . && chmod 755 containerrunner


FROM archlinux

ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN pacman -Syu --noconfirm --needed rsync reflector
RUN reflector -c DE -f 12 -l 12 -n 12 -p https --completion-percent 98 --save /etc/pacman.d/mirrorlist
RUN sed -i -r '/^NoExtract / d' /etc/pacman.conf
RUN pacman -Qqn | pacman -S --noconfirm -
RUN pacman -S --noconfirm --needed sudo

RUN groupadd --gid $USER_GID $USERNAME && \
    useradd --uid $USER_UID --gid $USER_GID -m $USERNAME && \
    echo "$USERNAME ALL=(root) NOPASSWD:ALL" >/etc/sudoers.d/$USERNAME && \
    chmod 0440 /etc/sudoers.d/$USERNAME

COPY --from=installer-builder /source/installer /home/$USERNAME
COPY --from=containerrunner-builder /source/containerrunner /home/$USERNAME

RUN sudo -i -u $USERNAME /home/$USERNAME/installer -v && rm /home/$USERNAME/installer

RUN usermod -s /bin/zsh $USERNAME

USER $USERNAME
