# Docker build usage:
#   docker build . -t opensdsio/opensds-authchecker:latest

FROM centos:7
MAINTAINER Leon Wang <wanghui71leon@gmail.com>

COPY start.sh /start.sh

RUN yum -y update && \
    yum -y install centos-release-openstack-queens && \
    yum -y install openstack-keystone openstack-utils python-openstackclient httpd mod_wsgi initscripts && \
    yum clean all

RUN chmod +x /start.sh && \
    chmod +x /etc/rc.d/rc.local && \
    sed -i '13i bash \/start\.sh >> \/start.log' /etc/rc.d/rc.local

CMD ["/sbin/init"]

