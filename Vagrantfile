# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  # Use Ubuntu 22.04 LTS as the base box
  # Detect architecture: on Apple Silicon, use Parallels (libvirt is Linux-only)
  host_arch = `uname -m`.strip
  if host_arch == 'arm64' || host_arch == 'aarch64'
    # For Apple Silicon: Parallels (Pro/Business) or VMware Fusion (free for personal use)
    # bento/ubuntu-22.04 v202401.31.0 has parallels/arm64 and vmware_desktop/arm64
    config.vm.box = "bento/ubuntu-22.04"
    config.vm.box_version = "~> 202401.31.0"
    config.vm.provider "parallels" do |prl|
      prl.memory = 2048
      prl.cpus = 2
    end
    config.vm.provider "vmware_desktop" do |vmw|
      vmw.memory = 2048
      vmw.cpus = 2
    end
  else
    # For Intel/AMD: use VirtualBox with standard Ubuntu box
    config.vm.box = "ubuntu/jammy64"
    config.vm.provider "virtualbox" do |vb|
      vb.name = "sre-chat-api-production"
      vb.memory = "2048"
      vb.cpus = 2
    end
  end
  config.vm.box_check_update = false

  # VM Configuration
  config.vm.hostname = "sre-chat-api-prod"
  config.vm.network "forwarded_port", guest: 80, host: 8080
  config.vm.network "forwarded_port", guest: 8080, host: 8081
  config.vm.network "forwarded_port", guest: 5432, host: 5433

  # Provisioning: Install dependencies and configure services
  config.vm.provision "shell", path: "provision/setup.sh", privileged: true

  # Sync project directory to VM
  config.vm.synced_folder ".", "/vagrant", disabled: false

  # Post-provisioning: Start services
  config.vm.provision "shell", path: "provision/deploy.sh", privileged: false, run: "always"
end
