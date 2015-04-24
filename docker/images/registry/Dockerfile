FROM opensuse:13.2
MAINTAINER Duncan Mac-Vicar P. "dmacvicar@suse.com"
RUN zypper --non-interactive ar -f http://download.opensuse.org/repositories/Virtualization:/containers/openSUSE_13.2/Virtualization:containers.repo
RUN zypper --gpg-auto-import-keys --non-interactive in --no-recommends \
   docker-distribution-registry
EXPOSE 80

CMD ["/usr/bin/registry", "/etc/registry-config.yml"]
