# Build: sudo docker build . -t opensds/osdsdock:integration-test
# Run: sudo docker run -d --net=host opensds/osdsdock:integration-test

FROM leonwanghui/golang:1.9

MAINTAINER Leon Wang <wanghui71leon@gmail.com>

COPY opensds.conf /etc/opensds/opensds.conf
COPY osdsdock /usr/bin

RUN mkdir -p /var/log/opensds
# Install some packages before running command.
RUN apt-get update && apt-get install -y librados-dev librbd-dev

# Define default command.
CMD ["/usr/bin/osdsdock"]
