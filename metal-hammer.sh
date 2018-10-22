#!/bbin/elvish
# rngd is not needed anymore in a vm because we added:
# domain.random :model => 'random' in the Vagrantfile
# for reference see:
# https://github.com/vagrant-libvirt/vagrant-libvirt#random-number-generator-passthrough
#
# /usr/sbin/rngd -r /dev/urandom -p /rngd.pid

/bbin/sshd -port 22 &

# golang hardware gather lib needs syslog file to get physical memory.
/bbin/dmesg > /var/log/syslog

/bbin/metal-hammer
