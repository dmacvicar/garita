# -*- mode: ruby -*-
# vi: set ft=ruby :
Vagrant.require_version '>= 1.6.0'
VAGRANTFILE_API_VERSION = '2'
ENV['VAGRANT_DEFAULT_PROVIDER'] = 'docker'
DOMAIN = 'test.lan'

# Create and configure the Docker container(s)
Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|

  config.vm.define 'garita' do |container|
    container.vm.provider 'docker' do |d|
      d.image           = 'opensuse:13.2'
      d.volumes         = [File.join(ENV['GOPATH'], 'bin/garita') +
       ':/usr/local/bin/garita']
      d.cmd             = ['/usr/local/bin/garita']
      d.name            = 'garita'
      d.create_args     = ['-h', d.name + ".#{DOMAIN}", '--dns-search', DOMAIN]
      d.expose          = [80]
    end
  end

  config.vm.define 'registry' do |container|
    container.vm.provider 'docker' do |d|
      d.build_dir       = './docker/images/registry'
      d.name            = 'registry'
      d.volumes         =
        [File.expand_path('vagrant/conf/registry-config.yml') +
          ':/etc/registry-config.yml']
      d.create_args     = ['-h', d.name + ".#{DOMAIN}", '--dns-search', DOMAIN]
      d.link 'garita:garita.test.lan'
      d.expose          = [80]
    end
  end

  config.vm.define 'dockerd' do |container|
    container.vm.provider 'docker' do |d|
      d.build_dir       = './docker/images/docker'
      d.privileged      = true
      d.name            = 'dockerd'
      d.create_args     = ['-h', d.name + ".#{DOMAIN}", '--dns-search', DOMAIN]
      d.ports           = ['23750:2375']
      d.link 'registry:registry.test.lan'
    end
  end

end
