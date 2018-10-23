# -*- mode: ruby -*-
# vi: set ft=ruby :

current_file_path = File.dirname(__FILE__)
kernel_path = File.join(current_file_path, "metal-hammer-kernel")
initrd_path = File.join(current_file_path, "metal-hammer-initrd.img.lz4")

Vagrant.configure("2") do |config|
  config.vm.provider "libvirt" do |domain|
    domain.default_prefix = "metal-hammer"
    domain.keymap = 'de'
    domain.random :model => 'random'
  end

  config.vm.define :pxeclient do |pxeclient|
    pxeclient.trigger.before :up do |trigger|
      trigger.info = "Download kernel..."
      trigger.run = {path: "download-kernel.sh"}
    end

    pxeclient.vm.provider :libvirt do |domain|
      domain.cpus = 1
      domain.memory = 1024
      domain.storage :file, :size => '1500M', :bus => 'sata'
      domain.storage :file, :size => '10M', :bus => 'sata'
      domain.boot 'hd'
      domain.kernel = kernel_path
      domain.initrd = initrd_path
      domain.cmd_line = "console=tty0 console=ttyS0 ip=dhcp METAL_CORE_URL=http://192.168.121.110:4242 IMAGE_URL=http://192.168.121.1:4711/images/os/ubuntu/18.10/img.tar.lz4 BGP=1"
      domain.loader = "/usr/share/OVMF/OVMF_CODE.fd"
    end
  end
end
