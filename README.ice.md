# Intel E810 (100GB) network cards

## DDP package

The dynamic device personalization package (`ice.pkg`) required by the `ice` driver is downloaded
during the docker build from [intel/ethernet-linux-ice](https://github.com/intel/ethernet-linux-ice)
and placed in the initrd at `/lib/firmware/intel/ice/ddp/ice.pkg`, see `ICE_VERSION` in the
`Dockerfile`.

## Firmware (NVM) update

metal-hammer may update the NVM of e810 based network cards on every boot, before the machine registers
itself. The update is best effort, it never fails the provisioning of a machine: whatever can not be
inspected, flashed or activated is logged and retried on the next boot.

- cards bound to the `ice` driver are detected through sysfs, their NVM version is read with
  `ethtool -i <interface>`. A card which can not be inspected is skipped
- the desired version is the one of the update package (that can be added to the docker build)
- versions are compared as semver with the minor padded to two digits: `4.8` is `4.80` and therefore
  newer than `4.60`, not older
- if the lowest version found is behind, all cards are flashed with one run of
  `/intel/nvmupdate64e -u -s -l /tmp/nvmupdate.log -o /tmp/nvmupdate.xml -c /intel/nvmupdate.cfg -a /intel`.
  The utility inventories the adapters itself and leaves cards alone which are up to date or have no
  image in the package, e.g. an x710
- whether and how the machine is restarted afterwards is read from the machine readable XML report the
  utility writes with `-o`: its top level `PowerCycleRequired` and `RebootRequired` counters state
  whether a card was flashed and how the new firmware is activated. The e810 nvm is written into an
  inactive bank which only becomes active once the card lost power, so it requests a **power cycle**
  through the BMC. Both counters are zero when the run flashed nothing, so a card which can not be
  flashed at all, e.g. an oem locked nvm, stays behind the packaged version without triggering a restart
- a card which is behind the packaged version but did not take the update is reported to metal-api as
  a provisioning event.

The update utility is not redistributable, so it is not vendored into this repository. It has to be
downloaded from the
[Intel download page](https://www.intel.com/content/www/us/en/download/19626) and the resulting url
must be passed to the docker build:

```bash
docker buildx build \
  --build-arg E810_NVM_URL=<url-of-E810_NVMUpdatePackage_vX_YY_Linux.tar.gz> \
  --build-arg E810_NVM_SHA256=<sha256-of-that-package> .
```

Intel changes the numeric mirror id of that url with every release. If `E810_NVM_URL` is not set,
the utility is not part of the initrd and the firmware update is skipped at runtime.

The package ships a binary which metal-hammer executes as root, so `E810_NVM_SHA256` is mandatory
if `E810_NVM_URL` is set and the download is verified against it before it is unpacked.

### Raising the firmware version

`E810_NVM_URL` and `E810_NVM_SHA256` have to be changed together, the build fails if the download
does not match the digest. The version is derived from the filename of the package
(`E810_NVMUpdatePackage_v4_80_Linux.tar.gz` becomes `4.80`) and stored next to the utility, so the
shipped firmware and the version metal-hammer expects can not get out of sync.

[Intel Driver Download](https://www.intel.com/content/www/us/en/search.html?ws=text#q=e810&t=Downloads&layout=table)
