# Discover a device

in order to be able to register a new device, or check whether a device is already registered, we execute from the pxeboot image a binary which does the hardware discovery. This is done with the *lshw* command which is available on various linux distributions, call it with the `-json` option as root and send the output to the maas api.

## Build

```bash
make
```

## Usage

```bash
sudo .bin/discover
INFO[09-24|14:09:59] configuration                            debug=false reportURL=http://localhost:8080/device/register
INFO[09-24|14:10:00] device already registered                uuid=4C3CEF61-F536-B211-A85C-B765E03E138F caller=lshw.go:63
```
