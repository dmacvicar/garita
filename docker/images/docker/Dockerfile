FROM opensuse:13.2
MAINTAINER Duncan Mac-Vicar P. "dmacvicar@suse.com"
RUN zypper --non-interactive ar -f http://download.opensuse.org/repositories/Virtualization:/containers/openSUSE_13.2/Virtualization:containers.repo
RUN zypper --gpg-auto-import-keys --non-interactive in --no-recommends \
   docker
EXPOSE 2375
CMD ["/usr/bin/docker", "-d", "-H", "tcp://0.0.0.0:2375", "--insecure-registry", "registry.test.lan"]
