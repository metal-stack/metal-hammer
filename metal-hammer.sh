#!/bbin/elvish
/usr/sbin/rngd -p /rngd.pid

/bbin/sshd -port 22 &

/bbin/metal-hammer
