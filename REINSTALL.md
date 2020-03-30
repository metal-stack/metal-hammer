# OS Reinstallation process

Triggering an OS reinstallation starts by calling the **metal-api** REST endpoint `/v1/machine/<id>/reinstall` by providing the image ID of the new OS. It only proceeds if the given machine is already allocated and the new image ID is valid too.

**metal-api** marks the machine to get reinstalled by setting `allocation.Reinstall = true`. It then informs **metal-core** via NSQ about the desired reinstallation by sending the machine command `REINSTALL` event with the machine ID.

**metal-core** simply queries the IPMI details of the machine, set the boot order to PXE and power resets the machine.

**metal-hammer** reboots in PXE mode, brings all interfaces up, read the hardware details - and therewith creates a new password for the `metal` user - and registers the machine, just as usual.

It then fetches the machine data from **metal-api** and evaluates the `allocation.Reinstall` flag. If it's `false` it continues as usual, i.e. wiping all disks, etc. If it's `true`, which is the case in this scenario, it skips the usual process and first checks if there is an `allocation.BootInfo` struct given, which contains data of the currently given OS, i.e  the current `imageID`, `primaryDisk`, `osPartition`, `initrd`, `cmdline`, `kernel` and `bootloaderID` parameters.  
**metal-hammer** continues to wipe only the primary disk holding the current OS and leaving all other disks untouched! For this it has to check on beforehand if the current primary disk is the same as the one that will be used for the new OS. Therefore at least the current `imageID` or `primaryDisk` data is needed from the `BootInfo` struct. If they are both not available the procedure stops, since it would be too risky to continue regarding disk wiping.  
If only the `imageID` is given it tries to guess the primary disk of the old OS.  

After wiping the primary disk the reinstall procedure continues with the usual installation process up from the `installImage` method that eventually ends with the `finalizeAllocation` call, which now includes the previous mentioned `BootInfo` parameters.

**metal-core** passes-through the request to **metal-api**, sets the boot order to HD and power cycles the machine again, which in turn boots the new OS.
 
**metal-api** removes the `allocation.Reinstall` mark and stores the `BootInfo` details together with the newly installed `imageID` in the `allcation.MachineSetup` struct.

This was the happy-path. But of course, things can go wrong. If for any reason the reinstallation process fails, we are potentially in one of the following two states: Either the primary disk has been wiped already (and therewith the existing OS) or not. In both cases **metal-hammer** calls **metal-core** via the `/machine/abort-reinstall/<id>` endpoint delivering the bool value `primaryDiskWiped` that indicates the actual state.  
If **metal-core** fails to respond or the OS has already been wiped the machine reboots. Otherwise it gets the `BootInfo` of the previous installed OS stored in the DS and reboots with these details into the existing OS, just as nothing had happened at all.  

**metal-core** passes-through the abort request to **metal-api**, which in turn removes the `allocation.Reinstall` flag and returns the `BootInfo` if the OS has not been wiped yet. Otherwise it simply returns nothing.  
The latter case results in a new PXE boot and reinstallation process, which now could be succeed or again fail.  
This can potentially result in an endless reinstallation loop, but it ensures that no other disk than the one holding the OS will be wiped ever wiped!
