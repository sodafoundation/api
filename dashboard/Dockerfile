# Docker build usage:
# 	docker build . -t opensdsio/dashboard:latest
# Docker run usage:
# 	docker run -d -p 8088:8088 opensdsio/dashboard:latest

FROM ubuntu:16.04
MAINTAINER Leon Wang <wanghui71leon@gmail.com>

ARG DEBIAN_FRONTEND=noninteractive

# Download and install some packages.
RUN apt-get update && apt-get install -y --no-install-recommends \
  sudo \
  wget \
  make \
  g++ \
  nginx \
  && rm -rf /var/lib/apt/lists/* \
  && apt-get clean
RUN wget --no-check-certificate https://deb.nodesource.com/setup_8.x \
  && chmod +x setup_8.x && ./setup_8.x \
  && apt-get install -y nodejs

# Current directory is always /opt/dashboard.
WORKDIR /opt/dashboard

# Copy dashboard source code into container before running command.
COPY ./ ./

RUN chmod 755 ./image_builder.sh \
  && sudo ./image_builder.sh
RUN sudo ./image_builder.sh --rebuild

# Define default command.
CMD /usr/sbin/nginx -g "daemon off;"
