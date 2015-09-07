FROM opensuse:13.2
MAINTAINER Duncan Mac-Vicar P. "dmacvicar@suse.com"
RUN zypper --non-interactive ar -f http://download.opensuse.org/repositories/Virtualization/openSUSE_13.2/Virtualization.repo
#RUN zypper --non-interactive ar -f http://download.opensuse.org/repositories/Virt#ualization:/containers/openSUSE_13.2/Virtualization:containers.repo
RUN zypper --gpg-auto-import-keys --non-interactive in --no-recommends \
   docker curl
EXPOSE 2375
CMD ["/usr/bin/docker", "daemon", "-H", "tcp://0.0.0.0:2375", "--insecure-registry", "registry.test.lan", "-s", "overlay"]
