# Use the official Go image with the latest version
FROM golang:latest

# Set the working directory inside the container
ARG USERNAME
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME -s /bin/bash

ENV USER=$USERNAME
USER $USERNAME

# Install any necessary Go tools (optional, customize as needed)
RUN go install golang.org/x/tools/gopls@latest

WORKDIR /akhanda