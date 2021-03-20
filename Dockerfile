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

ADD install/go-tools.sh /home/$USERNAME
RUN chmod a+rx /home/$USERNAME/go-tools.sh && sudo -i -u $USERNAME /home/$USERNAME/go-tools.sh && rm /home/$USERNAME/go-tools.sh

ADD install/protoc.sh /home/$USERNAME
RUN chmod a+rx /home/$USERNAME/protoc.sh && sudo -i -u $USERNAME /home/$USERNAME/protoc.sh && rm /home/$USERNAME/protoc.sh

ADD install/npm.sh /home/$USERNAME
RUN chmod a+rx /home/$USERNAME/npm.sh && sudo -i -u $USERNAME /home/$USERNAME/npm.sh && rm /home/$USERNAME/npm.sh

ADD install/postcss.sh /home/$USERNAME
RUN chmod a+rx /home/$USERNAME/postcss.sh && sudo -i -u $USERNAME /home/$USERNAME/postcss.sh && rm /home/$USERNAME/postcss.sh

ADD install/java.sh /home/$USERNAME
RUN chmod a+rx /home/$USERNAME/java.sh && sudo -i -u $USERNAME /home/$USERNAME/java.sh && rm /home/$USERNAME/java.sh

ADD install/maven.sh /home/$USERNAME
RUN chmod a+rx /home/$USERNAME/maven.sh && sudo -i -u $USERNAME /home/$USERNAME/maven.sh && rm /home/$USERNAME/maven.sh

ADD install/volumes.sh /home/$USERNAME
RUN chmod a+rx /home/$USERNAME/volumes.sh && sudo -i -u $USERNAME /home/$USERNAME/volumes.sh && rm /home/$USERNAME/volumes.sh

ADD install/misc.sh /home/$USERNAME
RUN chmod a+rx /home/$USERNAME/misc.sh && sudo -i -u $USERNAME /home/$USERNAME/misc.sh && rm /home/$USERNAME/misc.sh

USER $USERNAME
