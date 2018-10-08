#!/bbin/elvish
/usr/sbin/rngd -r /dev/urandom -p /rngd.pid

/bbin/sshd -port 22 &

/bbin/metal-hammer
