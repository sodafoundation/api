# Build: sudo docker build . -t opensds/osdslet:integration-test
# Run: sudo docker run -d --net=host opensds/osdslet:integration-test

FROM leonwanghui/golang:1.9

MAINTAINER Leon Wang <wanghui71leon@gmail.com>

COPY opensds.conf /etc/opensds/opensds.conf
COPY osdslet /usr/bin

RUN mkdir -p /var/log/opensds

# Define default command.
CMD ["/usr/bin/osdslet"]
