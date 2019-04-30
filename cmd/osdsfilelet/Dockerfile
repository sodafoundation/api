# Docker build usage:
# 	docker build . -t opensdsio/opensds-filecontroller:latest
# Docker run usage:
# 	docker run -d --net=host -v /etc/opensds:/etc/opensds opensdsio/opensds-filecontroller:latest

FROM ubuntu:16.04

COPY osdsfilelet /usr/bin

# Define default command.
CMD ["/usr/bin/osdsfilelet"]
