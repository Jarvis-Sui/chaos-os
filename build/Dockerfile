FROM centos:7

ENV PATH "$PATH:/usr/local/go/bin"

RUN yum install -y wget && \
    yum install -y gcc && \
    yum install -y make && \
    wget https://golang.org/dl/go1.16.7.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.16.7.linux-amd64.tar.gz && \
    rm -f go1.16.7.linux-amd64.tar.gz
