FROM ubuntu:focal
CMD ["bash"]
RUN /bin/sh -c set -eux && apt-get update && apt-get install -y --no-install-recommends ca-certificates curl netbase wget 
RUN apt-get install -y debconf-utils
RUN curl -fLO https://go.dev/dl/go1.17.13.linux-amd64.tar.gz && tar xzf go1.17.13.linux-amd64.tar.gz && mv go /usr/local/ && ln -s /usr/local/go/bin/go /usr/local/bin/go 
RUN /bin/sh -c set -ex; export DEBIAN_FRONTEND=noninteractive; export DEBCONF_NONINTERACTIVE_SEEN=true; echo 'tzdata tzdata/Areas select America' | debconf-set-selections;echo 'tzdata tzdata/Zones/Etc select UTC' | debconf-set-selections;echo 'tzdata tzdata/Zones/America select Los_Angeles' | debconf-set-selections;apt-get install -y git automake make
RUN /bin/sh -c set -eux && rm -rf /var/lib/apt/lists/*
RUN git clone https://github.com/shastrinator/wsp -b initial_branch && cd wsp && make build-server && make build-client; cp wsp_server /usr/bin/; cp wsp_client /usr/bin/
RUN rm -rf wsp
