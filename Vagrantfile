# -*- mode: ruby -*-
# vi: set ft=ruby :

current_file_path = File.dirname(__FILE__)
kernel_path = File.join(current_file_path, "pxeboot-kernel")
initrd_path = File.join(current_file_path, "pxeboot-initrd.img")

Vagrant.configure("2") do |config|
  config.vm.provider "libvirt" do |domain|
    domain.default_prefix = "metal"
    domain.keymap = 'de'
  end

  config.vm.define :pxeclient do |pxeclient|
    pxeclient.vm.provider :libvirt do |domain|
      domain.cpus = 1
      domain.memory = 1024
      domain.storage :file, :size => '6G', :bus => 'sata'
      domain.boot 'hd'
      domain.kernel = kernel_path
      domain.initrd = initrd_path
      domain.cmd_line = "console=tty0 console=ttyS0 METAL_CORE_URL=http://192.168.121.110:4242 IMAGE_URL=http://192.168.121.1:4711/images/os/ubuntu/18.04/img.tar.gz"
      domain.loader = "/usr/share/OVMF/OVMF_CODE.fd"
    end
  end
end
