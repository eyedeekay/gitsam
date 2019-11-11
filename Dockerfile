FROM alpine:edge
ARG REPO_DIR=eyedeekay
COPY $REPO_DIR /var/www
COPY .gitsam_secure/id_rsa.pub /etc/gitsam_secure/id_rsa.pub
COPY gitsam/gitsam /usr/bin/gitsam
WORKDIR /var/www
CMD gitsam -il 1 -ol 1 -iq 8 -oq 8 -ib 3 -ob 3 -pk /etc/gitsam_secure/id_rsa.pub