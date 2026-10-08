Vagrant.configure("2") do |config|
  config.vm.define "ubuntu" do |u|
    u.vm.box = "ubuntu/jammy64"
    u.vm.hostname = "devstack-ubuntu"

    u.vm.provider "virtualbox" do |vbox|
      vbox.name = "vitualizz-devstack-ubuntu"
      vbox.cpus = 2
      vbox.memory = 2048
      vbox.default_nic_type = "virtio"
    end

    u.vm.synced_folder ".", "/vagrant", disabled: false
  end

  config.vm.define "arch" do |a|
    a.vm.box = "archlinux/archlinux"
    a.vm.hostname = "devstack-arch"

    a.vm.provider "virtualbox" do |vbox|
      vbox.name = "vitualizz-devstack-arch"
      vbox.cpus = 2
      vbox.memory = 2048
      vbox.default_nic_type = "virtio"
    end

    a.vm.synced_folder ".", "/vagrant", disabled: false
  end

  config.vm.synced_folder ".", "/vagrant",
    type: "rsync",
    rsync__exclude: [".git/"]
end
