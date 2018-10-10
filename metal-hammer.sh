#!/bbin/elvish
# TODO only start if running in a virtual machine.
/usr/sbin/rngd -r /dev/urandom -p /rngd.pid

/bbin/sshd -port 22 &

# golang hardware gather lib needs syslog file to get physical memory.
/bbin/dmesg > /var/log/syslog

/bbin/metal-hammer
