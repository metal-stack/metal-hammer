# Discover a device with metal-hammer

in order to be able to register a new device, or check whether a device is already registered, we execute from the pxeboot image a binary which does the hardware discovery. This is done with the *lshw* command which is available on various linux distributions, call it with the `-json` option as root and send the output to the maas api.

## Build

```bash
make
```

## Usage

### Configuration

Is done via Environment variables.

```bash
bin/metal-hammer -h

This application is configured via the environment. The following environment
variables can be used:

KEY                   TYPE             DEFAULT                                  REQUIRED    DESCRIPTION
METAL_HAMMER_DEBUG    True or False    false                                    False       turn on debug log
```

### Execution

```bash
sudo bin/metal-hammer
INFO[09-24|14:09:59] configuration                            debug=false reportURL=http://localhost:8080/device/register
INFO[09-24|14:10:00] device already registered                uuid=4C3CEF61-F536-B211-A85C-B765E03E138F caller=lshw.go:63
```
