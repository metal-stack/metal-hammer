#!/bbin/elvish
# TODO only start if running in a virtual machine.
/usr/sbin/rngd -r /dev/urandom -p /rngd.pid

/bbin/sshd -port 22 &

/bbin/metal-hammer
