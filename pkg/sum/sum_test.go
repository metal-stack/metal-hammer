package sum

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnmarshalS2BiosCfg(t *testing.T) {
	// given
	s := &sum{}
	s.biosCfgXML = testS2BiosCfg

	// when
	err := s.unmarshalBiosCfg()

	// when
	s.determineMachineType()

	// then
	require.Equal(t, s2, s.machineType)

	// then
	require.Nil(t, err)

	// when
	err = s.findUEFINetworkBootOption()

	// then
	require.Nil(t, err)
	require.Equal(t, "UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection", s.uefiNetworkBootOption)
}

func TestUnmarshalBigTwinBiosCfg(t *testing.T) {
	// given
	s := &sum{}
	s.biosCfgXML = testBigTwinBiosCfg

	// when
	err := s.unmarshalBiosCfg()

	// when
	s.determineMachineType()

	// then
	require.Equal(t, bigTwin, s.machineType)

	// when
	s.determineSecureBoot()

	// then
	require.True(t, s.secureBootEnabled)

	// then
	require.Nil(t, err)

	// when
	err = s.findUEFINetworkBootOption()

	// then
	require.Nil(t, err)
	require.Equal(t, "UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28", s.uefiNetworkBootOption)
}

const (
	testS2BiosCfg = `<?xml version="1.0" encoding="ISO-8859-1" standalone="yes"?>
<BiosCfg>
  <!--Supermicro Update Manager 2.3.0 (2019/08/08)-->
  <!--File generated at 2019-12-18_12:31:40-->
  <Menu name="Main">
    <Information />
    <Subtitle></Subtitle>
    <Subtitle></Subtitle>
    <Subtitle>Supermicro X11SDV-8C-TP8F</Subtitle>
    <Text>BIOS Version(1.1a)</Text>
    <Text>Build Date(05/17/2019)</Text>
    <Subtitle></Subtitle>
    <Subtitle>Memory Information</Subtitle>
    <Text>Total Memory(131072 MB)</Text>
  </Menu>
  <Menu name="Advanced">
    <Information />
    <Menu name="Boot Feature">
      <Information>
        <Help><![CDATA[Boot Feature Configuration Page]]></Help>
      </Information>
      <Subtitle></Subtitle>
      <Setting name="Quiet Boot" checkedStatus="Checked" type="CheckBox">
        <!--Checked/Unchecked-->
        <Information>
          <DefaultStatus>Checked</DefaultStatus>
          <Help><![CDATA[Enables or disables Quiet Boot option]]></Help>
        </Information>
      </Setting>
      <Subtitle></Subtitle>
      <Setting name="Option ROM Messages" selectedOption="Force BIOS" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Force BIOS</Option>
            <Option value="0">Keep Current</Option>
          </AvailableOptions>
          <DefaultOption>Force BIOS</DefaultOption>
          <Help><![CDATA[Set display mode for Option ROM]]></Help>
        </Information>
      </Setting>
      <Setting name="Bootup NumLock State" selectedOption="On" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">On</Option>
            <Option value="0">Off</Option>
          </AvailableOptions>
          <DefaultOption>On</DefaultOption>
          <Help><![CDATA[Select the keyboard NumLock state]]></Help>
        </Information>
      </Setting>
      <Setting name="Wait For &quot;F1&quot; If Error" selectedOption="Enabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Enabled</DefaultOption>
          <Help><![CDATA[Enable- BIOS will wait for user to press "F1" if some error happens. Disable- BIOS will continue to POST, user interaction not required]]></Help>
        </Information>
      </Setting>
      <Setting name="INT19 Trap Response" selectedOption="Immediate" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Immediate</Option>
            <Option value="0">Postponed</Option>
          </AvailableOptions>
          <DefaultOption>Immediate</DefaultOption>
          <Help><![CDATA[BIOS reaction on INT19 trapping by Option ROM: IMMEDIATE - execute the trap right away; POSTPONED - execute the trap during legacy boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="Re-try Boot" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Legacy Boot</Option>
            <Option value="2">EFI Boot</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[Decide how to retry boot devices which fail to boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="Port 61h Bit-4 Emulation" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[Emulation of Port 61h bit-4 toggling in SMM]]></Help>
        </Information>
      </Setting>
      <Subtitle></Subtitle>
      <Subtitle>Power Configuration</Subtitle>
      <Setting name="Watch Dog Function" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[Enable or disable to turn on 5-minute watch dog timer. Upon timeout, JWD1 jumper determines system behavior.]]></Help>
        </Information>
      </Setting>
      <Setting name="Power Button Function" selectedOption="Instant Off" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Instant Off</Option>
            <Option value="0">4 Seconds Override</Option>
          </AvailableOptions>
          <DefaultOption>Instant Off</DefaultOption>
          <Help><![CDATA[Instant Off: Turn off system immediately in legacy OS.
4 Seconds Override: Turn off system after depressed for 4 seconds.]]></Help>
        </Information>
      </Setting>
      <Setting name="Restore on AC Power Loss" selectedOption="Last State" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Stay Off</Option>
            <Option value="1">Power On</Option>
            <Option value="2">Last State</Option>
          </AvailableOptions>
          <DefaultOption>Last State</DefaultOption>
          <Help><![CDATA[Stay Off: System always remains off.
Power On: System always turns on.
Last State: System returns to previous state before AC lost.]]></Help>
        </Information>
      </Setting>
    </Menu>
    <Menu name="CPU Configuration" order="1">
      <Information>
        <Help><![CDATA[CPU Configuration]]></Help>
      </Information>
      <Subtitle>Processor Configuration</Subtitle>
      <Subtitle>--------------------------------------------------</Subtitle>
      <Text>Processor BSP Revision(50654 - SKX M0)</Text>
      <Text>Processor Socket(CPU1)</Text>
      <Text>Processor ID(00050654*)</Text>
      <Text>Processor Frequency(2.300GHz)</Text>
      <Text>Processor Max Ratio(     17H)</Text>
      <Text>Processor Min Ratio(     0AH)</Text>
      <Text>Microcode Revision(0200005E)</Text>
      <Text>L1 Cache RAM(    64KB)</Text>
      <Text>L2 Cache RAM(  1024KB)</Text>
      <Text>L3 Cache RAM( 11264KB)</Text>
      <Subtitle>Processor 0 Version</Subtitle>
      <Subtitle>Intel(R) Xeon(R) D-2146NT CPU @ 2.30GHz</Subtitle>
      <Subtitle></Subtitle>
      <Setting name="Hyper-Threading [ALL]" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Disable</Option>
            <Option value="0">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Enables Hyper Threading (Software Method to Enable/Disable Logical Processor threads.]]></Help>
        </Information>
      </Setting>
      <Setting name="Cores Enabled" numericValue="0" type="Numeric">
        <Information>
          <MaxValue>28</MaxValue>
          <MinValue>0</MinValue>
          <StepSize>1</StepSize>
          <DefaultValue>0</DefaultValue>
          <Help><![CDATA[Number of Cores to Enable in each Processor Package. 0 means all cores. Total 8 cores available in each CPU package.]]></Help>
        </Information>
      </Setting>
      <Setting name="Execute Disable Bit" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[When disabled, forces the XD feature flag to always return 0.]]></Help>
        </Information>
      </Setting>
      <Setting name="Intel Virtualization Technology" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[When enabled, a VMM can utilize the additional hardware capabilities provided by Vanderpool Technology]]></Help>
        </Information>
      </Setting>
      <Setting name="PPIN Control" selectedOption="Unlock/Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Unlock/Disable</Option>
            <Option value="1">Unlock/Enable</Option>
          </AvailableOptions>
          <DefaultOption>Unlock/Enable</DefaultOption>
          <Help><![CDATA[Unlock and Enable/Disable PPIN Control]]></Help>
        </Information>
      </Setting>
      <Setting name="Hardware Prefetcher" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Enable</Option>
            <Option value="0">Disable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[= MLC Streamer Prefetcher (MSR 1A4h Bit[0])]]></Help>
        </Information>
      </Setting>
      <Setting name="Adjacent Cache Prefetch" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Enable</Option>
            <Option value="0">Disable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[= MLC Spatial Prefetcher (MSR 1A4h Bit[1])]]></Help>
        </Information>
      </Setting>
      <Setting name="DCU Streamer Prefetcher" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Enable</Option>
            <Option value="0">Disable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[DCU streamer prefetcher is an L1 data cache prefetcher (MSR 1A4h [2]).]]></Help>
        </Information>
      </Setting>
      <Setting name="DCU IP Prefetcher" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Enable</Option>
            <Option value="0">Disable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[DCU IP prefetcher is an L1 data cache prefetcher (MSR 1A4h [3]).]]></Help>
        </Information>
      </Setting>
      <Setting name="LLC Prefetch" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[Enable/Disable LLC Prefetch on all threads]]></Help>
        </Information>
      </Setting>
      <Setting name="Extended APIC" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[Enable/disable extended APIC support]]></Help>
        </Information>
      </Setting>
      <Setting name="AES-NI" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Enable/disable AES-NI support]]></Help>
        </Information>
      </Setting>
      <Menu name="Advanced Power Management Configuration">
        <Information>
          <Help><![CDATA[Displays and provides option to change the Power Management Settings]]></Help>
        </Information>
        <Subtitle>Advanced Power Management Configuration</Subtitle>
        <Subtitle>--------------------------------------------------</Subtitle>
        <Setting name="Power Technology" selectedOption="Energy Efficient" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">Disable</Option>
              <Option value="1">Energy Efficient</Option>
              <Option value="2">Custom</Option>
            </AvailableOptions>
            <DefaultOption>Energy Efficient</DefaultOption>
            <Help><![CDATA[Switch CPU Power Management profile]]></Help>
          </Information>
        </Setting>
        <Setting name="Power Performance Tuning" selectedOption="OS Controls EPB" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">OS Controls EPB</Option>
              <Option value="1">BIOS Controls EPB</Option>
            </AvailableOptions>
            <DefaultOption>OS Controls EPB</DefaultOption>
            <Help><![CDATA[MSR 1FCh Bit[25] = PWR_PERF_TUNING_CFG_MODE. Enable - Use IA32_ENERGY_PERF_BIAS input from the core;
Disable - Use alternate perf BIAS input from ENERGY_PERF_BIAS_CONFIG]]></Help>
            <WorkIf><![CDATA[ ( 2 == Power Technology )  and  ( 2 != Hardware P-States ) ]]></WorkIf>
          </Information>
        </Setting>
        <Setting name="ENERGY_PERF_BIAS_CFG mode" selectedOption="Balanced Performance" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="3">Maximum Performance</Option>
              <Option value="0">Performance</Option>
              <Option value="7">Balanced Performance</Option>
              <Option value="8">Balanced Power</Option>
              <Option value="15">Power</Option>
            </AvailableOptions>
            <DefaultOption>Balanced Performance</DefaultOption>
            <Help><![CDATA[Use input from ENERGY_PERF_BIAS_CONFIG mode selection. PERF/Balanced Perf/Balanced Power/Power]]></Help>
            <WorkIf><![CDATA[ ( 2 == Power Technology )  and  ( 0 != Power Performance Tuning ) ]]></WorkIf>
          </Information>
        </Setting>
        <Menu name="CPU P State Control">
          <Information>
            <Help><![CDATA[P State Control Configuration Sub Menu, include Turbo, XE and etc.]]></Help>
            <WorkIf><![CDATA[  2 == Power Technology  ]]></WorkIf>
          </Information>
          <Subtitle>CPU P State Control</Subtitle>
          <Subtitle></Subtitle>
          <Setting name="Uncore Freq Scaling (UFS)" selectedOption="Enable" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Enable</Option>
                <Option value="1">Disable</Option>
              </AvailableOptions>
              <DefaultOption>Enable</DefaultOption>
              <Help><![CDATA[Enable/Disable autonomous uncore frequency scaling]]></Help>
            </Information>
          </Setting>
          <Setting name="SpeedStep (Pstates)" selectedOption="Enable" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
              </AvailableOptions>
              <DefaultOption>Enable</DefaultOption>
              <Help><![CDATA[Enable/Disable EIST (P-States)]]></Help>
            </Information>
          </Setting>
          <Setting name="Config TDP" selectedOption="Nominal" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Nominal</Option>
                <Option value="1">Level 1</Option>
                <Option value="2">Level 2</Option>
              </AvailableOptions>
              <DefaultOption>Nominal</DefaultOption>
              <Help><![CDATA[Config TDP level selection]]></Help>
              <WorkIf><![CDATA[  0 != SpeedStep (Pstates)  ]]></WorkIf>
            </Information>
          </Setting>
          <Setting name="EIST PSD Function" selectedOption="HW_ALL" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">HW_ALL</Option>
                <Option value="1">SW_ALL</Option>
              </AvailableOptions>
              <DefaultOption>HW_ALL</DefaultOption>
              <Help><![CDATA[Choose HW_ALL/SW_ALL in _PSD return]]></Help>
              <WorkIf><![CDATA[  0 != SpeedStep (Pstates)  ]]></WorkIf>
            </Information>
          </Setting>
          <Setting name="Energy Efficient Turbo" selectedOption="Enable" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Enable</Option>
                <Option value="1">Disable</Option>
              </AvailableOptions>
              <DefaultOption>Enable</DefaultOption>
              <Help><![CDATA[Energy Efficient Turbo Disable, MSR 0x1FC [19]]]></Help>
            </Information>
          </Setting>
          <Setting name="Turbo Mode" selectedOption="Enable" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
              </AvailableOptions>
              <DefaultOption>Enable</DefaultOption>
              <Help><![CDATA[Enable/Disable processor Turbo Mode (requires EMTTM enabled too).]]></Help>
              <WorkIf><![CDATA[  0 != SpeedStep (Pstates)  ]]></WorkIf>
            </Information>
          </Setting>
        </Menu>
        <Menu name="Hardware PM State Control">
          <Information>
            <Help><![CDATA[Hardware P-State setting]]></Help>
            <WorkIf><![CDATA[  2 == Power Technology  ]]></WorkIf>
          </Information>
          <Subtitle>Hardware PM State Control</Subtitle>
          <Subtitle></Subtitle>
          <Setting name="Hardware P-States" selectedOption="Disable" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Native Mode</Option>
                <Option value="2">Out of Band Mode</Option>
                <Option value="3">Native Mode with No Legacy Support</Option>
              </AvailableOptions>
              <DefaultOption>Disable</DefaultOption>
              <Help><![CDATA[Disable: Hardware chooses a P-state based on OS Request (Legacy P-States)
Native Mode:Hardware chooses a P-state based on OS guidance
Out of Band Mode:Hardware autonomously chooses a P-state (no OS guidance)]]></Help>
            </Information>
          </Setting>
        </Menu>
        <Menu name="CPU C State Control">
          <Information>
            <Help><![CDATA[CPU C State setting]]></Help>
            <WorkIf><![CDATA[  2 == Power Technology  ]]></WorkIf>
          </Information>
          <Subtitle>CPU C State Control</Subtitle>
          <Subtitle></Subtitle>
          <Setting name="Autonomous Core C-State" selectedOption="Disable" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
              </AvailableOptions>
              <DefaultOption>Disable</DefaultOption>
              <Help><![CDATA[Autonomous Core C-State Control]]></Help>
            </Information>
          </Setting>
          <Setting name="CPU C6 report" selectedOption="Auto" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
                <Option value="255">Auto</Option>
              </AvailableOptions>
              <DefaultOption>Auto</DefaultOption>
              <Help><![CDATA[Enable/Disable CPU C6(ACPI C3) report to OS]]></Help>
              <WorkIf><![CDATA[  1 != Autonomous Core C-State  ]]></WorkIf>
            </Information>
          </Setting>
          <Setting name="Enhanced Halt State (C1E)" selectedOption="Enable" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
              </AvailableOptions>
              <DefaultOption>Enable</DefaultOption>
              <Help><![CDATA[Core C1E auto promotion Control. Takes effect after reboot.]]></Help>
              <WorkIf><![CDATA[  1 != Autonomous Core C-State  ]]></WorkIf>
            </Information>
          </Setting>
        </Menu>
        <Menu name="Package C State Control">
          <Information>
            <Help><![CDATA[Package C State setting]]></Help>
            <WorkIf><![CDATA[  2 == Power Technology  ]]></WorkIf>
          </Information>
          <Subtitle>Package C State Control</Subtitle>
          <Subtitle></Subtitle>
          <Setting name="Package C State" selectedOption="Auto" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">C0/C1 state</Option>
                <Option value="1">C2 state</Option>
                <Option value="2">C6(non Retention) state</Option>
                <Option value="3">C6(Retention) state</Option>
                <Option value="7">No Limit</Option>
                <Option value="255">Auto</Option>
              </AvailableOptions>
              <DefaultOption>Auto</DefaultOption>
              <Help><![CDATA[Package C State limit]]></Help>
            </Information>
          </Setting>
        </Menu>
        <Menu name="CPU T State Control">
          <Information>
            <Help><![CDATA[CPU T State setting]]></Help>
            <WorkIf><![CDATA[  2 == Power Technology  ]]></WorkIf>
          </Information>
          <Subtitle>CPU T State Control</Subtitle>
          <Subtitle></Subtitle>
          <Setting name="Software Controlled T-States" selectedOption="Disable" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
              </AvailableOptions>
              <DefaultOption>Disable</DefaultOption>
              <Help><![CDATA[Enable/Disable Software Controlled T-States]]></Help>
            </Information>
          </Setting>
        </Menu>
      </Menu>
    </Menu>
    <Menu name="Chipset Configuration">
      <Information>
        <Help><![CDATA[System Chipset configuration.]]></Help>
      </Information>
      <Subtitle>WARNING: Setting wrong values in below sections may cause</Subtitle>
      <Subtitle>         system to malfunction.</Subtitle>
      <Menu name="North Bridge">
        <Information>
          <Help><![CDATA[North Bridge Parameters]]></Help>
        </Information>
        <Menu name="Memory Configuration">
          <Information>
            <Help><![CDATA[Displays and provides option to change the Memory Settings]]></Help>
          </Information>
          <Subtitle></Subtitle>
          <Subtitle>--------------------------------------------------</Subtitle>
          <Subtitle>Integrated Memory Controller (iMC)</Subtitle>
          <Subtitle>--------------------------------------------------</Subtitle>
          <Subtitle></Subtitle>
          <Setting name="Enforce POR" selectedOption="POR" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">POR</Option>
                <Option value="2">Disable</Option>
              </AvailableOptions>
              <DefaultOption>POR</DefaultOption>
              <Help><![CDATA[Enable - Enforces Plan Of Record restrictions for DDR4 frequency and voltage programming. Disable - Disables this feature. Auto - Sets it to the MRC default setting; current default is Enable.]]></Help>
            </Information>
          </Setting>
          <Setting name="Memory Frequency" selectedOption="Auto" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Auto</Option>
                <Option value="11">2133</Option>
                <Option value="13">2400</Option>
                <Option value="15">2666</Option>
              </AvailableOptions>
              <DefaultOption>Auto</DefaultOption>
              <Help><![CDATA[Maximum Memory Frequency Selections in Mhz. Do not select Reserved]]></Help>
            </Information>
          </Setting>
          <Setting name="Data Scrambling for DDR4" selectedOption="Auto" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="2">Auto</Option>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
              </AvailableOptions>
              <DefaultOption>Auto</DefaultOption>
              <Help><![CDATA[Enable - Enables data scrambling for DDR4. Disable - Disables this feature. Auto - Sets it to the MRC default setting; current default is Enable.]]></Help>
            </Information>
          </Setting>
          <Setting name="tCCD_L Relaxation" selectedOption="Auto" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Auto</Option>
              </AvailableOptions>
              <DefaultOption>Auto</DefaultOption>
              <Help><![CDATA[Auto: tCCD_L Relaxation will be enabled for selected & capable DIMMs. Disable: tCCD_L Relaxation is always disabled.]]></Help>
            </Information>
          </Setting>
          <Setting name="2X REFRESH" selectedOption="Auto" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Auto</Option>
                <Option value="1">Enable</Option>
              </AvailableOptions>
              <DefaultOption>Auto</DefaultOption>
              <Help><![CDATA[Select 2x refresh mode]]></Help>
            </Information>
          </Setting>
          <Menu name="Memory Topology">
            <Information>
              <Help><![CDATA[Displays memory topology with Dimm population information]]></Help>
            </Information>
            <Subtitle>--------------------------------------------------</Subtitle>
            <Subtitle>DIMMA1:  2132MT/s Micron DRx4 32GB RDIMM</Subtitle>
            <Subtitle>DIMMB1:  2132MT/s Micron DRx4 32GB RDIMM</Subtitle>
            <Subtitle>DIMMD1:  2132MT/s Micron DRx4 32GB RDIMM</Subtitle>
            <Subtitle>DIMME1:  2132MT/s Micron DRx4 32GB RDIMM</Subtitle>
            <Subtitle></Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
          </Menu>
          <Menu name="Memory RAS Configuration">
            <Information>
              <Help><![CDATA[Displays and provides option to change the Memory Ras Settings]]></Help>
            </Information>
            <Subtitle></Subtitle>
            <Subtitle>--------------------------------------------------</Subtitle>
            <Subtitle>Memory RAS Configuration Setup</Subtitle>
            <Subtitle>--------------------------------------------------</Subtitle>
            <Subtitle></Subtitle>
            <Setting name="Static Virtual Lockstep Mode" selectedOption="Disable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">Disable</Option>
                  <Option value="3">Enable</Option>
                </AvailableOptions>
                <DefaultOption>Disable</DefaultOption>
                <Help><![CDATA[Enable Static Virtual Lockstep mode]]></Help>
              </Information>
            </Setting>
            <Setting name="Mirror mode" selectedOption="Disable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">Disable</Option>
                  <Option value="1">Enable Mirror Mode (1LM)</Option>
                </AvailableOptions>
                <DefaultOption>Disable</DefaultOption>
                <Help><![CDATA[Mirror Mode will set entire 1LM memory in system to be mirrored, consequently reducing the memory capacity by half. Mirror Enable will disable XPT Prefetch]]></Help>
                <WorkIf><![CDATA[   (  (  ( 1 != 0 )  &&  ( 1 != 0 )  )  &&  ( 1 != 0 )  )  &&  ( 1 != ADDDC Sparing )   ]]></WorkIf>
              </Information>
            </Setting>
            <Setting name="Memory Rank Sparing" selectedOption="Disable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">Disable</Option>
                  <Option value="1">Enable</Option>
                </AvailableOptions>
                <DefaultOption>Disable</DefaultOption>
                <Help><![CDATA[Enable/Disable Memory Rank Sparing. This feature is only available on 1LM]]></Help>
                <WorkIf><![CDATA[   ( 1 != Mirror mode )  &&  ( 1 != Volatile Memory Mode )   ]]></WorkIf>
              </Information>
            </Setting>
            <Setting name="Multi Rank Sparing" selectedOption="Two Rank" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="1">One Rank</Option>
                  <Option value="2">Two Rank</Option>
                </AvailableOptions>
                <DefaultOption>Two Rank</DefaultOption>
                <Help><![CDATA[Set Multi Rank Sparing number, default and the maximum is 2 ranks per channel]]></Help>
                <WorkIf><![CDATA[  0 != Memory Rank Sparing  ]]></WorkIf>
              </Information>
            </Setting>
            <Setting name="Correctable Error Threshold" numericValue="100" type="Numeric">
              <Information>
                <MaxValue>32767</MaxValue>
                <MinValue>0</MinValue>
                <StepSize>1</StepSize>
                <DefaultValue>100</DefaultValue>
                <Help><![CDATA[Correctable Error Threshold (1 - 32767) used for sparing, tagging, and leaky bucket]]></Help>
              </Information>
            </Setting>
            <Setting name="SDDC" selectedOption="Disable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">Disable</Option>
                  <Option value="1">Enable</Option>
                </AvailableOptions>
                <DefaultOption>Disable</DefaultOption>
                <Help><![CDATA[Enable/Disable SDDC. Not supported when AEP dimm present!]]></Help>
              </Information>
            </Setting>
            <Setting name="ADDDC Sparing" selectedOption="Disable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">Disable</Option>
                  <Option value="1">Enable</Option>
                </AvailableOptions>
                <DefaultOption>Disable</DefaultOption>
                <Help><![CDATA[Enable/Disable ADDDC Sparing]]></Help>
                <WorkIf><![CDATA[   (  (  ( 1 != 0 )  &&  ( 1 != Mirror mode )  )  &&  ( 1 != 0 )  )  &&  ( 1 != 0 )   ]]></WorkIf>
              </Information>
            </Setting>
            <Setting name="Patrol Scrub" selectedOption="Enable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">Disable</Option>
                  <Option value="1">Enable</Option>
                </AvailableOptions>
                <DefaultOption>Enable</DefaultOption>
                <Help><![CDATA[Enable/Disable Patrol Scrub]]></Help>
              </Information>
            </Setting>
            <Setting name="Patrol Scrub Interval" numericValue="24" type="Numeric">
              <Information>
                <MaxValue>24</MaxValue>
                <MinValue>0</MinValue>
                <StepSize>0</StepSize>
                <DefaultValue>24</DefaultValue>
                <Help><![CDATA[Selects the number of hours (1-24) required to complete full scrub. A value of zero means auto!]]></Help>
                <WorkIf><![CDATA[  0 != Patrol Scrub  ]]></WorkIf>
              </Information>
            </Setting>
          </Menu>
        </Menu>
        <Menu name="IIO Configuration">
          <Information>
            <Help><![CDATA[Displays and provides option to change the IIO Settings]]></Help>
          </Information>
          <Subtitle>IIO Configuration</Subtitle>
          <Subtitle>--------------------------------------------------</Subtitle>
          <Subtitle></Subtitle>
          <Setting name="EV DFX Features" selectedOption="Disable" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
              </AvailableOptions>
              <DefaultOption>Disable</DefaultOption>
              <Help><![CDATA[Expose IIO DFX devices and other CPU devices like PMON]]></Help>
            </Information>
          </Setting>
          <Menu name="CPU Configuration" order="2">
            <Information>
              <Help><![CDATA[]]></Help>
            </Information>
            <Setting name="IOU0 (IIO PCIe Br1)" selectedOption="Auto" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">x4x4x4x4</Option>
                  <Option value="1">x4x4x8</Option>
                  <Option value="2">x8x4x4</Option>
                  <Option value="3">x8x8</Option>
                  <Option value="4">x16</Option>
                  <Option value="255">Auto</Option>
                </AvailableOptions>
                <DefaultOption>Auto</DefaultOption>
                <Help><![CDATA[Selects PCIe port Bifurcation for selected slot(s)]]></Help>
              </Information>
            </Setting>
            <Setting name="IOU1 (IIO PCIe Br2)" selectedOption="Auto" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">x4x4x4x4</Option>
                  <Option value="1">x4x4x8</Option>
                  <Option value="255">Auto</Option>
                </AvailableOptions>
                <DefaultOption>Auto</DefaultOption>
                <Help><![CDATA[Selects PCIe port Bifurcation for selected slot(s)]]></Help>
              </Information>
            </Setting>
            <Menu name="CPU SLOT6 PCI-E 3.0 X16">
              <Information>
                <Help><![CDATA[Settings related to PCI Express PortS (0/1A/1B/1C/1D/2A/2B/2C/2D/3A/3B/3C/3D/4A/5A) ]]></Help>
              </Information>
              <Subtitle>CPU SLOT6 PCI-E 3.0 X16</Subtitle>
              <Subtitle>--------------------------------------------------</Subtitle>
              <Subtitle></Subtitle>
              <Setting name="Link Speed" order="1" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Auto</Option>
                    <Option value="1">Gen 1 (2.5 GT/s)</Option>
                    <Option value="2">Gen 2 (5 GT/s)</Option>
                    <Option value="3">Gen 3 (8 GT/s)</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Choose Link Speed for this PCIe port]]></Help>
                </Information>
              </Setting>
              <Text>PCI-E Port Link Status(Link Did Not Train)</Text>
              <Text>PCI-E Port Link Max(Max Width x16)</Text>
              <Text>PCI-E Port Link Speed(Link Did Not Train)</Text>
              <Setting name="PCI-E Port Max Payload Size" order="1" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">128B</Option>
                    <Option value="1">256B</Option>
                    <Option value="2">Auto</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Set Maxpayload size to 256B if possible]]></Help>
                </Information>
              </Setting>
            </Menu>
            <Menu name="CPU SLOT7 PCI-E 3.0 X8">
              <Information>
                <Help><![CDATA[Settings related to PCI Express PortS (0/1A/1B/1C/1D/2A/2B/2C/2D/3A/3B/3C/3D/4A/5A) ]]></Help>
              </Information>
              <Subtitle>CPU SLOT7 PCI-E 3.0 X8</Subtitle>
              <Subtitle>--------------------------------------------------</Subtitle>
              <Subtitle></Subtitle>
              <Setting name="Link Speed" order="2" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Auto</Option>
                    <Option value="1">Gen 1 (2.5 GT/s)</Option>
                    <Option value="2">Gen 2 (5 GT/s)</Option>
                    <Option value="3">Gen 3 (8 GT/s)</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Choose Link Speed for this PCIe port]]></Help>
                </Information>
              </Setting>
              <Text>PCI-E Port Link Status(Link Did Not Train)</Text>
              <Text>PCI-E Port Link Max(Max Width x8)</Text>
              <Text>PCI-E Port Link Speed(Link Did Not Train)</Text>
              <Setting name="PCI-E Port Max Payload Size" order="2" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">128B</Option>
                    <Option value="1">256B</Option>
                    <Option value="2">Auto</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Set Maxpayload size to 256B if possible]]></Help>
                </Information>
              </Setting>
            </Menu>
          </Menu>
          <Menu name="IOAT Configuration">
            <Information>
              <Help><![CDATA[All IOAT configuration options]]></Help>
            </Information>
            <Setting name="Disable TPH" selectedOption="No" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">No</Option>
                  <Option value="1">Yes</Option>
                </AvailableOptions>
                <DefaultOption>No</DefaultOption>
                <Help><![CDATA[TLP Processing Hint disable]]></Help>
              </Information>
            </Setting>
            <Setting name="Prioritize TPH" selectedOption="Disable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="1">Enable</Option>
                  <Option value="0">Disable</Option>
                </AvailableOptions>
                <DefaultOption>Disable</DefaultOption>
                <Help><![CDATA[Prioritize TPH]]></Help>
                <WorkIf><![CDATA[  1 != Disable TPH  ]]></WorkIf>
              </Information>
            </Setting>
            <Setting name="Relaxed Ordering" selectedOption="Disable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">Disable</Option>
                  <Option value="1">Enable</Option>
                </AvailableOptions>
                <DefaultOption>Disable</DefaultOption>
                <Help><![CDATA[Relaxed Ordering Enable/Disable]]></Help>
              </Information>
            </Setting>
          </Menu>
          <Menu name="Intel® VT for Directed I/O (VT-d)" order="1">
            <Information>
              <Help><![CDATA[Press <Enter> to bring up the Intel® VT for Directed I/O (VT-d) Configuration menu.]]></Help>
            </Information>
            <Subtitle>Intel® VT for Directed I/O (VT-d)</Subtitle>
            <Subtitle>--------------------------------------------------</Subtitle>
            <Subtitle></Subtitle>
            <Setting name="Intel® VT for Directed I/O (VT-d)" order="2" selectedOption="Enable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="1">Enable</Option>
                  <Option value="0">Disable</Option>
                </AvailableOptions>
                <DefaultOption>Enable</DefaultOption>
                <Help><![CDATA[Enable/Disable Intel® Virtualization Technology for Directed I/O (VT-d) by reporting the I/O device assignment to VMM through DMAR ACPI Tables.]]></Help>
              </Information>
            </Setting>
            <Setting name="Interrupt Remapping" selectedOption="Enable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="1">Enable</Option>
                  <Option value="0">Disable</Option>
                </AvailableOptions>
                <DefaultOption>Enable</DefaultOption>
                <Help><![CDATA[Enable/Disable VT_D Interrupt Remapping Support]]></Help>
                <WorkIf><![CDATA[  0 != Intel® VT for Directed I/O (VT-d)$2  ]]></WorkIf>
              </Information>
            </Setting>
            <Setting name="PassThrough DMA" selectedOption="Enable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="1">Enable</Option>
                  <Option value="0">Disable</Option>
                </AvailableOptions>
                <DefaultOption>Enable</DefaultOption>
                <Help><![CDATA[Enable/Disable Non-Isoch VT_D Engine Pass Through DMA support]]></Help>
                <WorkIf><![CDATA[  0 != Intel® VT for Directed I/O (VT-d)$2  ]]></WorkIf>
              </Information>
            </Setting>
            <Setting name="ATS" selectedOption="Enable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="1">Enable</Option>
                  <Option value="0">Disable</Option>
                </AvailableOptions>
                <DefaultOption>Enable</DefaultOption>
                <Help><![CDATA[Enable/Disable Non-Isoch VT_D Engine ATS support]]></Help>
                <WorkIf><![CDATA[  0 != Intel® VT for Directed I/O (VT-d)$2  ]]></WorkIf>
              </Information>
            </Setting>
            <Setting name="Posted Interrupt" selectedOption="Enable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="1">Enable</Option>
                  <Option value="0">Disable</Option>
                </AvailableOptions>
                <DefaultOption>Enable</DefaultOption>
                <Help><![CDATA[Enable/Disable VT_D posted interrupt]]></Help>
                <WorkIf><![CDATA[  0 != Intel® VT for Directed I/O (VT-d)$2  ]]></WorkIf>
              </Information>
            </Setting>
            <Setting name="Coherency Support (Non-Isoch)" selectedOption="Enable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="1">Enable</Option>
                  <Option value="0">Disable</Option>
                </AvailableOptions>
                <DefaultOption>Enable</DefaultOption>
                <Help><![CDATA[Enable/Disable Non-Isoch VT_D Engine Coherency support]]></Help>
                <WorkIf><![CDATA[  0 != Intel® VT for Directed I/O (VT-d)$2  ]]></WorkIf>
              </Information>
            </Setting>
          </Menu>
          <Menu name="Intel® VMD technology">
            <Information>
              <Help><![CDATA[Press <Enter> to bring up the Intel® VMD for Volume Management Device Configuration menu.]]></Help>
            </Information>
            <Subtitle>Intel® VMD technology</Subtitle>
            <Subtitle>--------------------------------------------------</Subtitle>
            <Subtitle></Subtitle>
            <Menu name="Intel® VMD for Volume Management Device on CPU">
              <Information>
                <Help><![CDATA[]]></Help>
              </Information>
              <Subtitle>VMD Config for PStack0</Subtitle>
              <Subtitle>--------------------------------------------------</Subtitle>
              <Setting name="Intel® VMD for Volume Management Device for PStack0" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology in this Stack.]]></Help>
                </Information>
              </Setting>
              <Setting name="CPU SLOT6 PCI-E 3.0 X16 VMD" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology on specific root port]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for PStack0  ]]></WorkIf>
                </Information>
              </Setting>
              <Setting name="Hot Plug Capable" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Hot Plug for PCIe Root Ports 1A-1D]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for PStack0  ]]></WorkIf>
                </Information>
              </Setting>
              <Subtitle></Subtitle>
              <!--Valid if:   0 != Intel® VMD for Volume Management Device for PStack0  -->
            </Menu>
          </Menu>
          <Subtitle></Subtitle>
          <Subtitle> IIO-PCIE Express Global Options</Subtitle>
          <Subtitle>========================================</Subtitle>
          <Setting name="PCI-E Completion Timeout Disable" selectedOption="No" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="1">Yes</Option>
                <Option value="0">No</Option>
                <Option value="2">Per-Port</Option>
              </AvailableOptions>
              <DefaultOption>No</DefaultOption>
              <Help><![CDATA[Enable / disable the Completion Timeout (D:x F:0 O:B8h B:4) where x is 0-3]]></Help>
            </Information>
          </Setting>
        </Menu>
      </Menu>
      <Menu name="South Bridge">
        <Information>
          <Help><![CDATA[South Bridge Parameters]]></Help>
        </Information>
        <Subtitle></Subtitle>
        <Text>USB Module Version(20)</Text>
        <Subtitle></Subtitle>
        <Text>USB Devices:()</Text>
        <Subtitle>      2 Keyboards, 2 Mice, 4 Hubs</Subtitle>
        <Subtitle></Subtitle>
        <Setting name="Legacy USB Support" selectedOption="Enabled" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">Enabled</Option>
              <Option value="1">Disabled</Option>
              <Option value="2">Auto</Option>
            </AvailableOptions>
            <DefaultOption>Enabled</DefaultOption>
            <Help><![CDATA[Enables Legacy USB support. AUTO option disables legacy support if no USB devices are connected. DISABLE option will keep USB devices available only for EFI applications.]]></Help>
          </Information>
        </Setting>
        <Setting name="XHCI Hand-off" selectedOption="Disabled" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="1">Enabled</Option>
              <Option value="0">Disabled</Option>
            </AvailableOptions>
            <DefaultOption>Disabled</DefaultOption>
            <Help><![CDATA[This is a workaround for OSes without XHCI hand-off support. The XHCI ownership change should be claimed by XHCI driver.]]></Help>
          </Information>
        </Setting>
        <Setting name="Port 60/64 Emulation" selectedOption="Enabled" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">Disabled</Option>
              <Option value="1">Enabled</Option>
            </AvailableOptions>
            <DefaultOption>Enabled</DefaultOption>
            <Help><![CDATA[Enables I/O port 60h/64h emulation support. This should be enabled for the complete USB keyboard legacy support for non-USB aware OSes.]]></Help>
          </Information>
        </Setting>
      </Menu>
    </Menu>
    <Menu name="Server ME Configuration">
      <Information>
        <Help><![CDATA[Configure Server ME Technology Parameters]]></Help>
      </Information>
      <Subtitle>General ME Configuration</Subtitle>
      <Text>Oper. Firmware Version(0E:4.0.4.97)</Text>
      <Text>Backup Firmware Version(N/A)</Text>
      <Text>Recovery Firmware Version(0E:4.0.4.97)</Text>
      <Text>ME Firmware Status #1(0x000F0245)</Text>
      <Text>ME Firmware Status #2(0x8811C026)</Text>
      <Text>  Current State(Operational)</Text>
      <Text>  Error Code(No Error)</Text>
    </Menu>
    <Menu name="PCH SATA Configuration">
      <Information>
        <Help><![CDATA[SATA devices and settings]]></Help>
      </Information>
      <Subtitle>PCH SATA Configuration</Subtitle>
      <Subtitle>--------------------------------------------------</Subtitle>
      <Subtitle></Subtitle>
      <Setting name="SATA Controller" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Enable or Disable SATA Controller]]></Help>
        </Information>
      </Setting>
      <Setting name="Configure SATA as" selectedOption="AHCI" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">AHCI</Option>
            <Option value="1">RAID</Option>
          </AvailableOptions>
          <DefaultOption>AHCI</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
          <WorkIf><![CDATA[  0 != SATA Controller  ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="SATA HDD Unlock" order="1" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Enable: HDD password unlock is enabled in the OS]]></Help>
          <WorkIf><![CDATA[  0 != SATA Controller  ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="SATA RSTe Boot Info" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Enable setting provides full int13h support for SATA controller attached devices. CSM storage OPROM policy should be set to legacy to make this selection effective.]]></Help>
          <WorkIf><![CDATA[ ( 0 != SATA Controller )  and  ( 1 == Configure SATA as ) ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="Aggressive Link Power Management" order="1" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[Enables/Disables SALP]]></Help>
        </Information>
      </Setting>
      <Setting name="SATA RAID Option ROM/UEFI Driver" selectedOption="Legacy" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">EFI</Option>
            <Option value="2">Legacy</Option>
          </AvailableOptions>
          <DefaultOption>Legacy</DefaultOption>
          <Help><![CDATA[In RAID mode load EFI driver. (If disabled loads LEGACY OPROM)]]></Help>
          <WorkIf><![CDATA[  1 == Configure SATA as  ]]></WorkIf>
        </Information>
      </Setting>
      <Subtitle></Subtitle>
      <Text>SATA Port 0(TOSHIBA MG06AC - 10000.8 GB)</Text>
      <Setting name="  Hot Plug" order="1" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="1" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  SATA Device Type" order="1" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>SATA Port 1(TOSHIBA MG06AC - 10000.8 GB)</Text>
      <Setting name="  Hot Plug" order="2" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="2" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  SATA Device Type" order="2" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>SATA Port 2(TOSHIBA MG06AC - 10000.8 GB)</Text>
      <Setting name="  Hot Plug" order="3" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="3" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  SATA Device Type" order="3" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>SATA Port 3(TOSHIBA MG06AC - 10000.8 GB)</Text>
      <Setting name="  Hot Plug" order="4" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="4" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  SATA Device Type" order="4" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>SATA Port 4(TOSHIBA MG06AC - 10000.8 GB)</Text>
      <Setting name="  Hot Plug" order="5" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="5" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  SATA Device Type" order="5" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>SATA Port 5(TOSHIBA MG06AC - 10000.8 GB)</Text>
      <Setting name="  Hot Plug" order="6" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="6" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  SATA Device Type" order="6" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>SATA Port 6(TOSHIBA MG06AC - 10000.8 GB)</Text>
      <Setting name="  Hot Plug" order="7" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="7" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  SATA Device Type" order="7" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>SATA Port 7(TOSHIBA MG06AC - 10000.8 GB)</Text>
      <Setting name="  Hot Plug" order="8" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="8" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  SATA Device Type" order="8" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
    </Menu>
    <Menu name="PCH sSATA Configuration">
      <Information>
        <Help><![CDATA[sSATA devices and settings]]></Help>
      </Information>
      <Subtitle>PCH sSATA Configuration</Subtitle>
      <Subtitle>--------------------------------------------------</Subtitle>
      <Subtitle></Subtitle>
      <Setting name="sSATA Controller" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Enable</Option>
            <Option value="0">Disable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Enable or Disable SATA Controller]]></Help>
        </Information>
      </Setting>
      <Setting name="Configure sSATA as" selectedOption="AHCI" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">AHCI</Option>
            <Option value="1">RAID</Option>
          </AvailableOptions>
          <DefaultOption>AHCI</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
          <WorkIf><![CDATA[  0 != sSATA Controller  ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="SATA HDD Unlock" order="2" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Enable: HDD password unlock is enabled in the OS]]></Help>
          <WorkIf><![CDATA[  0 != sSATA Controller  ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="sSATA RSTe Boot Info" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Enable setting provides full int13h support for SATA controller attached devices. CSM storage OPROM policy should be set to legacy to make this selection effective.]]></Help>
          <WorkIf><![CDATA[ ( 0 != sSATA Controller )  and  ( 1 == Configure sSATA as ) ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="Aggressive Link Power Management" order="2" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[Enables/Disables SALP]]></Help>
        </Information>
      </Setting>
      <Setting name="sSATA RAID Option ROM/UEFI Driver" selectedOption="Legacy" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">EFI</Option>
            <Option value="2">Legacy</Option>
          </AvailableOptions>
          <DefaultOption>Legacy</DefaultOption>
          <Help><![CDATA[In RAID mode load EFI driver. (If disabled loads LEGACY OPROM)]]></Help>
          <WorkIf><![CDATA[  1 == Configure sSATA as  ]]></WorkIf>
        </Information>
      </Setting>
      <Subtitle></Subtitle>
      <Text>sSATA Port 0(TOSHIBA MG06AC - 10000.8 GB)</Text>
      <Setting name="  Hot Plug" order="9" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="9" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  sSATA Device Type" order="1" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>sSATA Port 1(TOSHIBA MG06AC - 10000.8 GB)</Text>
      <Setting name="  Hot Plug" order="10" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="10" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  sSATA Device Type" order="2" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>sSATA Port 2(TOSHIBA MG06AC - 10000.8 GB)</Text>
      <Setting name="  Hot Plug" order="11" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="11" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  sSATA Device Type" order="3" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>sSATA Port 3(TOSHIBA MG06AC - 10000.8 GB)</Text>
      <Setting name="  Hot Plug" order="12" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="12" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  sSATA Device Type" order="4" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>sSATA Port 4(TS64GMTS400S   - 64.0 GB)</Text>
      <Setting name="  Hot Plug" order="13" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="13" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  sSATA Device Type" order="5" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>sSATA Port 5([Not Installed])</Text>
      <Setting name="  Hot Plug" order="14" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="14" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  sSATA Device Type" order="6" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
    </Menu>
    <Menu name="AMI Graphic Output Protocol Policy">
      <Information>
        <Help><![CDATA[User Select Monitor Output by Graphic Output Protocol]]></Help>
      </Information>
      <Subtitle>P3:TOSHIBA MG06ACA10TE</Subtitle>
      <Subtitle>P7:TOSHIBA MG06ACA10TE</Subtitle>
      <Setting name="Output Select" selectedOption="Unknown Device" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Unknown Device</Option>
          </AvailableOptions>
          <DefaultOption>Unknown Device</DefaultOption>
          <Help><![CDATA[Output Interface]]></Help>
        </Information>
      </Setting>
    </Menu>
    <Menu name="PCIe/PCI/PnP Configuration">
      <Information>
        <Help><![CDATA[PCI, PCI-X and PCI Express Settings.]]></Help>
      </Information>
      <Text>PCI Bus Driver Version(A5.01.16)</Text>
      <Subtitle></Subtitle>
      <Subtitle>PCI Devices Common Settings:</Subtitle>
      <Setting name="Above 4G Decoding" selectedOption="Enabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Enabled</DefaultOption>
          <Help><![CDATA[Enables or Disables 64bit capable Devices to be Decoded in Above 4G Address Space (Only if System Supports 64 bit PCI Decoding).]]></Help>
        </Information>
      </Setting>
      <Setting name="SR-IOV Support" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[If system has SR-IOV capable PCIe Devices, this option Enables or Disables Single Root IO Virtualization Support.]]></Help>
        </Information>
      </Setting>
      <Setting name="BME DMA Mitigation" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[Re-enable Bus Master Attribute disabled during Pci enumeration for PCI Bridges after SMM Locked ]]></Help>
        </Information>
      </Setting>
      <Setting name="MMIO High Base" selectedOption="56T" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">56T</Option>
            <Option value="1">40T</Option>
            <Option value="2">24T</Option>
            <Option value="3">16T</Option>
            <Option value="4">4T</Option>
            <Option value="5">1T</Option>
          </AvailableOptions>
          <DefaultOption>56T</DefaultOption>
          <Help><![CDATA[Select MMIO High Base]]></Help>
        </Information>
      </Setting>
      <Setting name="MMIO High Granularity Size" selectedOption="256G" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">1G</Option>
            <Option value="1">4G</Option>
            <Option value="2">16G</Option>
            <Option value="3">64G</Option>
            <Option value="4">256G</Option>
            <Option value="5">1024G</Option>
          </AvailableOptions>
          <DefaultOption>256G</DefaultOption>
          <Help><![CDATA[Selects the allocation size used to assign mmioh resources.
Total mmioh space can be up to 32xgranularity.
Per stack mmioh resource assignments are multiples of the granularity where 1 unit per stack is the default allocation.]]></Help>
        </Information>
      </Setting>
      <Setting name="Maximum Read Request" selectedOption="Auto" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="55">Auto</Option>
            <Option value="0">128 Bytes</Option>
            <Option value="1">256 Bytes</Option>
            <Option value="2">512 Bytes</Option>
            <Option value="3">1024 Bytes</Option>
            <Option value="4">2048 Bytes</Option>
            <Option value="5">4096 Bytes</Option>
          </AvailableOptions>
          <DefaultOption>Auto</DefaultOption>
          <Help><![CDATA[Set Maximum Read Request Size of PCI Express Device or allow System BIOS to select the value.]]></Help>
        </Information>
      </Setting>
      <Setting name="MMCFG BASE" selectedOption="2G" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">1G</Option>
            <Option value="1">1.5G</Option>
            <Option value="2">1.75G</Option>
            <Option value="3">2G</Option>
            <Option value="4">2.25G</Option>
            <Option value="5">3G</Option>
          </AvailableOptions>
          <DefaultOption>2G</DefaultOption>
          <Help><![CDATA[Select MMCFG Base]]></Help>
        </Information>
      </Setting>
      <Setting name="NVMe Firmware Source" selectedOption="Vendor Defined Firmware" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Vendor Defined Firmware</Option>
            <Option value="1">AMI Native Support</Option>
          </AvailableOptions>
          <DefaultOption>Vendor Defined Firmware</DefaultOption>
          <Help><![CDATA[AMI Native FW Support or Device Vendor Defined FW Support]]></Help>
        </Information>
      </Setting>
      <Setting name="VGA Priority" selectedOption="Onboard" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Onboard</Option>
            <Option value="2">Offboard</Option>
          </AvailableOptions>
          <DefaultOption>Onboard</DefaultOption>
          <Help><![CDATA[Select active Video type]]></Help>
        </Information>
      </Setting>
      <Setting name="Primary PCIE VGA" selectedOption="CPU SLOT6 PCI-E 3.0 X16" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">CPU SLOT6 PCI-E 3.0 X16</Option>
            <Option value="2">CPU SLOT7 PCI-E 3.0 X8</Option>
          </AvailableOptions>
          <DefaultOption>CPU SLOT6 PCI-E 3.0 X16</DefaultOption>
          <Help><![CDATA[Select the primary PCIE VGA.]]></Help>
          <WorkIf><![CDATA[  1 != VGA Priority  ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="Consistent Device Name Support" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[Enable ACPI _DSM device name support for onbard devices and slots.]]></Help>
        </Information>
      </Setting>
      <Setting name="JMD2:M.2-H PCI-E 3.0 X2 lane 1 Type" selectedOption="PCIE" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">PCIE</Option>
            <Option value="1">USB 3.0</Option>
          </AvailableOptions>
          <DefaultOption>PCIE</DefaultOption>
          <Help><![CDATA[Select the IO type of JMD2:M.2-H PCI-E 3.0 X2 lane 1. ]]></Help>
        </Information>
      </Setting>
      <Subtitle>RSC-RR1U-E8</Subtitle>
      <Setting name="CPU SLOT6 PCI-E 3.0 X16 OPROM" selectedOption="Legacy" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Legacy</Option>
            <Option value="2">EFI</Option>
          </AvailableOptions>
          <DefaultOption>Legacy</DefaultOption>
          <Help><![CDATA[Enables or disables CPU SLOT6 PCI-E 3.0 X16 OPROM option.]]></Help>
        </Information>
      </Setting>
      <Setting name="CPU SLOT7 PCI-E 3.0 X8 OPROM" selectedOption="Legacy" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Legacy</Option>
            <Option value="2">EFI</Option>
          </AvailableOptions>
          <DefaultOption>Legacy</DefaultOption>
          <Help><![CDATA[Enables or disables CPU SLOT7 PCI-E 3.0 X8 OPROM option.]]></Help>
        </Information>
      </Setting>
      <Setting name="JMD1:M.2-HC PCI-E 3.0 X4 OPROM" selectedOption="Legacy" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Legacy</Option>
            <Option value="2">EFI</Option>
          </AvailableOptions>
          <DefaultOption>Legacy</DefaultOption>
          <Help><![CDATA[Enables or disables JMD1:M.2-HC PCI-E 3.0 X4 OPROM option.]]></Help>
        </Information>
      </Setting>
      <Setting name="JMD2:M.2-H PCI-E 3.0 X2 OPROM" selectedOption="Legacy" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Legacy</Option>
            <Option value="2">EFI</Option>
          </AvailableOptions>
          <DefaultOption>Legacy</DefaultOption>
          <Help><![CDATA[Enables or disables JMD2:M.2-H PCI-E 3.0 X2 OPROM option.]]></Help>
        </Information>
      </Setting>
      <Setting name="PCI-E 3.0 X1 OPROM" selectedOption="Legacy" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Legacy</Option>
            <Option value="2">EFI</Option>
          </AvailableOptions>
          <DefaultOption>Legacy</DefaultOption>
          <Help><![CDATA[Enables or disables mini PCI-E 3.0 X1 OPROM option.]]></Help>
        </Information>
      </Setting>
      <Setting name="Onboard LAN Option ROM Type" selectedOption="Legacy" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Legacy</Option>
            <Option value="1">EFI</Option>
          </AvailableOptions>
          <DefaultOption>Legacy</DefaultOption>
          <Help><![CDATA[Select which firmware type to be loaded for onboard LANs]]></Help>
        </Information>
      </Setting>
      <Setting name="Onboard LAN1 Option ROM" selectedOption="PXE" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">PXE</Option>
            <Option value="2">iSCSI</Option>
          </AvailableOptions>
          <DefaultOption>PXE</DefaultOption>
          <Help><![CDATA[Select which firmware function to be loaded for onboard LAN1.]]></Help>
          <WorkIf><![CDATA[  1 != Onboard LAN Option ROM Type  ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="Onboard LAN2 Option ROM" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">PXE</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[Select which firmware function to be loaded for onboard LAN2]]></Help>
          <WorkIf><![CDATA[ ( 1 != Onboard LAN Option ROM Type )  and  ( Onboard LAN1 Option ROM <= 1 ) ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="Onboard LAN3 Option ROM" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">PXE</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[Select which firmware function to be loaded for onboard LAN3]]></Help>
          <WorkIf><![CDATA[ ( 1 != Onboard LAN Option ROM Type )  and  ( Onboard LAN1 Option ROM <= 1 ) ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="Onboard LAN4 Option ROM" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">PXE</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[Select which firmware function to be loaded for onboard LAN4]]></Help>
          <WorkIf><![CDATA[ ( 1 != Onboard LAN Option ROM Type )  and  ( Onboard LAN1 Option ROM <= 1 ) ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="Onboard LAN5 Option ROM" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Legacy</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[Select which firmware function to be loaded for onboard LAN5]]></Help>
          <WorkIf><![CDATA[  1 != Onboard LAN Option ROM Type  ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="Onboard LAN6 Option ROM" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Legacy</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[Select which firmware function to be loaded for onboard LAN6]]></Help>
          <WorkIf><![CDATA[  1 != Onboard LAN Option ROM Type  ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="Onboard LAN7 Option ROM" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Legacy</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[Select which firmware function to be loaded for onboard LAN7]]></Help>
          <WorkIf><![CDATA[  1 != Onboard LAN Option ROM Type  ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="Onboard LAN8 Option ROM" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Legacy</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[Select which firmware function to be loaded for onboard LAN8]]></Help>
          <WorkIf><![CDATA[  1 != Onboard LAN Option ROM Type  ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="Onboard Video Option ROM" selectedOption="Legacy" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Legacy</Option>
            <Option value="2">EFI</Option>
          </AvailableOptions>
          <DefaultOption>Legacy</DefaultOption>
          <Help><![CDATA[Select which onboard video firmware type to be loaded.]]></Help>
        </Information>
      </Setting>
      <Menu name="Network Stack Configuration">
        <Information>
          <Help><![CDATA[Network Stack Settings]]></Help>
        </Information>
        <Setting name="Network Stack" selectedOption="Enabled" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">Disabled</Option>
              <Option value="1">Enabled</Option>
            </AvailableOptions>
            <DefaultOption>Enabled</DefaultOption>
            <Help><![CDATA[Enable/Disable UEFI Network Stack]]></Help>
          </Information>
        </Setting>
        <Setting name="Ipv4 PXE Support" selectedOption="Enabled" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">Disabled</Option>
              <Option value="1">Enabled</Option>
            </AvailableOptions>
            <DefaultOption>Enabled</DefaultOption>
            <Help><![CDATA[Enable/Disable IPv4 PXE boot support. If disabled, IPv4 PXE boot support will not be available.]]></Help>
            <WorkIf><![CDATA[  0 != Network Stack  ]]></WorkIf>
          </Information>
        </Setting>
        <Setting name="Ipv4 HTTP Support" selectedOption="Disabled" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">Disabled</Option>
              <Option value="1">Enabled</Option>
            </AvailableOptions>
            <DefaultOption>Disabled</DefaultOption>
            <Help><![CDATA[Enable/Disable IPv4 HTTP boot support. If disabled, IPv4 HTTP boot support will not be available.]]></Help>
            <WorkIf><![CDATA[  0 != Network Stack  ]]></WorkIf>
          </Information>
        </Setting>
        <Setting name="Ipv6 PXE Support" selectedOption="Disabled" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">Disabled</Option>
              <Option value="1">Enabled</Option>
            </AvailableOptions>
            <DefaultOption>Disabled</DefaultOption>
            <Help><![CDATA[Enable/Disable IPv6 PXE boot support. If disabled, IPv6 PXE boot support will not be available.]]></Help>
            <WorkIf><![CDATA[  0 != Network Stack  ]]></WorkIf>
          </Information>
        </Setting>
        <Setting name="Ipv6 HTTP Support" selectedOption="Disabled" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">Disabled</Option>
              <Option value="1">Enabled</Option>
            </AvailableOptions>
            <DefaultOption>Disabled</DefaultOption>
            <Help><![CDATA[Enable/Disable IPv6 HTTP boot support. If disabled, IPv6 HTTP boot support will not be available.]]></Help>
            <WorkIf><![CDATA[  0 != Network Stack  ]]></WorkIf>
          </Information>
        </Setting>
        <Setting name="PXE boot wait time" numericValue="0" type="Numeric">
          <Information>
            <MaxValue>5</MaxValue>
            <MinValue>0</MinValue>
            <StepSize>1</StepSize>
            <DefaultValue>0</DefaultValue>
            <Help><![CDATA[Wait time in seconds to press ESC key to abort the PXE boot. Use either +/- or numeric keys to set the value.]]></Help>
            <WorkIf><![CDATA[  0 != Network Stack  ]]></WorkIf>
          </Information>
        </Setting>
        <Setting name="Media detect count" numericValue="1" type="Numeric">
          <Information>
            <MaxValue>50</MaxValue>
            <MinValue>1</MinValue>
            <StepSize>1</StepSize>
            <DefaultValue>1</DefaultValue>
            <Help><![CDATA[Number of times the presence of media will be checked. Use either +/- or numeric keys to set the value.]]></Help>
            <WorkIf><![CDATA[  0 != Network Stack  ]]></WorkIf>
          </Information>
        </Setting>
      </Menu>
      <Subtitle></Subtitle>
      <Subtitle></Subtitle>
      <Subtitle></Subtitle>
    </Menu>
    <Menu name="Super IO Configuration">
      <Information>
        <Help><![CDATA[System Super IO Chip Parameters.]]></Help>
      </Information>
      <Subtitle>Super IO Configuration</Subtitle>
      <Subtitle></Subtitle>
      <Text>Super IO Chip(AST2500)</Text>
      <Menu name="Serial Port 1 Configuration">
        <Information>
          <Help><![CDATA[Set Parameters of Serial Port 1 (COMA)]]></Help>
        </Information>
        <Subtitle>Serial Port 1 Configuration</Subtitle>
        <Subtitle></Subtitle>
        <Setting name="Serial Port 1" checkedStatus="Checked" type="CheckBox">
          <!--Checked/Unchecked-->
          <Information>
            <DefaultStatus>Checked</DefaultStatus>
            <Help><![CDATA[Enable or Disable Serial Port (COM)]]></Help>
          </Information>
        </Setting>
        <Text>Device Settings(IO=3F8h; IRQ=4;)</Text>
        <!--Valid if:   0 != Serial Port 1  -->
        <Subtitle></Subtitle>
        <Setting name="Change Settings" order="1" selectedOption="Auto" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">Auto</Option>
              <Option value="1">IO=3F8h; IRQ=4;</Option>
              <Option value="3">IO=2F8h; IRQ=4;</Option>
              <Option value="4">IO=3E8h; IRQ=4;</Option>
              <Option value="5">IO=2E8h; IRQ=4;</Option>
            </AvailableOptions>
            <DefaultOption>Auto</DefaultOption>
            <Help><![CDATA[Select an optimal settings for Super IO Device]]></Help>
            <WorkIf><![CDATA[  0 != Serial Port 1  ]]></WorkIf>
          </Information>
        </Setting>
      </Menu>
      <Menu name="Serial Port 2 Configuration">
        <Information>
          <Help><![CDATA[Set Parameters of Serial Port 2 (COMB)]]></Help>
        </Information>
        <Subtitle>Serial Port 2 Configuration</Subtitle>
        <Subtitle></Subtitle>
        <Setting name="Serial Port 2" checkedStatus="Checked" type="CheckBox">
          <!--Checked/Unchecked-->
          <Information>
            <DefaultStatus>Checked</DefaultStatus>
            <Help><![CDATA[Enable or Disable Serial Port (COM)]]></Help>
          </Information>
        </Setting>
        <Text>Device Settings(IO=2F8h; IRQ=3;)</Text>
        <!--Valid if:   0 != Serial Port 2  -->
        <Subtitle></Subtitle>
        <Setting name="Change Settings" order="2" selectedOption="Auto" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">Auto</Option>
              <Option value="1">IO=2F8h; IRQ=3;</Option>
              <Option value="2">IO=3F8h; IRQ=3;</Option>
              <Option value="4">IO=3E8h; IRQ=3;</Option>
              <Option value="5">IO=2E8h; IRQ=3;</Option>
            </AvailableOptions>
            <DefaultOption>Auto</DefaultOption>
            <Help><![CDATA[Select an optimal settings for Super IO Device]]></Help>
            <WorkIf><![CDATA[  0 != Serial Port 2  ]]></WorkIf>
          </Information>
        </Setting>
      </Menu>
    </Menu>
    <Menu name="Serial Port Console Redirection">
      <Information>
        <Help><![CDATA[Serial Port Console Redirection]]></Help>
      </Information>
      <Subtitle></Subtitle>
      <Subtitle>COM1</Subtitle>
      <Setting name="Console Redirection" order="1" checkedStatus="Unchecked" type="CheckBox">
        <!--Checked/Unchecked-->
        <Information>
          <DefaultStatus>Unchecked</DefaultStatus>
          <Help><![CDATA[Console Redirection Enable or Disable.]]></Help>
        </Information>
      </Setting>
      <Menu name="Console Redirection Settings" order="1">
        <Information>
          <Help><![CDATA[The settings specify how the host computer and the remote computer (which the user is using) will exchange data. Both computers should have the same or compatible settings.]]></Help>
          <WorkIf><![CDATA[ ( 0 != Console Redirection$1 ) ]]></WorkIf>
        </Information>
        <Subtitle>COM1</Subtitle>
        <Subtitle>Console Redirection Settings</Subtitle>
        <Subtitle></Subtitle>
        <Setting name="Terminal Type" order="1" selectedOption="VT100+" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">VT100</Option>
              <Option value="1">VT100+</Option>
              <Option value="2">VT-UTF8</Option>
              <Option value="3">ANSI</Option>
            </AvailableOptions>
            <DefaultOption>VT100+</DefaultOption>
            <Help><![CDATA[Emulation: ANSI: Extended ASCII char set. VT100: ASCII char set. VT100+: Extends VT100 to support color, function keys, etc. VT-UTF8: Uses UTF8 encoding to map Unicode chars onto 1 or more bytes.]]></Help>
          </Information>
        </Setting>
        <Setting name="Bits per second" order="1" selectedOption="115200" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="3">9600</Option>
              <Option value="4">19200</Option>
              <Option value="5">38400</Option>
              <Option value="6">57600</Option>
              <Option value="7">115200</Option>
            </AvailableOptions>
            <DefaultOption>115200</DefaultOption>
            <Help><![CDATA[Selects serial port transmission speed. The speed must be matched on the other side. Long or noisy lines may require lower speeds.]]></Help>
          </Information>
        </Setting>
        <Setting name="Data Bits" order="1" selectedOption="8" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="7">7</Option>
              <Option value="8">8</Option>
            </AvailableOptions>
            <DefaultOption>8</DefaultOption>
            <Help><![CDATA[Data Bits]]></Help>
          </Information>
        </Setting>
        <Setting name="Parity" order="1" selectedOption="None" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="1">None</Option>
              <Option value="2">Even</Option>
              <Option value="3">Odd</Option>
              <Option value="4">Mark</Option>
              <Option value="5">Space</Option>
            </AvailableOptions>
            <DefaultOption>None</DefaultOption>
            <Help><![CDATA[A parity bit can be sent with the data bits to detect some transmission errors. Even: parity bit is 0 if the num of 1's in the data bits is even. Odd: parity bit is 0 if num of 1's in the data bits is odd.  Mark: parity bit is always 1. Space: Parity bit is always 0. Mark and Space Parity do not allow for error detection. They can be used as an additional data bit.]]></Help>
          </Information>
        </Setting>
        <Setting name="Stop Bits" order="1" selectedOption="1" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="1">1</Option>
              <Option value="3">2</Option>
            </AvailableOptions>
            <DefaultOption>1</DefaultOption>
            <Help><![CDATA[Stop bits indicate the end of a serial data packet. (A start bit indicates the beginning). The standard setting is 1 stop bit. Communication with slow devices may require more than 1 stop bit.]]></Help>
          </Information>
        </Setting>
        <Setting name="Flow Control" order="1" selectedOption="None" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">None</Option>
              <Option value="1">Hardware RTS/CTS</Option>
            </AvailableOptions>
            <DefaultOption>None</DefaultOption>
            <Help><![CDATA[Flow control can prevent data loss from buffer overflow. When sending data, if the receiving buffers are full, a 'stop' signal can be sent to stop the data flow. Once the buffers are empty, a 'start' signal can be sent to re-start the flow. Hardware flow control uses two wires to send start/stop signals.]]></Help>
          </Information>
        </Setting>
        <Setting name="VT-UTF8 Combo Key Support" order="1" checkedStatus="Checked" type="CheckBox">
          <!--Checked/Unchecked-->
          <Information>
            <DefaultStatus>Checked</DefaultStatus>
            <Help><![CDATA[Enable VT-UTF8 Combination Key Support for ANSI/VT100 terminals]]></Help>
          </Information>
        </Setting>
        <Setting name="Recorder Mode" order="1" checkedStatus="Unchecked" type="CheckBox">
          <!--Checked/Unchecked-->
          <Information>
            <DefaultStatus>Unchecked</DefaultStatus>
            <Help><![CDATA[With this mode enabled only text will be sent. This is to capture Terminal data.]]></Help>
          </Information>
        </Setting>
        <Setting name="Resolution 100x31" order="1" checkedStatus="Checked" type="CheckBox">
          <!--Checked/Unchecked-->
          <Information>
            <DefaultStatus>Checked</DefaultStatus>
            <Help><![CDATA[Enables or disables extended terminal resolution]]></Help>
          </Information>
        </Setting>
        <Setting name="Legacy OS Redirection Resolution" order="1" selectedOption="80x24" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">80x24</Option>
              <Option value="1">80x25</Option>
            </AvailableOptions>
            <DefaultOption>80x24</DefaultOption>
            <Help><![CDATA[On Legacy OS, the Number of Rows and Columns supported redirection]]></Help>
          </Information>
        </Setting>
        <Setting name="Putty KeyPad" order="1" selectedOption="VT100" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="1">VT100</Option>
              <Option value="2">LINUX</Option>
              <Option value="4">XTERMR6</Option>
              <Option value="8">SCO</Option>
              <Option value="16">ESCN</Option>
              <Option value="32">VT400</Option>
            </AvailableOptions>
            <DefaultOption>VT100</DefaultOption>
            <Help><![CDATA[Select FunctionKey and KeyPad on Putty.]]></Help>
          </Information>
        </Setting>
        <Setting name="Redirection After BIOS POST" order="1" selectedOption="Always Enable" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">Always Enable</Option>
              <Option value="1">BootLoader</Option>
            </AvailableOptions>
            <DefaultOption>Always Enable</DefaultOption>
            <Help><![CDATA[When Bootloader is selected, then Legacy Console Redirection is disabled before booting to legacy OS. When Always Enable is selected, then Legacy Console Redirection is enabled for legacy OS. Default setting for this option is set to Always Enable.]]></Help>
          </Information>
        </Setting>
        <Subtitle></Subtitle>
      </Menu>
      <Subtitle></Subtitle>
      <Subtitle>SOL</Subtitle>
      <Setting name="Console Redirection" order="2" checkedStatus="Checked" type="CheckBox">
        <!--Checked/Unchecked-->
        <Information>
          <DefaultStatus>Checked</DefaultStatus>
          <Help><![CDATA[Console Redirection Enable or Disable.]]></Help>
        </Information>
      </Setting>
      <Menu name="Console Redirection Settings" order="2">
        <Information>
          <Help><![CDATA[The settings specify how the host computer and the remote computer (which the user is using) will exchange data. Both computers should have the same or compatible settings.]]></Help>
          <WorkIf><![CDATA[ ( 0 != Console Redirection$2 ) ]]></WorkIf>
        </Information>
        <Subtitle>SOL</Subtitle>
        <Subtitle>Console Redirection Settings</Subtitle>
        <Subtitle></Subtitle>
        <Setting name="Terminal Type" order="2" selectedOption="VT100+" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">VT100</Option>
              <Option value="1">VT100+</Option>
              <Option value="2">VT-UTF8</Option>
              <Option value="3">ANSI</Option>
            </AvailableOptions>
            <DefaultOption>VT100+</DefaultOption>
            <Help><![CDATA[Emulation: ANSI: Extended ASCII char set. VT100: ASCII char set. VT100+: Extends VT100 to support color, function keys, etc. VT-UTF8: Uses UTF8 encoding to map Unicode chars onto 1 or more bytes.]]></Help>
          </Information>
        </Setting>
        <Setting name="Bits per second" order="2" selectedOption="115200" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="3">9600</Option>
              <Option value="4">19200</Option>
              <Option value="5">38400</Option>
              <Option value="6">57600</Option>
              <Option value="7">115200</Option>
            </AvailableOptions>
            <DefaultOption>115200</DefaultOption>
            <Help><![CDATA[Selects serial port transmission speed. The speed must be matched on the other side. Long or noisy lines may require lower speeds.]]></Help>
          </Information>
        </Setting>
        <Setting name="Data Bits" order="2" selectedOption="8" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="7">7</Option>
              <Option value="8">8</Option>
            </AvailableOptions>
            <DefaultOption>8</DefaultOption>
            <Help><![CDATA[Data Bits]]></Help>
          </Information>
        </Setting>
        <Setting name="Parity" order="2" selectedOption="None" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="1">None</Option>
              <Option value="2">Even</Option>
              <Option value="3">Odd</Option>
              <Option value="4">Mark</Option>
              <Option value="5">Space</Option>
            </AvailableOptions>
            <DefaultOption>None</DefaultOption>
            <Help><![CDATA[A parity bit can be sent with the data bits to detect some transmission errors. Even: parity bit is 0 if the num of 1's in the data bits is even. Odd: parity bit is 0 if num of 1's in the data bits is odd.  Mark: parity bit is always 1. Space: Parity bit is always 0. Mark and Space Parity do not allow for error detection. They can be used as an additional data bit.]]></Help>
          </Information>
        </Setting>
        <Setting name="Stop Bits" order="2" selectedOption="1" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="1">1</Option>
              <Option value="3">2</Option>
            </AvailableOptions>
            <DefaultOption>1</DefaultOption>
            <Help><![CDATA[Stop bits indicate the end of a serial data packet. (A start bit indicates the beginning). The standard setting is 1 stop bit. Communication with slow devices may require more than 1 stop bit.]]></Help>
          </Information>
        </Setting>
        <Setting name="Flow Control" order="2" selectedOption="None" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">None</Option>
              <Option value="1">Hardware RTS/CTS</Option>
            </AvailableOptions>
            <DefaultOption>None</DefaultOption>
            <Help><![CDATA[Flow control can prevent data loss from buffer overflow. When sending data, if the receiving buffers are full, a 'stop' signal can be sent to stop the data flow. Once the buffers are empty, a 'start' signal can be sent to re-start the flow. Hardware flow control uses two wires to send start/stop signals.]]></Help>
          </Information>
        </Setting>
        <Setting name="VT-UTF8 Combo Key Support" order="2" checkedStatus="Checked" type="CheckBox">
          <!--Checked/Unchecked-->
          <Information>
            <DefaultStatus>Checked</DefaultStatus>
            <Help><![CDATA[Enable VT-UTF8 Combination Key Support for ANSI/VT100 terminals]]></Help>
          </Information>
        </Setting>
        <Setting name="Recorder Mode" order="2" checkedStatus="Unchecked" type="CheckBox">
          <!--Checked/Unchecked-->
          <Information>
            <DefaultStatus>Unchecked</DefaultStatus>
            <Help><![CDATA[With this mode enabled only text will be sent. This is to capture Terminal data.]]></Help>
          </Information>
        </Setting>
        <Setting name="Resolution 100x31" order="2" checkedStatus="Checked" type="CheckBox">
          <!--Checked/Unchecked-->
          <Information>
            <DefaultStatus>Checked</DefaultStatus>
            <Help><![CDATA[Enables or disables extended terminal resolution]]></Help>
          </Information>
        </Setting>
        <Setting name="Legacy OS Redirection Resolution" order="2" selectedOption="80x24" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">80x24</Option>
              <Option value="1">80x25</Option>
            </AvailableOptions>
            <DefaultOption>80x24</DefaultOption>
            <Help><![CDATA[On Legacy OS, the Number of Rows and Columns supported redirection]]></Help>
          </Information>
        </Setting>
        <Setting name="Putty KeyPad" order="2" selectedOption="VT100" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="1">VT100</Option>
              <Option value="2">LINUX</Option>
              <Option value="4">XTERMR6</Option>
              <Option value="8">SCO</Option>
              <Option value="16">ESCN</Option>
              <Option value="32">VT400</Option>
            </AvailableOptions>
            <DefaultOption>VT100</DefaultOption>
            <Help><![CDATA[Select FunctionKey and KeyPad on Putty.]]></Help>
          </Information>
        </Setting>
        <Setting name="Redirection After BIOS POST" order="2" selectedOption="Always Enable" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">Always Enable</Option>
              <Option value="1">BootLoader</Option>
            </AvailableOptions>
            <DefaultOption>Always Enable</DefaultOption>
            <Help><![CDATA[When Bootloader is selected, then Legacy Console Redirection is disabled before booting to legacy OS. When Always Enable is selected, then Legacy Console Redirection is enabled for legacy OS. Default setting for this option is set to Always Enable.]]></Help>
          </Information>
        </Setting>
        <Subtitle></Subtitle>
      </Menu>
      <Subtitle></Subtitle>
      <Subtitle>Legacy Console Redirection</Subtitle>
      <Setting name="Redirection COM Port" selectedOption="COM1" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">COM1</Option>
            <Option value="1">SOL</Option>
          </AvailableOptions>
          <DefaultOption>COM1</DefaultOption>
          <Help><![CDATA[Select a COM port to display redirection of Legacy OS and Legacy OPROM Messages]]></Help>
        </Information>
      </Setting>
      <Subtitle></Subtitle>
      <Subtitle>Serial Port for Out-of-Band Management/</Subtitle>
      <Subtitle>Windows Emergency Management Services (EMS)</Subtitle>
      <Setting name="Console Redirection" order="3" checkedStatus="Unchecked" type="CheckBox">
        <!--Checked/Unchecked-->
        <Information>
          <DefaultStatus>Unchecked</DefaultStatus>
          <Help><![CDATA[Console Redirection Enable or Disable.]]></Help>
        </Information>
      </Setting>
      <Menu name="Console Redirection Settings" order="3">
        <Information>
          <Help><![CDATA[The settings specify how the host computer and the remote computer (which the user is using) will exchange data. Both computers should have the same or compatible settings.]]></Help>
          <WorkIf><![CDATA[  0 != Console Redirection$3  ]]></WorkIf>
        </Information>
        <Setting name="Out-of-Band Mgmt Port" selectedOption="COM1" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">COM1</Option>
              <Option value="1">SOL</Option>
            </AvailableOptions>
            <DefaultOption>COM1</DefaultOption>
            <Help><![CDATA[Microsoft Windows Emergency Management Services (EMS) allows for remote management of a Windows Server OS through a serial port.]]></Help>
          </Information>
        </Setting>
        <Setting name="Terminal Type" order="3" selectedOption="VT-UTF8" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">VT100</Option>
              <Option value="1">VT100+</Option>
              <Option value="2">VT-UTF8</Option>
              <Option value="3">ANSI</Option>
            </AvailableOptions>
            <DefaultOption>VT-UTF8</DefaultOption>
            <Help><![CDATA[VT-UTF8 is the preferred terminal type for out-of-band management. The next best choice is VT100+ and then VT100. See above, in Console Redirection Settings page, for more Help with Terminal Type/Emulation.]]></Help>
            <WorkIf><![CDATA[  0 != Console Redirection$3  ]]></WorkIf>
          </Information>
        </Setting>
        <Setting name="Bits per second" order="3" selectedOption="115200" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="3">9600</Option>
              <Option value="4">19200</Option>
              <Option value="6">57600</Option>
              <Option value="7">115200</Option>
            </AvailableOptions>
            <DefaultOption>115200</DefaultOption>
            <Help><![CDATA[Selects serial port transmission speed. The speed must be matched on the other side. Long or noisy lines may require lower speeds.]]></Help>
            <WorkIf><![CDATA[  0 != Console Redirection$3  ]]></WorkIf>
          </Information>
        </Setting>
        <Setting name="Flow Control" order="3" selectedOption="None" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">None</Option>
              <Option value="1">Hardware RTS/CTS</Option>
              <Option value="2">Software Xon/Xoff</Option>
            </AvailableOptions>
            <DefaultOption>None</DefaultOption>
            <Help><![CDATA[Flow control can prevent data loss from buffer overflow. When sending data, if the receiving buffers are full, a 'stop' signal can be sent to stop the data flow. Once the buffers are empty, a 'start' signal can be sent to re-start the flow. Hardware flow control uses two wires to send start/stop signals.]]></Help>
            <WorkIf><![CDATA[  0 != Console Redirection$3  ]]></WorkIf>
          </Information>
        </Setting>
        <Text>Data Bits(8)</Text>
        <!--Valid if:   0 != Console Redirection  -->
        <Text>Parity(None)</Text>
        <!--Valid if:   0 != Console Redirection  -->
        <Text>Stop Bits(1)</Text>
        <!--Valid if:   0 != Console Redirection  -->
      </Menu>
    </Menu>
    <Menu name="ACPI Settings">
      <Information>
        <Help><![CDATA[System ACPI Parameters.]]></Help>
      </Information>
      <Subtitle>ACPI Settings</Subtitle>
      <Subtitle></Subtitle>
      <Setting name="PCI AER Support" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[$SMCUNHIDE$Enable/Disable ACPI OS to natively manage PCI Advanced Error Reporting.]]></Help>
        </Information>
      </Setting>
      <Setting name="Memory Corrected Error Enabling" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[$SMCUNHIDE$Enable/Disable Memory Corrected Error]]></Help>
        </Information>
      </Setting>
      <Setting name="Headless Support" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[Enable/Disable ACPI OS indicates the system cannot detect the monitor or keyboard / mouse devices.]]></Help>
        </Information>
      </Setting>
      <Setting name="WHEA Support" selectedOption="Enabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Enabled</DefaultOption>
          <Help><![CDATA[Enable/Disable WHEA support]]></Help>
        </Information>
      </Setting>
      <Setting name="High Precision Event Timer" selectedOption="Enabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Enabled</DefaultOption>
          <Help><![CDATA[Enable or Disable the High Precision Event Timer.]]></Help>
        </Information>
      </Setting>
    </Menu>
    <Menu name="Trusted Computing">
      <Information>
        <Help><![CDATA[Trusted Computing Settings]]></Help>
      </Information>
      <Subtitle>Configuration</Subtitle>
      <Setting name="  Security Device Support" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Enables or Disables BIOS support for security device. O.S. will not show Security Device. TCG EFI protocol and INT1A interface will not be available.]]></Help>
        </Information>
      </Setting>
      <Text>  NO Security Device Found()</Text>
    </Menu>
    <Menu name="Network Configuration">
      <Information>
        <Help><![CDATA[Network Configuration Settings]]></Help>
        <WorkIf><![CDATA[  0 != Onboard LAN Option ROM Type  ]]></WorkIf>
      </Information>
    </Menu>
    <Menu name="Driver Health">
      <Information>
        <Help><![CDATA[Provides Health Status for the Drivers/Controllers]]></Help>
      </Information>
      <Menu name="">
        <Information>
          <Help><![CDATA[Provides Health Status for the Drivers/Controllers]]></Help>
        </Information>
      </Menu>
    </Menu>
  </Menu>
  <Menu name="Event Logs">
    <Information />
    <Menu name="Change SMBIOS Event Log Settings">
      <Information>
        <Help><![CDATA[Press <Enter> to change the SMBIOS Event Log configuration.]]></Help>
      </Information>
      <Subtitle>Enabling/Disabling Options</Subtitle>
      <Setting name="SMBIOS Event Log" selectedOption="Enabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Enabled</DefaultOption>
          <Help><![CDATA[Change this to enable or disable all features of SMBIOS Event Logging during boot.]]></Help>
        </Information>
      </Setting>
      <Subtitle></Subtitle>
      <Subtitle>Erasing Settings</Subtitle>
      <Setting name="Erase Event Log" selectedOption="No" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">No</Option>
            <Option value="1">Yes, Next reset</Option>
            <Option value="2">Yes, Every reset</Option>
          </AvailableOptions>
          <DefaultOption>No</DefaultOption>
          <Help><![CDATA[Choose options for erasing SMBIOS Event Log.  Erasing is done prior to any logging activation during reset.]]></Help>
          <WorkIf><![CDATA[ ( 0 != SMBIOS Event Log ) ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="When Log is Full" selectedOption="Do Nothing" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Do Nothing</Option>
            <Option value="1">Erase Immediately</Option>
          </AvailableOptions>
          <DefaultOption>Do Nothing</DefaultOption>
          <Help><![CDATA[Choose options for reactions to a full SMBIOS Event Log.]]></Help>
          <WorkIf><![CDATA[ ( 0 != SMBIOS Event Log ) ]]></WorkIf>
        </Information>
      </Setting>
      <Subtitle></Subtitle>
      <Subtitle>SMBIOS Event Log Standard Settings</Subtitle>
      <Setting name="Log System Boot Event" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Enabled</Option>
            <Option value="0">Disabled</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[Choose option to enable/disable logging of System boot event]]></Help>
          <WorkIf><![CDATA[ ( 0 != SMBIOS Event Log ) ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="MECI" numericValue="1" type="Numeric">
        <Information>
          <MaxValue>255</MaxValue>
          <MinValue>1</MinValue>
          <StepSize>1</StepSize>
          <DefaultValue>1</DefaultValue>
          <Help><![CDATA[Mutiple Event Count Increment:  The number of occurrences of a duplicate event that must pass before the multiple-event counter of log entry is updated.The value ranges from 1 to 255.]]></Help>
          <WorkIf><![CDATA[ ( 0 != SMBIOS Event Log ) ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="METW" numericValue="60" type="Numeric">
        <Information>
          <MaxValue>99</MaxValue>
          <MinValue>0</MinValue>
          <StepSize>1</StepSize>
          <DefaultValue>60</DefaultValue>
          <Help><![CDATA[Multiple Event Time Window:  The number of minutes which must pass between duplicate log entries which utilize a multiple-event counter. The value ranges from 0 to 99 minutes.]]></Help>
          <WorkIf><![CDATA[ ( 0 != SMBIOS Event Log ) ]]></WorkIf>
        </Information>
      </Setting>
      <Subtitle></Subtitle>
      <Subtitle>NOTE: All values changed here do not take effect</Subtitle>
      <Subtitle>      until computer is restarted.</Subtitle>
    </Menu>
    <Menu name="View SMBIOS Event Log">
      <Information>
        <Help><![CDATA[Press <Enter> to view the SMBIOS Event Log records.]]></Help>
      </Information>
      <Subtitle>DATE      TIME       ERROR CODE     SEVERITY</Subtitle>
      <Subtitle></Subtitle>
    </Menu>
    <Subtitle></Subtitle>
  </Menu>
  <Menu name="Security">
    <Information />
    <Text>Administrator Password(Not Installed)</Text>
    <!--Valid if:   0 == Administrator Password  -->
    <Text>Administrator Password(Installed)</Text>
    <!--Valid if:   0 != Administrator Password  -->
    <Text>User Password(Not Installed)</Text>
    <!--Valid if:   0 == User Password  -->
    <Text>User Password(Installed)</Text>
    <!--Valid if:   0 != User Password  -->
    <Subtitle></Subtitle>
    <Subtitle>Password Description</Subtitle>
    <Subtitle></Subtitle>
    <Subtitle>If the Administrator's / User's password is set, </Subtitle>
    <Subtitle>then this only limits access to Setup and is</Subtitle>
    <Subtitle>asked for when entering Setup.</Subtitle>
    <Subtitle>Please set Administrator's password first in order </Subtitle>
    <Subtitle>to set User's password, if clear Administrator's </Subtitle>
    <Subtitle>password, the User's password will be cleared as well.</Subtitle>
    <Subtitle></Subtitle>
    <Subtitle>The password length must be in the following range:</Subtitle>
    <Subtitle>in the following range:</Subtitle>
    <Text>Minimum length(3)</Text>
    <Text>Maximum length(20)</Text>
    <Subtitle></Subtitle>
    <Setting name="Administrator Password" type="Password">
      <Information>
        <Help>Set Administrator Password</Help>
        <MinSize>3</MinSize>
        <MaxSize>20</MaxSize>
        <HasPassword>False</HasPassword>
      </Information>
      <NewPassword><![CDATA[]]></NewPassword>
      <ConfirmNewPassword><![CDATA[]]></ConfirmNewPassword>
    </Setting>
    <Setting name="User Password" type="Password">
      <Information>
        <Help>Set User Password</Help>
        <WorkIf><![CDATA[  0 != Administrator Password  ]]></WorkIf>
        <MinSize>3</MinSize>
        <MaxSize>20</MaxSize>
        <HasPassword>False</HasPassword>
      </Information>
      <NewPassword><![CDATA[]]></NewPassword>
      <ConfirmNewPassword><![CDATA[]]></ConfirmNewPassword>
    </Setting>
    <Setting name="Password Check" selectedOption="Setup" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Setup</Option>
          <Option value="1">Always</Option>
        </AvailableOptions>
        <DefaultOption>Setup</DefaultOption>
        <Help><![CDATA[Setup: Check password while invoking setup. Always: Check password while invoking setup as well as on each boot.]]></Help>
      </Information>
    </Setting>
    <Setting name="Hard Drive Security Frozen" selectedOption="Enable" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="1">Enable</Option>
          <Option value="0">Disable</Option>
        </AvailableOptions>
        <DefaultOption>Enable</DefaultOption>
        <Help><![CDATA[Enable/Disable BIOS Security Frozen Command to SATA and NVME Devices]]></Help>
      </Information>
    </Setting>
    <Subtitle>HDD Security Configuration:</Subtitle>
    <Subtitle></Subtitle>
  </Menu>
  <Menu name="Boot">
    <Information />
    <Subtitle>Boot Configuration</Subtitle>
    <Subtitle></Subtitle>
    <Setting name="Driver Option #%d" selectedOption="" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0"></Option>
          <Option value="1"></Option>
        </AvailableOptions>
        <DefaultOption></DefaultOption>
        <Help><![CDATA[Sets the system driver order]]></Help>
      </Information>
    </Setting>
    <Setting name="Boot mode select" selectedOption="Dual" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Legacy</Option>
          <Option value="1">UEFI</Option>
          <Option value="2">Dual</Option>
        </AvailableOptions>
        <DefaultOption>Dual</DefaultOption>
        <Help><![CDATA[Select boot mode Legacy/UEFI]]></Help>
      </Information>
    </Setting>
    <Setting name="Legacy to EFI support" selectedOption="Disabled" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Disabled</Option>
          <Option value="1">Enabled</Option>
        </AvailableOptions>
        <DefaultOption>Disabled</DefaultOption>
        <Help><![CDATA[Enabled: System is able to boot to EFI OS after boot failed from Legacy boot order.]]></Help>
      </Information>
    </Setting>
    <Subtitle></Subtitle>
    <Subtitle>FIXED BOOT ORDER Priorities</Subtitle>
    <Setting name="UEFI Boot Option #1" selectedOption="UEFI Hard Disk:metal-ubuntu" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="1">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="2">UEFI CD/DVD</Option>
          <Option value="3">UEFI USB Hard Disk</Option>
          <Option value="4">UEFI USB CD/DVD</Option>
          <Option value="5">UEFI USB Key</Option>
          <Option value="6">UEFI USB Floppy</Option>
          <Option value="7">UEFI USB Lan</Option>
          <Option value="8">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="9">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI Hard Disk:metal-ubuntu</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="UEFI Boot Option #2" selectedOption="UEFI AP:UEFI: Built-in EFI Shell" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="1">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="2">UEFI CD/DVD</Option>
          <Option value="3">UEFI USB Hard Disk</Option>
          <Option value="4">UEFI USB CD/DVD</Option>
          <Option value="5">UEFI USB Key</Option>
          <Option value="6">UEFI USB Floppy</Option>
          <Option value="7">UEFI USB Lan</Option>
          <Option value="8">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="9">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI AP:UEFI: Built-in EFI Shell</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="UEFI Boot Option #3" selectedOption="UEFI CD/DVD" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="1">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="2">UEFI CD/DVD</Option>
          <Option value="3">UEFI USB Hard Disk</Option>
          <Option value="4">UEFI USB CD/DVD</Option>
          <Option value="5">UEFI USB Key</Option>
          <Option value="6">UEFI USB Floppy</Option>
          <Option value="7">UEFI USB Lan</Option>
          <Option value="8">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="9">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI CD/DVD</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="UEFI Boot Option #4" selectedOption="UEFI USB Hard Disk" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="1">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="2">UEFI CD/DVD</Option>
          <Option value="3">UEFI USB Hard Disk</Option>
          <Option value="4">UEFI USB CD/DVD</Option>
          <Option value="5">UEFI USB Key</Option>
          <Option value="6">UEFI USB Floppy</Option>
          <Option value="7">UEFI USB Lan</Option>
          <Option value="8">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="9">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI USB Hard Disk</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="UEFI Boot Option #5" selectedOption="UEFI USB CD/DVD" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="1">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="2">UEFI CD/DVD</Option>
          <Option value="3">UEFI USB Hard Disk</Option>
          <Option value="4">UEFI USB CD/DVD</Option>
          <Option value="5">UEFI USB Key</Option>
          <Option value="6">UEFI USB Floppy</Option>
          <Option value="7">UEFI USB Lan</Option>
          <Option value="8">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="9">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI USB CD/DVD</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="UEFI Boot Option #6" selectedOption="UEFI USB Key" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="1">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="2">UEFI CD/DVD</Option>
          <Option value="3">UEFI USB Hard Disk</Option>
          <Option value="4">UEFI USB CD/DVD</Option>
          <Option value="5">UEFI USB Key</Option>
          <Option value="6">UEFI USB Floppy</Option>
          <Option value="7">UEFI USB Lan</Option>
          <Option value="8">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="9">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI USB Key</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="UEFI Boot Option #7" selectedOption="UEFI USB Floppy" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="1">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="2">UEFI CD/DVD</Option>
          <Option value="3">UEFI USB Hard Disk</Option>
          <Option value="4">UEFI USB CD/DVD</Option>
          <Option value="5">UEFI USB Key</Option>
          <Option value="6">UEFI USB Floppy</Option>
          <Option value="7">UEFI USB Lan</Option>
          <Option value="8">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="9">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI USB Floppy</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="UEFI Boot Option #8" selectedOption="UEFI USB Lan" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="1">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="2">UEFI CD/DVD</Option>
          <Option value="3">UEFI USB Hard Disk</Option>
          <Option value="4">UEFI USB CD/DVD</Option>
          <Option value="5">UEFI USB Key</Option>
          <Option value="6">UEFI USB Floppy</Option>
          <Option value="7">UEFI USB Lan</Option>
          <Option value="8">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="9">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI USB Lan</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="UEFI Boot Option #9" selectedOption="UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="1">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="2">UEFI CD/DVD</Option>
          <Option value="3">UEFI USB Hard Disk</Option>
          <Option value="4">UEFI USB CD/DVD</Option>
          <Option value="5">UEFI USB Key</Option>
          <Option value="6">UEFI USB Floppy</Option>
          <Option value="7">UEFI USB Lan</Option>
          <Option value="8">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="9">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Legacy Boot Option #1" selectedOption="Hard Disk" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>Hard Disk</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 1 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Legacy Boot Option #2" selectedOption="CD/DVD" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>CD/DVD</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 1 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Legacy Boot Option #3" selectedOption="USB Hard Disk" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>USB Hard Disk</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 1 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Legacy Boot Option #4" selectedOption="USB CD/DVD" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>USB CD/DVD</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 1 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Legacy Boot Option #5" selectedOption="USB Key" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>USB Key</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 1 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Legacy Boot Option #6" selectedOption="USB Floppy" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>USB Floppy</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 1 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Legacy Boot Option #7" selectedOption="USB Lan" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>USB Lan</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 1 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Legacy Boot Option #8" selectedOption="Network" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>Network</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 1 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Dual Boot Option #1" selectedOption="UEFI Hard Disk:metal-ubuntu" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="16">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>Hard Disk</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Dual Boot Option #2" selectedOption="CD/DVD" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="16">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>CD/DVD</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Dual Boot Option #3" selectedOption="USB Hard Disk" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="16">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>USB Hard Disk</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Dual Boot Option #4" selectedOption="USB CD/DVD" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="16">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>USB CD/DVD</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Dual Boot Option #5" selectedOption="USB Key" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="16">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>USB Key</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Dual Boot Option #6" selectedOption="USB Floppy" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="16">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>USB Floppy</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Dual Boot Option #7" selectedOption="USB Lan" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="16">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>USB Lan</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Dual Boot Option #8" selectedOption="Network" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="16">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>Network</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Dual Boot Option #9" selectedOption="UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="16">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI Hard Disk:metal-ubuntu</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Dual Boot Option #10" selectedOption="UEFI CD/DVD" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="16">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI CD/DVD</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Dual Boot Option #11" selectedOption="UEFI USB Hard Disk" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="16">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI USB Hard Disk</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Dual Boot Option #12" selectedOption="UEFI USB CD/DVD" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="16">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI USB CD/DVD</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Dual Boot Option #13" selectedOption="UEFI USB Key" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="16">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI USB Key</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Dual Boot Option #14" selectedOption="UEFI USB Floppy" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="16">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI USB Floppy</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Dual Boot Option #15" selectedOption="UEFI USB Lan" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="16">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI USB Lan</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Dual Boot Option #16" selectedOption="Hard Disk" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="16">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Dual Boot Option #17" selectedOption="UEFI AP:UEFI: Built-in EFI Shell" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk</Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk:metal-ubuntu</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection</Option>
          <Option value="16">UEFI AP:UEFI: Built-in EFI Shell</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI AP:UEFI: Built-in EFI Shell</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Subtitle></Subtitle>
    <Subtitle></Subtitle>
    <Menu name="UEFI Hard Disk Drive BBS Priorities">
      <Information>
        <Help><![CDATA[Specifies the Boot Device Priority sequence from available UEFI Hard Disk Drives.]]></Help>
        <WorkIf><![CDATA[ ( Boot mode select is not in 0  ) ]]></WorkIf>
      </Information>
      <Setting name="Boot Option #1" order="1" selectedOption="metal-ubuntu(SATA,Port:4)" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">metal-ubuntu(SATA,Port:4)</Option>
            <Option value="1">UEFI OS(SATA,Port:4)</Option>
            <Option value="18">Disabled</Option>
          </AvailableOptions>
          <DefaultOption>metal-ubuntu(SATA,Port:4)</DefaultOption>
          <Help><![CDATA[Sets the system boot order]]></Help>
        </Information>
      </Setting>
      <Setting name="Boot Option #2" order="1" selectedOption="UEFI OS(SATA,Port:4)" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">metal-ubuntu(SATA,Port:4)</Option>
            <Option value="1">UEFI OS(SATA,Port:4)</Option>
            <Option value="18">Disabled</Option>
          </AvailableOptions>
          <DefaultOption>UEFI OS(SATA,Port:4)</DefaultOption>
          <Help><![CDATA[Sets the system boot order]]></Help>
        </Information>
      </Setting>
    </Menu>
    <Menu name="UEFI Application Boot Priorities">
      <Information>
        <Help><![CDATA[Specifies the Boot Device Priority sequence from available UEFI Application.]]></Help>
        <WorkIf><![CDATA[ ( Boot mode select is not in 0  ) ]]></WorkIf>
      </Information>
      <Setting name="Boot Option #1" order="2" selectedOption="UEFI: Built-in EFI Shell" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">UEFI: Built-in EFI Shell</Option>
            <Option value="18">Disabled</Option>
          </AvailableOptions>
          <DefaultOption>UEFI: Built-in EFI Shell</DefaultOption>
          <Help><![CDATA[Sets the system boot order]]></Help>
        </Information>
      </Setting>
    </Menu>
    <Menu name="UEFI NETWORK Drive BBS Priorities">
      <Information>
        <Help><![CDATA[Specifies the Boot Device Priority sequence from available UEFI NETWORK Drives.]]></Help>
        <WorkIf><![CDATA[ ( Boot mode select is not in 0  ) ]]></WorkIf>
      </Information>
      <Setting name="Boot Option #1" order="3" selectedOption="UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efa)" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efa)</Option>
            <Option value="1">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efb)</Option>
            <Option value="2">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efc)</Option>
            <Option value="3">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efd)</Option>
            <Option value="4">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GBASE-T(MAC,Address:ac1f6b7d840a)</Option>
            <Option value="5">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GBASE-T(MAC,Address:ac1f6b7d840b)</Option>
            <Option value="6">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GbE SFP+(MAC,Address:ac1f6b7d840c)</Option>
            <Option value="7">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GbE SFP+(MAC,Address:ac1f6b7d840d)</Option>
            <Option value="18">Disabled</Option>
          </AvailableOptions>
          <DefaultOption>UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efa)</DefaultOption>
          <Help><![CDATA[Sets the system boot order]]></Help>
        </Information>
      </Setting>
      <Setting name="Boot Option #2" order="2" selectedOption="UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efb)" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efa)</Option>
            <Option value="1">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efb)</Option>
            <Option value="2">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efc)</Option>
            <Option value="3">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efd)</Option>
            <Option value="4">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GBASE-T(MAC,Address:ac1f6b7d840a)</Option>
            <Option value="5">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GBASE-T(MAC,Address:ac1f6b7d840b)</Option>
            <Option value="6">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GbE SFP+(MAC,Address:ac1f6b7d840c)</Option>
            <Option value="7">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GbE SFP+(MAC,Address:ac1f6b7d840d)</Option>
            <Option value="18">Disabled</Option>
          </AvailableOptions>
          <DefaultOption>UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efb)</DefaultOption>
          <Help><![CDATA[Sets the system boot order]]></Help>
        </Information>
      </Setting>
      <Setting name="Boot Option #3" selectedOption="UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efc)" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efa)</Option>
            <Option value="1">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efb)</Option>
            <Option value="2">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efc)</Option>
            <Option value="3">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efd)</Option>
            <Option value="4">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GBASE-T(MAC,Address:ac1f6b7d840a)</Option>
            <Option value="5">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GBASE-T(MAC,Address:ac1f6b7d840b)</Option>
            <Option value="6">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GbE SFP+(MAC,Address:ac1f6b7d840c)</Option>
            <Option value="7">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GbE SFP+(MAC,Address:ac1f6b7d840d)</Option>
            <Option value="18">Disabled</Option>
          </AvailableOptions>
          <DefaultOption>UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efc)</DefaultOption>
          <Help><![CDATA[Sets the system boot order]]></Help>
        </Information>
      </Setting>
      <Setting name="Boot Option #4" selectedOption="UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efd)" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efa)</Option>
            <Option value="1">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efb)</Option>
            <Option value="2">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efc)</Option>
            <Option value="3">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efd)</Option>
            <Option value="4">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GBASE-T(MAC,Address:ac1f6b7d840a)</Option>
            <Option value="5">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GBASE-T(MAC,Address:ac1f6b7d840b)</Option>
            <Option value="6">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GbE SFP+(MAC,Address:ac1f6b7d840c)</Option>
            <Option value="7">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GbE SFP+(MAC,Address:ac1f6b7d840d)</Option>
            <Option value="18">Disabled</Option>
          </AvailableOptions>
          <DefaultOption>UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efd)</DefaultOption>
          <Help><![CDATA[Sets the system boot order]]></Help>
        </Information>
      </Setting>
      <Setting name="Boot Option #5" selectedOption="UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GBASE-T(MAC,Address:ac1f6b7d840a)" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efa)</Option>
            <Option value="1">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efb)</Option>
            <Option value="2">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efc)</Option>
            <Option value="3">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efd)</Option>
            <Option value="4">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GBASE-T(MAC,Address:ac1f6b7d840a)</Option>
            <Option value="5">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GBASE-T(MAC,Address:ac1f6b7d840b)</Option>
            <Option value="6">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GbE SFP+(MAC,Address:ac1f6b7d840c)</Option>
            <Option value="7">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GbE SFP+(MAC,Address:ac1f6b7d840d)</Option>
            <Option value="18">Disabled</Option>
          </AvailableOptions>
          <DefaultOption>UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GBASE-T(MAC,Address:ac1f6b7d840a)</DefaultOption>
          <Help><![CDATA[Sets the system boot order]]></Help>
        </Information>
      </Setting>
      <Setting name="Boot Option #6" selectedOption="UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GBASE-T(MAC,Address:ac1f6b7d840b)" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efa)</Option>
            <Option value="1">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efb)</Option>
            <Option value="2">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efc)</Option>
            <Option value="3">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efd)</Option>
            <Option value="4">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GBASE-T(MAC,Address:ac1f6b7d840a)</Option>
            <Option value="5">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GBASE-T(MAC,Address:ac1f6b7d840b)</Option>
            <Option value="6">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GbE SFP+(MAC,Address:ac1f6b7d840c)</Option>
            <Option value="7">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GbE SFP+(MAC,Address:ac1f6b7d840d)</Option>
            <Option value="18">Disabled</Option>
          </AvailableOptions>
          <DefaultOption>UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GBASE-T(MAC,Address:ac1f6b7d840b)</DefaultOption>
          <Help><![CDATA[Sets the system boot order]]></Help>
        </Information>
      </Setting>
      <Setting name="Boot Option #7" selectedOption="UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GbE SFP+(MAC,Address:ac1f6b7d840c)" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efa)</Option>
            <Option value="1">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efb)</Option>
            <Option value="2">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efc)</Option>
            <Option value="3">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efd)</Option>
            <Option value="4">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GBASE-T(MAC,Address:ac1f6b7d840a)</Option>
            <Option value="5">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GBASE-T(MAC,Address:ac1f6b7d840b)</Option>
            <Option value="6">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GbE SFP+(MAC,Address:ac1f6b7d840c)</Option>
            <Option value="7">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GbE SFP+(MAC,Address:ac1f6b7d840d)</Option>
            <Option value="18">Disabled</Option>
          </AvailableOptions>
          <DefaultOption>UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GbE SFP+(MAC,Address:ac1f6b7d840c)</DefaultOption>
          <Help><![CDATA[Sets the system boot order]]></Help>
        </Information>
      </Setting>
      <Setting name="Boot Option #8" selectedOption="UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GbE SFP+(MAC,Address:ac1f6b7d840d)" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efa)</Option>
            <Option value="1">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efb)</Option>
            <Option value="2">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efc)</Option>
            <Option value="3">UEFI: PXE IPv4 Intel(R) I350 Gigabit Network Connection(MAC,Address:ac1f6b7d7efd)</Option>
            <Option value="4">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GBASE-T(MAC,Address:ac1f6b7d840a)</Option>
            <Option value="5">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GBASE-T(MAC,Address:ac1f6b7d840b)</Option>
            <Option value="6">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GbE SFP+(MAC,Address:ac1f6b7d840c)</Option>
            <Option value="7">UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GbE SFP+(MAC,Address:ac1f6b7d840d)</Option>
            <Option value="18">Disabled</Option>
          </AvailableOptions>
          <DefaultOption>UEFI: PXE IPv4 Intel(R) Ethernet Connection X722 for 10GbE SFP+(MAC,Address:ac1f6b7d840d)</DefaultOption>
          <Help><![CDATA[Sets the system boot order]]></Help>
        </Information>
      </Setting>
    </Menu>
  </Menu>
</BiosCfg>
`
	testBigTwinBiosCfg = `<?xml version="1.0" encoding="ISO-8859-1" standalone="yes"?>
<BiosCfg>
  <Menu name="Main">
    <Information />
    <Subtitle></Subtitle>
    <Subtitle></Subtitle>
    <Subtitle>Supermicro X11DPT-B</Subtitle>
    <Text>BIOS Version(3.0a)</Text>
    <Text>Build Date(02/20/2019)</Text>
    <Text>CPLD Version(03.B0.09)</Text>
    <Subtitle></Subtitle>
    <Subtitle>Memory Information</Subtitle>
    <Text>Total Memory(98304 MB)</Text>
  </Menu>
  <Menu name="Advanced">
    <Information />
    <Menu name="Boot Feature">
      <Information>
        <Help><![CDATA[Boot Feature Configuration Page]]></Help>
      </Information>
      <Subtitle></Subtitle>
      <Setting name="Quiet Boot" checkedStatus="Unchecked" type="CheckBox">
        <!--Checked/Unchecked-->
        <Information>
          <DefaultStatus>Checked</DefaultStatus>
          <Help><![CDATA[Enables or disables Quiet Boot option]]></Help>
        </Information>
      </Setting>
      <Subtitle></Subtitle>
      <Setting name="Option ROM Messages" selectedOption="Force BIOS" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Force BIOS</Option>
            <Option value="0">Keep Current</Option>
          </AvailableOptions>
          <DefaultOption>Force BIOS</DefaultOption>
          <Help><![CDATA[Set display mode for Option ROM]]></Help>
        </Information>
      </Setting>
      <Setting name="Bootup NumLock State" selectedOption="On" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">On</Option>
            <Option value="0">Off</Option>
          </AvailableOptions>
          <DefaultOption>On</DefaultOption>
          <Help><![CDATA[Select the keyboard NumLock state]]></Help>
        </Information>
      </Setting>
      <Setting name="Wait For &quot;F1&quot; If Error" selectedOption="Enabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Enabled</DefaultOption>
          <Help><![CDATA[Enable- BIOS will wait for user to press "F1" if some error happens. Disable- BIOS will continue to POST, user interaction not required]]></Help>
        </Information>
      </Setting>
      <Setting name="INT19 Trap Response" selectedOption="Immediate" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Immediate</Option>
            <Option value="0">Postponed</Option>
          </AvailableOptions>
          <DefaultOption>Immediate</DefaultOption>
          <Help><![CDATA[BIOS reaction on INT19 trapping by Option ROM: IMMEDIATE - execute the trap right away; POSTPONED - execute the trap during legacy boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="Re-try Boot" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Legacy Boot</Option>
            <Option value="2">EFI Boot</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[Decide how to retry boot devices which fail to boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="Install Windows 7 USB Support" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[When enabled, install Windows 7 USB Keyboard/Mouse can be used. After install Windows 7 & XHCI driver please set to "Disabled"]]></Help>
        </Information>
      </Setting>
      <Setting name="Port 61h Bit-4 Emulation" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[Emulation of Port 61h bit-4 toggling in SMM]]></Help>
        </Information>
      </Setting>
      <Subtitle></Subtitle>
      <Subtitle>Power Configuration</Subtitle>
      <Setting name="Watch Dog Function" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[Enable or disable to turn on 5-minute watch dog timer. Upon timeout, JWD1 jumper determines system behavior.]]></Help>
        </Information>
      </Setting>
      <Setting name="Restore on AC Power Loss" selectedOption="Last State" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Stay Off</Option>
            <Option value="1">Power On</Option>
            <Option value="2">Last State</Option>
          </AvailableOptions>
          <DefaultOption>Last State</DefaultOption>
          <Help><![CDATA[Stay Off: System always remains off.
Power On: System always turns on.
Last State: System returns to previous state before AC lost.]]></Help>
        </Information>
      </Setting>
      <Setting name="Power Button Function" selectedOption="Instant Off" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Instant Off</Option>
            <Option value="0">4 Seconds Override</Option>
          </AvailableOptions>
          <DefaultOption>Instant Off</DefaultOption>
          <Help><![CDATA[Instant Off: Turn off system immediately in legacy OS.
4 Seconds Override: Turn off system after depressed for 4 seconds.]]></Help>
        </Information>
      </Setting>
      <Setting name="Throttle on Power Fail" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[To decrease system power by throttling CPU frequency when one power supply is failed]]></Help>
        </Information>
      </Setting>
      <Setting name="System Firmware Progress Log" selectedOption="Enabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[Enable System Firmware Progress Log.]]></Help>
        </Information>
      </Setting>
    </Menu>
    <Menu name="CPU Configuration">
      <Information>
        <Help><![CDATA[CPU Configuration]]></Help>
      </Information>
      <Subtitle>Processor Configuration</Subtitle>
      <Subtitle>--------------------------------------------------</Subtitle>
      <Text>Processor BSP Revision(50654 - SKX U0)</Text>
      <Text>Processor Socket(CPU1      |  CPU2    )</Text>
      <Text>Processor ID(00050654* |  00050654 )</Text>
      <Text>Processor Frequency(2.100GHz  |  2.100GHz)</Text>
      <Text>Processor Max Ratio(     15H  |  15H)</Text>
      <Text>Processor Min Ratio(     08H  |  08H)</Text>
      <Text>Microcode Revision(02000057  |  02000057)</Text>
      <Text>L1 Cache RAM(    64KB  |      64KB)</Text>
      <Text>L2 Cache RAM(  1024KB  |    1024KB)</Text>
      <Text>L3 Cache RAM( 11264KB  |   11264KB)</Text>
      <Subtitle>Processor 0 Version</Subtitle>
      <Subtitle>Intel(R) Xeon(R) Silver 4110 CPU @ 2.10GHz</Subtitle>
      <Subtitle>Processor 1 Version</Subtitle>
      <Subtitle>Intel(R) Xeon(R) Silver 4110 CPU @ 2.10GHz</Subtitle>
      <Subtitle></Subtitle>
      <Setting name="Hyper-Threading [ALL]" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Disable</Option>
            <Option value="0">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Enables Hyper Threading (Software Method to Enable/Disable Logical Processor threads.]]></Help>
        </Information>
      </Setting>
      <Setting name="Cores Enabled" numericValue="0" type="Numeric">
        <Information>
          <MaxValue>28</MaxValue>
          <MinValue>0</MinValue>
          <StepSize>1</StepSize>
          <DefaultValue>0</DefaultValue>
          <Help><![CDATA[Number of Cores to Enable in each Processor Package. 0 means all cores. Total 8 cores available in each CPU package.]]></Help>
        </Information>
      </Setting>
      <Setting name="Monitor/Mwait" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Enable or Disable the Monitor/Mwait instruction]]></Help>
        </Information>
      </Setting>
      <Setting name="Execute Disable Bit" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[When disabled, forces the XD feature flag to always return 0.]]></Help>
        </Information>
      </Setting>
      <Setting name="Intel Virtualization Technology" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[When enabled, a VMM can utilize the additional hardware capabilities provided by Vanderpool Technology]]></Help>
        </Information>
      </Setting>
      <Setting name="PPIN Control" selectedOption="Unlock/Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Unlock/Disable</Option>
            <Option value="1">Unlock/Enable</Option>
          </AvailableOptions>
          <DefaultOption>Unlock/Enable</DefaultOption>
          <Help><![CDATA[When Protected Processor Inventory Number (PPIN) is enabled, the processor will return a 64-bit ID number via the PPIN MSR.]]></Help>
        </Information>
      </Setting>
      <Setting name="Hardware Prefetcher" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Enable</Option>
            <Option value="0">Disable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[To turn on/off the Mid Level Cache (L2) streamer prefetcher.]]></Help>
        </Information>
      </Setting>
      <Setting name="Adjacent Cache Prefetch" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Enable</Option>
            <Option value="0">Disable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[To turn on/off prefetching of adjacent cache lines.]]></Help>
        </Information>
      </Setting>
      <Setting name="DCU Streamer Prefetcher" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Enable</Option>
            <Option value="0">Disable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Enable prefetch of next L1 Data line based upon multiple loads in same cache line.]]></Help>
        </Information>
      </Setting>
      <Setting name="DCU IP Prefetcher" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Enable</Option>
            <Option value="0">Disable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Enable prefetch of next L1 line based upon sequential load history.]]></Help>
        </Information>
      </Setting>
      <Setting name="LLC Prefetch" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If this feature is set to Enable, LLC (hardware cache) prefetching on all threads will be supported.]]></Help>
        </Information>
      </Setting>
      <Setting name="Extended APIC" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[Enable/disable extended APIC support]]></Help>
        </Information>
      </Setting>
      <Setting name="AES-NI" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Enable/disable AES-NI support]]></Help>
        </Information>
      </Setting>
      <Menu name="Advanced Power Management Configuration">
        <Information>
          <Help><![CDATA[Displays and provides option to change the Power Management Settings]]></Help>
        </Information>
        <Subtitle>Advanced Power Management Configuration</Subtitle>
        <Subtitle>--------------------------------------------------</Subtitle>
        <Setting name="Power Technology" selectedOption="Disable" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">Disable</Option>
              <Option value="1">Energy Efficient</Option>
              <Option value="2">Custom</Option>
            </AvailableOptions>
            <DefaultOption>Energy Efficient</DefaultOption>
            <Help><![CDATA[Enable processor power management features.]]></Help>
          </Information>
        </Setting>
        <Setting name="Power Performance Tuning" selectedOption="OS Controls EPB" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">OS Controls EPB</Option>
              <Option value="1">BIOS Controls EPB</Option>
            </AvailableOptions>
            <DefaultOption>OS Controls EPB</DefaultOption>
            <Help><![CDATA[Selects whether BIOS or Operatiing System chooses eneryg performance bias tuning.]]></Help>
            <WorkIf><![CDATA[ ( 2 == Power Technology )  and  ( 2 != Hardware P-States ) ]]></WorkIf>
          </Information>
        </Setting>
        <Setting name="ENERGY_PERF_BIAS_CFG mode" selectedOption="Balanced Performance" type="Option">
          <Information>
            <AvailableOptions>
              <!--Option ValidIf:  ( 1 == Monitor/Mwait ) -->
              <Option value="3">Maximum Performance</Option>
              <Option value="0">Performance</Option>
              <Option value="7">Balanced Performance</Option>
              <Option value="8">Balanced Power</Option>
              <Option value="15">Power</Option>
            </AvailableOptions>
            <DefaultOption>Balanced Performance</DefaultOption>
            <Help><![CDATA[Set Energy Performance BIAS, which overrides OS setting.]]></Help>
            <WorkIf><![CDATA[ ( 2 == Power Technology )  and  ( 0 != Power Performance Tuning ) ]]></WorkIf>
          </Information>
        </Setting>
        <Menu name="CPU P State Control">
          <Information>
            <Help><![CDATA[P State Control Configuration Sub Menu, include Turbo, XE and etc.]]></Help>
            <WorkIf><![CDATA[  2 == Power Technology  ]]></WorkIf>
          </Information>
          <Subtitle>CPU P State Control</Subtitle>
          <Subtitle></Subtitle>
          <Setting name="SpeedStep (P-States)" selectedOption="Disable" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
              </AvailableOptions>
              <DefaultOption>Enable</DefaultOption>
              <Help><![CDATA[EIST allows the processor to dynamically adjust frequency and voltage based on power versus performance needs.]]></Help>
            </Information>
          </Setting>
          <Setting name="EIST PSD Function" selectedOption="HW_ALL" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">HW_ALL</Option>
                <Option value="1">SW_ALL</Option>
                <Option value="2">SW_ANY</Option>
              </AvailableOptions>
              <DefaultOption>HW_ALL</DefaultOption>
              <Help><![CDATA[Determine how ACPI-aware OS coordinates P-State transitions between logical processors.]]></Help>
              <WorkIf><![CDATA[  0 != SpeedStep (P-States)  ]]></WorkIf>
            </Information>
          </Setting>
          <Setting name="Turbo Mode" selectedOption="Enable" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
              </AvailableOptions>
              <DefaultOption>Enable</DefaultOption>
              <Help><![CDATA[Turbo Mode opportunistically and automatically allows processor cores to run faster than the marked frequency if the physical processor is operation below power, temperature and current specification limits]]></Help>
              <WorkIf><![CDATA[  0 != SpeedStep (P-States)  ]]></WorkIf>
            </Information>
          </Setting>
        </Menu>
        <Menu name="Hardware PM State Control">
          <Information>
            <Help><![CDATA[Hardware P-State setting]]></Help>
            <WorkIf><![CDATA[  2 == Power Technology  ]]></WorkIf>
          </Information>
          <Subtitle>Hardware PM State Control</Subtitle>
          <Subtitle></Subtitle>
          <Setting name="Hardware P-States" selectedOption="Disable" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Native Mode</Option>
                <Option value="2">Out of Band Mode</Option>
                <Option value="3">Native Mode with No Legacy Support</Option>
              </AvailableOptions>
              <DefaultOption>Disable</DefaultOption>
              <Help><![CDATA[If set to Disable, hardware will choose a P-state setting for the system based on an OS request.
If set to Native Mode, hardware will choose a P-state setting based on OS guidance.
If set to Native Mode with No Legacy Support, hardware will choose a P-state setting independently without OS guidance.
If set to Out of Band Mode, hardware autonomously choose a P-state without OS guidance.]]></Help>
            </Information>
          </Setting>
        </Menu>
        <Menu name="CPU C State Control">
          <Information>
            <Help><![CDATA[CPU C State setting]]></Help>
            <WorkIf><![CDATA[  2 == Power Technology  ]]></WorkIf>
          </Information>
          <Subtitle>CPU C State Control</Subtitle>
          <Subtitle></Subtitle>
          <Setting name="Autonomous Core C-State" selectedOption="Disable" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
              </AvailableOptions>
              <DefaultOption>Disable</DefaultOption>
              <Help><![CDATA[Select Enable to support Autonomous Core C-State control which will allow the processor core to control its C-State setting automatically and independently.]]></Help>
            </Information>
          </Setting>
          <Setting name="CPU C6 report" selectedOption="Disable" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
                <Option value="255">Auto</Option>
              </AvailableOptions>
              <DefaultOption>Auto</DefaultOption>
              <Help><![CDATA[Enable/Disable CPU C6(ACPI C3) report to OS. Recommanded to be enabled.]]></Help>
              <WorkIf><![CDATA[  1 != Autonomous Core C-State  ]]></WorkIf>
            </Information>
          </Setting>
          <Setting name="Enhanced Halt State (C1E)" selectedOption="Disable" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
              </AvailableOptions>
              <DefaultOption>Enable</DefaultOption>
              <Help><![CDATA[Enable/Disable processors to transitions the voltage associated with Max Efficiency Ratio.Takes effect after reboot.]]></Help>
              <WorkIf><![CDATA[  1 != Autonomous Core C-State  ]]></WorkIf>
            </Information>
          </Setting>
        </Menu>
        <Menu name="Package C State Control">
          <Information>
            <Help><![CDATA[Package C State setting]]></Help>
            <WorkIf><![CDATA[  2 == Power Technology  ]]></WorkIf>
          </Information>
          <Subtitle>Package C State Control</Subtitle>
          <Subtitle></Subtitle>
          <Setting name="Package C State" selectedOption="C0/C1 state" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">C0/C1 state</Option>
                <Option value="1">C2 state</Option>
                <Option value="2">C6(non Retention) state</Option>
                <Option value="3">C6(Retention) state</Option>
                <Option value="7">No Limit</Option>
                <Option value="255">Auto</Option>
              </AvailableOptions>
              <DefaultOption>Auto</DefaultOption>
              <Help><![CDATA[Limit the lowest package level C-State to processors. Lower package C-State lower processor power consumption upon idle.]]></Help>
            </Information>
          </Setting>
        </Menu>
        <Menu name="CPU T State Control">
          <Information>
            <Help><![CDATA[CPU T State setting]]></Help>
            <WorkIf><![CDATA[  2 == Power Technology  ]]></WorkIf>
          </Information>
          <Subtitle>CPU T State Control</Subtitle>
          <Subtitle></Subtitle>
          <Setting name="Software Controlled T-States" selectedOption="Disable" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
              </AvailableOptions>
              <DefaultOption>Disable</DefaultOption>
              <Help><![CDATA[Enable/Disable CPU throttling by OS. Throttling reduces power consumption]]></Help>
            </Information>
          </Setting>
        </Menu>
      </Menu>
    </Menu>
    <Menu name="Chipset Configuration">
      <Information>
        <Help><![CDATA[System Chipset configuration.]]></Help>
      </Information>
      <Subtitle>WARNING: Setting wrong values in below sections may cause</Subtitle>
      <Subtitle>         system to malfunction.</Subtitle>
      <Menu name="North Bridge">
        <Information>
          <Help><![CDATA[North Bridge Parameters]]></Help>
        </Information>
        <Menu name="UPI Configuration">
          <Information>
            <Help><![CDATA[Displays and provides option to change the UPI Settings]]></Help>
          </Information>
          <Subtitle>UPI Configuration</Subtitle>
          <Subtitle>--------------------------------------------------</Subtitle>
          <Text>Number of CPU(2)</Text>
          <Text>Number of Active UPI Link(2)</Text>
          <Text>Current UPI Link Speed(Fast)</Text>
          <Text>Current UPI Link Frequency(9.6 GT/s)</Text>
          <Text>UPI Global MMIO Low Base / Limit(90000000 / FBFFFFFF)</Text>
          <Text>UPI Global MMIO High Base / Limit(0000000000000000 / 00000000FFFFFFFF)</Text>
          <Text>UPI Pci-e Configuration Base / Size(80000000 / 10000000)</Text>
          <Setting name="Degrade Precedence" selectedOption="Topology Precedence" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Topology Precedence</Option>
                <Option value="1">Feature Precedence</Option>
              </AvailableOptions>
              <DefaultOption>Topology Precedence</DefaultOption>
              <Help><![CDATA[Use this feature to select the degrading precedence option for Ultra Path Interconnect connections. Select Topology Precedent to degrade UPI features if system options are in conflict. Select Feature Precedent to degrade UPI topology if system options are in conflict.]]></Help>
            </Information>
          </Setting>
          <Setting name="Link L0p Enable" selectedOption="Auto" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
                <Option value="2">Auto</Option>
              </AvailableOptions>
              <DefaultOption>Auto</DefaultOption>
              <Help><![CDATA[Enable/Disable UPI to enter L0p state for power saving.]]></Help>
            </Information>
          </Setting>
          <Setting name="Link L1 Enable" selectedOption="Auto" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
                <Option value="2">Auto</Option>
              </AvailableOptions>
              <DefaultOption>Auto</DefaultOption>
              <Help><![CDATA[Enable/Disable UPI to enter L1 state for power saving.]]></Help>
            </Information>
          </Setting>
          <Setting name="IO Directory Cache (IODC)" selectedOption="Auto" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <!--Option ValidIf:  ( 1 == Snoopy mode for AD ) -->
                <Option value="1">Auto</Option>
              </AvailableOptions>
              <DefaultOption>Auto</DefaultOption>
              <Help><![CDATA[Select Enable for the IODC (I/O Directory Cache) to generate snoops instead of generating memory lockups for remote IIO (InvIToM) and/or WCiLF (Cores). Select Auto for the IODC to generate snoops (instead of memory lockups) for WCiLF (Cores).]]></Help>
            </Information>
          </Setting>
          <Setting name="SNC" selectedOption="Disable" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
                <Option value="2">Auto</Option>
              </AvailableOptions>
              <DefaultOption>Disable</DefaultOption>
              <Help><![CDATA[Sub NUMA Clustering (SNC) is a feature that breaks up the Last Level Cache (LLC) into clusters based on address range. Each cluster is connected to a subset of the memory controller. Enabling SNC improves average latency and reduces memory access congestion to achieve higher performance. Select Auto for 1-cluster or 2-clusters depending on IMC interleave. Select Enable for Full SNC (2-clusters and 1-way IMC interleave).]]></Help>
            </Information>
          </Setting>
          <Setting name="XPT Prefetch" selectedOption="Disable" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
              </AvailableOptions>
              <DefaultOption>Disable</DefaultOption>
              <Help><![CDATA[XPT Prefetch speculatively makes a copy to the memory controller of a read request being sent to the LLC. If the read request maps to the local memory address and the recent memory reads are likely to miss the LLC, a speculative read is sent to the local memory controller]]></Help>
            </Information>
          </Setting>
          <Setting name="KTI Prefetch" selectedOption="Enable" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
              </AvailableOptions>
              <DefaultOption>Enable</DefaultOption>
              <Help><![CDATA[KTI Pretech enables memory read to start early on a DDR bus, where the KTI Rx path will directly create a Memory Speculative Read command to the memory controller.]]></Help>
            </Information>
          </Setting>
          <Setting name="Local/Remote Threshold" selectedOption="Auto" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Auto</Option>
                <Option value="2">Low</Option>
                <Option value="3">Medium</Option>
                <Option value="4">High</Option>
              </AvailableOptions>
              <DefaultOption>Auto</DefaultOption>
              <Help><![CDATA[This feature allows the user to set the threshold for the Interrupt Request (IRQ) signal, which handles hardware interruptions.]]></Help>
            </Information>
          </Setting>
          <Setting name="Stale AtoS" selectedOption="Auto" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
                <Option value="2">Auto</Option>
              </AvailableOptions>
              <DefaultOption>Auto</DefaultOption>
              <Help><![CDATA[This feature optimizes A to S directory. When all snoop responses found in directory A are found to be RspI, then all data is moved to directory S and is returned in S-state.]]></Help>
            </Information>
          </Setting>
          <Setting name="LLC Dead Line Alloc" selectedOption="Enable" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
                <Option value="2">Auto</Option>
              </AvailableOptions>
              <DefaultOption>Enable</DefaultOption>
              <Help><![CDATA[Select Enable to optimally fill dead lines in LLC. Select Disable to never fill dead lines in LLC.]]></Help>
            </Information>
          </Setting>
          <Setting name="Isoc Mode" selectedOption="Auto" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
                <Option value="2">Auto</Option>
              </AvailableOptions>
              <DefaultOption>Auto</DefaultOption>
              <Help><![CDATA[Enable/Disable the isochronous mode to reduce/increase the credits available for memory traffic. Workstation & HEDT require Isoc enabled for audio and media performance.]]></Help>
            </Information>
          </Setting>
        </Menu>
        <Menu name="Memory Configuration">
          <Information>
            <Help><![CDATA[Displays and provides option to change the Memory Settings]]></Help>
          </Information>
          <Subtitle></Subtitle>
          <Subtitle>--------------------------------------------------</Subtitle>
          <Subtitle>Integrated Memory Controller (iMC)</Subtitle>
          <Subtitle>--------------------------------------------------</Subtitle>
          <Subtitle></Subtitle>
          <Setting name="Enforce POR" selectedOption="POR" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">POR</Option>
                <Option value="2">Disable</Option>
              </AvailableOptions>
              <DefaultOption>POR</DefaultOption>
              <Help><![CDATA[Select POR to enforce POR restrictions for DDR4 frequency and voltage programming]]></Help>
            </Information>
          </Setting>
          <Setting name="PPR Type" selectedOption="Auto" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="3">Auto</Option>
                <Option value="1">Hard PPR</Option>
                <Option value="2">Soft PPR</Option>
                <Option value="0">PPR Disabled</Option>
              </AvailableOptions>
              <DefaultOption>Auto</DefaultOption>
              <Help><![CDATA[Selects Post Package Repair Type - Hard / Soft / Disabled. Auto - Sets it to the MRC default setting; current default is Disabled.]]></Help>
            </Information>
          </Setting>
          <Setting name="Memory Frequency" selectedOption="Auto" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Auto</Option>
                <Option value="9">1866</Option>
                <Option value="10">2000</Option>
                <Option value="11">2133</Option>
                <Option value="13">2400</Option>
                <Option value="15">2666</Option>
                <Option value="17">2933</Option>
              </AvailableOptions>
              <DefaultOption>Auto</DefaultOption>
              <Help><![CDATA[Restrict maximum memory frequency below enforced POR.]]></Help>
            </Information>
          </Setting>
          <Setting name="Data Scrambling for DDR4" selectedOption="Auto" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="2">Auto</Option>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
              </AvailableOptions>
              <DefaultOption>Auto</DefaultOption>
              <Help><![CDATA[This feature improves detection of DDR4 memory address line errors and reduces the probability of occurrence.]]></Help>
            </Information>
          </Setting>
          <Setting name="tCCD_L Relaxation" selectedOption="Auto" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Auto</Option>
              </AvailableOptions>
              <DefaultOption>Auto</DefaultOption>
              <Help><![CDATA[If enabled, the tCCD_L is overridden by the SPD. Otherwise, it's enforced based on memory frequency]]></Help>
            </Information>
          </Setting>
          <Setting name="tRWSR Relaxation" selectedOption="Disable" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
              </AvailableOptions>
              <DefaultOption>Disable</DefaultOption>
              <Help><![CDATA[If enabled, the tRWSR is used worst case value which will improve the RX signal margin.]]></Help>
            </Information>
          </Setting>
          <Setting name="2x Refresh" selectedOption="Auto" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Auto</Option>
                <Option value="1">Enable</Option>
              </AvailableOptions>
              <DefaultOption>Auto</DefaultOption>
              <Help><![CDATA[This option allows the user to select 2X memory refresh mode.]]></Help>
            </Information>
          </Setting>
          <Setting name="Page Policy" selectedOption="Auto" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="3">Auto</Option>
                <Option value="1">Closed</Option>
                <Option value="2">Adaptive</Option>
              </AvailableOptions>
              <DefaultOption>Auto</DefaultOption>
              <Help><![CDATA[This feature allows the user to determine the desired page mode for IMC. When Auto is selected, the memory controller will close or open pages based on the current operation. Closed policy closes that page after reading or writing. Adaptive is similar to open page policy, but can be dynamically modified.]]></Help>
            </Information>
          </Setting>
          <Setting name="IMC Interleaving" selectedOption="Auto" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Auto</Option>
                <Option value="1">1-way Interleave</Option>
                <Option value="2">2-way Interleave</Option>
              </AvailableOptions>
              <DefaultOption>Auto</DefaultOption>
              <Help><![CDATA[This feature allows the user to configure Integrated Memory Controller (IMC) Interleaving settings.]]></Help>
            </Information>
          </Setting>
          <Menu name="Memory Topology">
            <Information>
              <Help><![CDATA[Displays memory topology with Dimm population information]]></Help>
            </Information>
            <Subtitle>--------------------------------------------------</Subtitle>
            <Subtitle>P1 DIMMA1:  2400MT/s Micron Technology DRx8 8GB RDIMM</Subtitle>
            <Subtitle>P1 DIMMB1:  2400MT/s Micron Technology DRx8 8GB RDIMM</Subtitle>
            <Subtitle>P1 DIMMC1:  2400MT/s Micron Technology DRx8 8GB RDIMM</Subtitle>
            <Subtitle>P1 DIMMD1:  2400MT/s Micron Technology DRx8 8GB RDIMM</Subtitle>
            <Subtitle>P1 DIMME1:  2400MT/s Micron Technology DRx8 8GB RDIMM</Subtitle>
            <Subtitle>P1 DIMMF1:  2400MT/s Micron Technology DRx8 8GB RDIMM</Subtitle>
            <Subtitle>P2 DIMMA1:  2400MT/s Micron Technology DRx8 8GB RDIMM</Subtitle>
            <Subtitle>P2 DIMMB1:  2400MT/s Micron Technology DRx8 8GB RDIMM</Subtitle>
            <Subtitle>P2 DIMMC1:  2400MT/s Micron Technology DRx8 8GB RDIMM</Subtitle>
            <Subtitle>P2 DIMMD1:  2400MT/s Micron Technology DRx8 8GB RDIMM</Subtitle>
            <Subtitle>P2 DIMME1:  2400MT/s Micron Technology DRx8 8GB RDIMM</Subtitle>
            <Subtitle>P2 DIMMF1:  2400MT/s Micron Technology DRx8 8GB RDIMM</Subtitle>
            <Subtitle></Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
            <Subtitle>     </Subtitle>
          </Menu>
          <Menu name="Memory RAS Configuration">
            <Information>
              <Help><![CDATA[Displays and provides option to change the Memory Ras Settings]]></Help>
            </Information>
            <Subtitle></Subtitle>
            <Subtitle>--------------------------------------------------</Subtitle>
            <Subtitle>Memory RAS Configuration Setup</Subtitle>
            <Subtitle>--------------------------------------------------</Subtitle>
            <Subtitle></Subtitle>
            <Setting name="Static Virtual Lockstep Mode" selectedOption="Disable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">Disable</Option>
                  <Option value="3">Enable</Option>
                </AvailableOptions>
                <DefaultOption>Disable</DefaultOption>
                <Help><![CDATA[Select Enable to support Static Virtual Lockstep mode to enhance memory performance.]]></Help>
              </Information>
            </Setting>
            <Setting name="Mirror mode" selectedOption="Disable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">Disable</Option>
                  <Option value="1">Mirror Mode 1LM</Option>
                  <Option value="2">Mirror Mode 2LM</Option>
                </AvailableOptions>
                <DefaultOption>Disable</DefaultOption>
                <Help><![CDATA[Select Enable to set all 1LM/2LM memory installed in the system on the mirror mode, which will create a duplicate copy of data stored in the memory to increase memory security, but it will reduce the memory capacity into half.]]></Help>
              </Information>
            </Setting>
            <Setting name="Memory Rank Sparing" selectedOption="Disable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">Disable</Option>
                  <Option value="1">Enable</Option>
                </AvailableOptions>
                <DefaultOption>Disable</DefaultOption>
                <Help><![CDATA[Rank sparing enables a failing rank to be replaced by ranks installed in an unoccupied space.]]></Help>
                <WorkIf><![CDATA[   ( 1 != Mirror mode )  &&  ( 1 != 0 )   ]]></WorkIf>
              </Information>
            </Setting>
            <Setting name="Multi Rank Sparing" selectedOption="Two Rank" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="1">One Rank</Option>
                  <Option value="2">Two Rank</Option>
                </AvailableOptions>
                <DefaultOption>Two Rank</DefaultOption>
                <Help><![CDATA[Rank sparing enables a failing rank to be replaced by ranks installed in an unoccupied space.]]></Help>
                <WorkIf><![CDATA[  0 != Memory Rank Sparing  ]]></WorkIf>
              </Information>
            </Setting>
            <Setting name="Correctable Error Threshold" numericValue="100" type="Numeric">
              <Information>
                <MaxValue>32767</MaxValue>
                <MinValue>1</MinValue>
                <StepSize>1</StepSize>
                <DefaultValue>100</DefaultValue>
                <Help><![CDATA[Correctable Error Threshold (0x01 - 0x7fff) used for sparing, tagging, and leaky bucket]]></Help>
              </Information>
            </Setting>
            <Setting name="SDDC" selectedOption="Disable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">Disable</Option>
                  <Option value="1">Enable</Option>
                </AvailableOptions>
                <DefaultOption>Disable</DefaultOption>
                <Help><![CDATA[Enable/Disable SDDC. Not supported when AEP dimm present!]]></Help>
              </Information>
            </Setting>
            <Setting name="Patrol Scrub" selectedOption="Enable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">Disable</Option>
                  <Option value="1">Enable</Option>
                </AvailableOptions>
                <DefaultOption>Enable</DefaultOption>
                <Help><![CDATA[Enable or disable the ability to proactively searches the system memory, repairing correctable errors.]]></Help>
              </Information>
            </Setting>
            <Setting name="Patrol Scrub Interval" numericValue="24" type="Numeric">
              <Information>
                <MaxValue>24</MaxValue>
                <MinValue>0</MinValue>
                <StepSize>0</StepSize>
                <DefaultValue>24</DefaultValue>
                <Help><![CDATA[Selects the number of hours (1-24) required to complete full scrub. A value of zero means auto!]]></Help>
                <WorkIf><![CDATA[  0 != Patrol Scrub  ]]></WorkIf>
              </Information>
            </Setting>
          </Menu>
        </Menu>
        <Menu name="IIO Configuration">
          <Information>
            <Help><![CDATA[Displays and provides option to change the IIO Settings]]></Help>
          </Information>
          <Subtitle>IIO Configuration</Subtitle>
          <Subtitle>--------------------------------------------------</Subtitle>
          <Subtitle></Subtitle>
          <Setting name="EV DFX Features" selectedOption="Disable" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="0">Disable</Option>
                <Option value="1">Enable</Option>
              </AvailableOptions>
              <DefaultOption>Disable</DefaultOption>
              <Help><![CDATA[Enable/Disable DFx logics(Design for Debug, validation, Test) for Intel Electrical Validation tools.]]></Help>
            </Information>
          </Setting>
          <Menu name="CPU1 Configuration">
            <Information>
              <Help><![CDATA[]]></Help>
            </Information>
            <Setting name="IOU0 (IIO PCIe Br1)" order="1" selectedOption="Auto" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">x4x4x4x4</Option>
                  <Option value="1">x4x4x8</Option>
                  <Option value="2">x8x4x4</Option>
                  <Option value="3">x8x8</Option>
                  <Option value="4">x16</Option>
                  <Option value="255">Auto</Option>
                </AvailableOptions>
                <DefaultOption>Auto</DefaultOption>
                <Help><![CDATA[Selects PCIe port Bifurcation for selected slot(s)]]></Help>
              </Information>
            </Setting>
            <Setting name="IOU1 (IIO PCIe Br2)" order="1" selectedOption="Auto" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">x4x4x4x4</Option>
                  <Option value="1">x4x4x8</Option>
                  <Option value="2">x8x4x4</Option>
                  <Option value="3">x8x8</Option>
                  <Option value="4">x16</Option>
                  <Option value="255">Auto</Option>
                </AvailableOptions>
                <DefaultOption>Auto</DefaultOption>
                <Help><![CDATA[Selects PCIe port Bifurcation for selected slot(s)]]></Help>
              </Information>
            </Setting>
            <Setting name="IOU2 (IIO PCIe Br3)" order="1" selectedOption="Auto" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">x4x4x4x4</Option>
                  <Option value="1">x4x4x8</Option>
                  <Option value="2">x8x4x4</Option>
                  <Option value="3">x8x8</Option>
                  <Option value="4">x16</Option>
                  <Option value="255">Auto</Option>
                </AvailableOptions>
                <DefaultOption>Auto</DefaultOption>
                <Help><![CDATA[Selects PCIe port Bifurcation for selected slot(s)]]></Help>
              </Information>
            </Setting>
            <Menu name="CPU1 PcieBr0D00F0 - Port 0/DMI">
              <Information>
                <Help><![CDATA[Settings related to PCI Express PortS (0/1A/1B/1C/1D/2A/2B/2C/2D/3A/3B/3C/3D/4A/5A) ]]></Help>
              </Information>
              <Subtitle>CPU1 PcieBr0D00F0 - Port 0/DMI</Subtitle>
              <Subtitle>--------------------------------------------------</Subtitle>
              <Subtitle></Subtitle>
              <Setting name="Link Speed" order="1" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Auto</Option>
                    <Option value="1">Gen 1 (2.5 GT/s)</Option>
                    <Option value="2">Gen 2 (5 GT/s)</Option>
                    <Option value="3">Gen 3 (8 GT/s)</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Choose Link Speed for this PCIe port]]></Help>
                </Information>
              </Setting>
              <Text>PCI-E Port Link Status(Linked as x4)</Text>
              <Text>PCI-E Port Link Max(Max Width x4)</Text>
              <Text>PCI-E Port Link Speed(Gen 3 (8.0 GT/s))</Text>
              <Setting name="PCI-E Port Max Payload Size" order="1" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">128B</Option>
                    <Option value="1">256B</Option>
                    <Option value="2">Auto</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Set Maxpayload size to 256B if possible]]></Help>
                </Information>
              </Setting>
            </Menu>
            <Menu name="CPU1 PcieBr1D00F0 - Port 1A">
              <Information>
                <Help><![CDATA[Settings related to PCI Express PortS (0/1A/1B/1C/1D/2A/2B/2C/2D/3A/3B/3C/3D/4A/5A) ]]></Help>
              </Information>
              <Subtitle>CPU1 PcieBr1D00F0 - Port 1A</Subtitle>
              <Subtitle>--------------------------------------------------</Subtitle>
              <Subtitle></Subtitle>
              <Setting name="Link Speed" order="2" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Auto</Option>
                    <Option value="1">Gen 1 (2.5 GT/s)</Option>
                    <Option value="2">Gen 2 (5 GT/s)</Option>
                    <Option value="3">Gen 3 (8 GT/s)</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Choose Link Speed for this PCIe port]]></Help>
                </Information>
              </Setting>
              <Text>PCI-E Port Link Status(Linked as x8)</Text>
              <Text>PCI-E Port Link Max(Max Width x16)</Text>
              <Text>PCI-E Port Link Speed(Gen 3 (8.0 GT/s))</Text>
              <Setting name="PCI-E Port Max Payload Size" order="2" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">128B</Option>
                    <Option value="1">256B</Option>
                    <Option value="2">Auto</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Set Maxpayload size to 256B if possible]]></Help>
                </Information>
              </Setting>
            </Menu>
            <Menu name="CPU1 PcieBr2D00F0 - Port 2A">
              <Information>
                <Help><![CDATA[Settings related to PCI Express PortS (0/1A/1B/1C/1D/2A/2B/2C/2D/3A/3B/3C/3D/4A/5A) ]]></Help>
              </Information>
              <Subtitle>CPU1 PcieBr2D00F0 - Port 2A</Subtitle>
              <Subtitle>--------------------------------------------------</Subtitle>
              <Subtitle></Subtitle>
              <Setting name="Link Speed" order="3" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Auto</Option>
                    <Option value="1">Gen 1 (2.5 GT/s)</Option>
                    <Option value="2">Gen 2 (5 GT/s)</Option>
                    <Option value="3">Gen 3 (8 GT/s)</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Choose Link Speed for this PCIe port]]></Help>
                </Information>
              </Setting>
              <Text>PCI-E Port Link Status(Linked as x4)</Text>
              <Text>PCI-E Port Link Max(Max Width x16)</Text>
              <Text>PCI-E Port Link Speed(Gen 2 (5.0 GT/s))</Text>
              <Setting name="PCI-E Port Max Payload Size" order="3" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">128B</Option>
                    <Option value="1">256B</Option>
                    <Option value="2">Auto</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Set Maxpayload size to 256B if possible]]></Help>
                </Information>
              </Setting>
            </Menu>
            <Menu name="CPU1 PcieBr3D00F0 - Port 3A">
              <Information>
                <Help><![CDATA[Settings related to PCI Express PortS (0/1A/1B/1C/1D/2A/2B/2C/2D/3A/3B/3C/3D/4A/5A) ]]></Help>
              </Information>
              <Subtitle>CPU1 PcieBr3D00F0 - Port 3A</Subtitle>
              <Subtitle>--------------------------------------------------</Subtitle>
              <Subtitle></Subtitle>
              <Setting name="Link Speed" order="4" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Auto</Option>
                    <Option value="1">Gen 1 (2.5 GT/s)</Option>
                    <Option value="2">Gen 2 (5 GT/s)</Option>
                    <Option value="3">Gen 3 (8 GT/s)</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Choose Link Speed for this PCIe port]]></Help>
                </Information>
              </Setting>
              <Text>PCI-E Port Link Status(Link Did Not Train)</Text>
              <Text>PCI-E Port Link Max(Max Width x8)</Text>
              <Text>PCI-E Port Link Speed(Link Did Not Train)</Text>
              <Setting name="PCI-E Port Max Payload Size" order="4" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">128B</Option>
                    <Option value="1">256B</Option>
                    <Option value="2">Auto</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Set Maxpayload size to 256B if possible]]></Help>
                </Information>
              </Setting>
            </Menu>
            <Menu name="CPU1 PcieBr3D02F0 - Port 3C">
              <Information>
                <Help><![CDATA[Settings related to PCI Express PortS (0/1A/1B/1C/1D/2A/2B/2C/2D/3A/3B/3C/3D/4A/5A) ]]></Help>
              </Information>
              <Subtitle>CPU1 PcieBr3D02F0 - Port 3C</Subtitle>
              <Subtitle>--------------------------------------------------</Subtitle>
              <Subtitle></Subtitle>
              <Setting name="Link Speed" order="5" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Auto</Option>
                    <Option value="1">Gen 1 (2.5 GT/s)</Option>
                    <Option value="2">Gen 2 (5 GT/s)</Option>
                    <Option value="3">Gen 3 (8 GT/s)</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Choose Link Speed for this PCIe port]]></Help>
                </Information>
              </Setting>
              <Text>PCI-E Port Link Status(Link Did Not Train)</Text>
              <Text>PCI-E Port Link Max(Max Width x4)</Text>
              <Text>PCI-E Port Link Speed(Link Did Not Train)</Text>
              <Setting name="PCI-E Port Max Payload Size" order="5" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">128B</Option>
                    <Option value="1">256B</Option>
                    <Option value="2">Auto</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Set Maxpayload size to 256B if possible]]></Help>
                </Information>
              </Setting>
            </Menu>
            <Menu name="CPU1 PcieBr3D03F0 - Port 3D">
              <Information>
                <Help><![CDATA[Settings related to PCI Express PortS (0/1A/1B/1C/1D/2A/2B/2C/2D/3A/3B/3C/3D/4A/5A) ]]></Help>
              </Information>
              <Subtitle>CPU1 PcieBr3D03F0 - Port 3D</Subtitle>
              <Subtitle>--------------------------------------------------</Subtitle>
              <Subtitle></Subtitle>
              <Setting name="Link Speed" order="6" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Auto</Option>
                    <Option value="1">Gen 1 (2.5 GT/s)</Option>
                    <Option value="2">Gen 2 (5 GT/s)</Option>
                    <Option value="3">Gen 3 (8 GT/s)</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Choose Link Speed for this PCIe port]]></Help>
                </Information>
              </Setting>
              <Text>PCI-E Port Link Status(Link Did Not Train)</Text>
              <Text>PCI-E Port Link Max(Max Width x4)</Text>
              <Text>PCI-E Port Link Speed(Link Did Not Train)</Text>
              <Setting name="PCI-E Port Max Payload Size" order="6" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">128B</Option>
                    <Option value="1">256B</Option>
                    <Option value="2">Auto</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Set Maxpayload size to 256B if possible]]></Help>
                </Information>
              </Setting>
            </Menu>
          </Menu>
          <Menu name="CPU2 Configuration">
            <Information>
              <Help><![CDATA[]]></Help>
            </Information>
            <Setting name="IOU0 (IIO PCIe Br1)" order="2" selectedOption="Auto" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">x4x4x4x4</Option>
                  <Option value="1">x4x4x8</Option>
                  <Option value="2">x8x4x4</Option>
                  <Option value="3">x8x8</Option>
                  <Option value="4">x16</Option>
                  <Option value="255">Auto</Option>
                </AvailableOptions>
                <DefaultOption>Auto</DefaultOption>
                <Help><![CDATA[Selects PCIe port Bifurcation for selected slot(s)]]></Help>
              </Information>
            </Setting>
            <Setting name="IOU1 (IIO PCIe Br2)" order="2" selectedOption="Auto" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">x4x4x4x4</Option>
                  <Option value="1">x4x4x8</Option>
                  <Option value="2">x8x4x4</Option>
                  <Option value="3">x8x8</Option>
                  <Option value="4">x16</Option>
                  <Option value="255">Auto</Option>
                </AvailableOptions>
                <DefaultOption>Auto</DefaultOption>
                <Help><![CDATA[Selects PCIe port Bifurcation for selected slot(s)]]></Help>
              </Information>
            </Setting>
            <Setting name="IOU2 (IIO PCIe Br3)" order="2" selectedOption="Auto" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">x4x4x4x4</Option>
                  <Option value="1">x4x4x8</Option>
                  <Option value="2">x8x4x4</Option>
                  <Option value="3">x8x8</Option>
                  <Option value="4">x16</Option>
                  <Option value="255">Auto</Option>
                </AvailableOptions>
                <DefaultOption>Auto</DefaultOption>
                <Help><![CDATA[Selects PCIe port Bifurcation for selected slot(s)]]></Help>
              </Information>
            </Setting>
            <Menu name="CPU2 PcieBr1D00F0 - Port 1A">
              <Information>
                <Help><![CDATA[Settings related to PCI Express PortS (0/1A/1B/1C/1D/2A/2B/2C/2D/3A/3B/3C/3D/4A/5A) ]]></Help>
              </Information>
              <Subtitle>CPU2 PcieBr1D00F0 - Port 1A</Subtitle>
              <Subtitle>--------------------------------------------------</Subtitle>
              <Subtitle></Subtitle>
              <Setting name="Link Speed" order="7" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Auto</Option>
                    <Option value="1">Gen 1 (2.5 GT/s)</Option>
                    <Option value="2">Gen 2 (5 GT/s)</Option>
                    <Option value="3">Gen 3 (8 GT/s)</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Choose Link Speed for this PCIe port]]></Help>
                </Information>
              </Setting>
              <Text>PCI-E Port Link Status(Link Did Not Train)</Text>
              <Text>PCI-E Port Link Max(Max Width x8)</Text>
              <Text>PCI-E Port Link Speed(Link Did Not Train)</Text>
              <Setting name="PCI-E Port Max Payload Size" order="7" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">128B</Option>
                    <Option value="1">256B</Option>
                    <Option value="2">Auto</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Set Maxpayload size to 256B if possible]]></Help>
                </Information>
              </Setting>
            </Menu>
            <Menu name="CPU2 PcieBr1D02F0 - Port 1C">
              <Information>
                <Help><![CDATA[Settings related to PCI Express PortS (0/1A/1B/1C/1D/2A/2B/2C/2D/3A/3B/3C/3D/4A/5A) ]]></Help>
              </Information>
              <Subtitle>CPU2 PcieBr1D02F0 - Port 1C</Subtitle>
              <Subtitle>--------------------------------------------------</Subtitle>
              <Subtitle></Subtitle>
              <Setting name="Link Speed" order="8" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Auto</Option>
                    <Option value="1">Gen 1 (2.5 GT/s)</Option>
                    <Option value="2">Gen 2 (5 GT/s)</Option>
                    <Option value="3">Gen 3 (8 GT/s)</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Choose Link Speed for this PCIe port]]></Help>
                </Information>
              </Setting>
              <Text>PCI-E Port Link Status(Link Did Not Train)</Text>
              <Text>PCI-E Port Link Max(Max Width x8)</Text>
              <Text>PCI-E Port Link Speed(Link Did Not Train)</Text>
              <Setting name="PCI-E Port Max Payload Size" order="8" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">128B</Option>
                    <Option value="1">256B</Option>
                    <Option value="2">Auto</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Set Maxpayload size to 256B if possible]]></Help>
                </Information>
              </Setting>
            </Menu>
            <Menu name="CPU2 PcieBr2D00F0 - Port 2A">
              <Information>
                <Help><![CDATA[Settings related to PCI Express PortS (0/1A/1B/1C/1D/2A/2B/2C/2D/3A/3B/3C/3D/4A/5A) ]]></Help>
              </Information>
              <Subtitle>CPU2 PcieBr2D00F0 - Port 2A</Subtitle>
              <Subtitle>--------------------------------------------------</Subtitle>
              <Subtitle></Subtitle>
              <Setting name="Link Speed" order="9" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Auto</Option>
                    <Option value="1">Gen 1 (2.5 GT/s)</Option>
                    <Option value="2">Gen 2 (5 GT/s)</Option>
                    <Option value="3">Gen 3 (8 GT/s)</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Choose Link Speed for this PCIe port]]></Help>
                </Information>
              </Setting>
              <Text>PCI-E Port Link Status(Link Did Not Train)</Text>
              <Text>PCI-E Port Link Max(Max Width x16)</Text>
              <Text>PCI-E Port Link Speed(Link Did Not Train)</Text>
              <Setting name="PCI-E Port Max Payload Size" order="9" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">128B</Option>
                    <Option value="1">256B</Option>
                    <Option value="2">Auto</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Set Maxpayload size to 256B if possible]]></Help>
                </Information>
              </Setting>
            </Menu>
            <Menu name="CPU2 PcieBr3D00F0 - Port 3A">
              <Information>
                <Help><![CDATA[Settings related to PCI Express PortS (0/1A/1B/1C/1D/2A/2B/2C/2D/3A/3B/3C/3D/4A/5A) ]]></Help>
              </Information>
              <Subtitle>CPU2 PcieBr3D00F0 - Port 3A</Subtitle>
              <Subtitle>--------------------------------------------------</Subtitle>
              <Subtitle></Subtitle>
              <Setting name="Link Speed" order="10" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Auto</Option>
                    <Option value="1">Gen 1 (2.5 GT/s)</Option>
                    <Option value="2">Gen 2 (5 GT/s)</Option>
                    <Option value="3">Gen 3 (8 GT/s)</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Choose Link Speed for this PCIe port]]></Help>
                </Information>
              </Setting>
              <Text>PCI-E Port Link Status(Linked as x4)</Text>
              <Text>PCI-E Port Link Max(Max Width x4)</Text>
              <Text>PCI-E Port Link Speed(Gen 3 (8.0 GT/s))</Text>
              <Setting name="PCI-E Port Max Payload Size" order="10" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">128B</Option>
                    <Option value="1">256B</Option>
                    <Option value="2">Auto</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Set Maxpayload size to 256B if possible]]></Help>
                </Information>
              </Setting>
            </Menu>
            <Menu name="CPU2 PcieBr3D01F0 - Port 3B">
              <Information>
                <Help><![CDATA[Settings related to PCI Express PortS (0/1A/1B/1C/1D/2A/2B/2C/2D/3A/3B/3C/3D/4A/5A) ]]></Help>
              </Information>
              <Subtitle>CPU2 PcieBr3D01F0 - Port 3B</Subtitle>
              <Subtitle>--------------------------------------------------</Subtitle>
              <Subtitle></Subtitle>
              <Setting name="Link Speed" order="11" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Auto</Option>
                    <Option value="1">Gen 1 (2.5 GT/s)</Option>
                    <Option value="2">Gen 2 (5 GT/s)</Option>
                    <Option value="3">Gen 3 (8 GT/s)</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Choose Link Speed for this PCIe port]]></Help>
                </Information>
              </Setting>
              <Text>PCI-E Port Link Status(Linked as x4)</Text>
              <Text>PCI-E Port Link Max(Max Width x4)</Text>
              <Text>PCI-E Port Link Speed(Gen 3 (8.0 GT/s))</Text>
              <Setting name="PCI-E Port Max Payload Size" order="11" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">128B</Option>
                    <Option value="1">256B</Option>
                    <Option value="2">Auto</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Set Maxpayload size to 256B if possible]]></Help>
                </Information>
              </Setting>
            </Menu>
            <Menu name="CPU2 PcieBr3D02F0 - Port 3C">
              <Information>
                <Help><![CDATA[Settings related to PCI Express PortS (0/1A/1B/1C/1D/2A/2B/2C/2D/3A/3B/3C/3D/4A/5A) ]]></Help>
              </Information>
              <Subtitle>CPU2 PcieBr3D02F0 - Port 3C</Subtitle>
              <Subtitle>--------------------------------------------------</Subtitle>
              <Subtitle></Subtitle>
              <Setting name="Link Speed" order="12" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Auto</Option>
                    <Option value="1">Gen 1 (2.5 GT/s)</Option>
                    <Option value="2">Gen 2 (5 GT/s)</Option>
                    <Option value="3">Gen 3 (8 GT/s)</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Choose Link Speed for this PCIe port]]></Help>
                </Information>
              </Setting>
              <Text>PCI-E Port Link Status(Link Did Not Train)</Text>
              <Text>PCI-E Port Link Max(Max Width x4)</Text>
              <Text>PCI-E Port Link Speed(Link Did Not Train)</Text>
              <Setting name="PCI-E Port Max Payload Size" order="12" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">128B</Option>
                    <Option value="1">256B</Option>
                    <Option value="2">Auto</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Set Maxpayload size to 256B if possible]]></Help>
                </Information>
              </Setting>
            </Menu>
            <Menu name="CPU2 PcieBr3D03F0 - Port 3D">
              <Information>
                <Help><![CDATA[Settings related to PCI Express PortS (0/1A/1B/1C/1D/2A/2B/2C/2D/3A/3B/3C/3D/4A/5A) ]]></Help>
              </Information>
              <Subtitle>CPU2 PcieBr3D03F0 - Port 3D</Subtitle>
              <Subtitle>--------------------------------------------------</Subtitle>
              <Subtitle></Subtitle>
              <Setting name="Link Speed" order="13" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Auto</Option>
                    <Option value="1">Gen 1 (2.5 GT/s)</Option>
                    <Option value="2">Gen 2 (5 GT/s)</Option>
                    <Option value="3">Gen 3 (8 GT/s)</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Choose Link Speed for this PCIe port]]></Help>
                </Information>
              </Setting>
              <Text>PCI-E Port Link Status(Link Did Not Train)</Text>
              <Text>PCI-E Port Link Max(Max Width x4)</Text>
              <Text>PCI-E Port Link Speed(Link Did Not Train)</Text>
              <Setting name="PCI-E Port Max Payload Size" order="13" selectedOption="Auto" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">128B</Option>
                    <Option value="1">256B</Option>
                    <Option value="2">Auto</Option>
                  </AvailableOptions>
                  <DefaultOption>Auto</DefaultOption>
                  <Help><![CDATA[Set Maxpayload size to 256B if possible]]></Help>
                </Information>
              </Setting>
            </Menu>
          </Menu>
          <Menu name="IOAT Configuration">
            <Information>
              <Help><![CDATA[All IOAT configuration options]]></Help>
            </Information>
            <Setting name="Disable TPH" selectedOption="No" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">No</Option>
                  <Option value="1">Yes</Option>
                </AvailableOptions>
                <DefaultOption>No</DefaultOption>
                <Help><![CDATA[TPH is used for data-tagging with a destination ID and a few important attributes. It can send critical data to a particular cache without writing through to memory. Select No in this item for TLP Processing Hint support, which will allow a "TLP request" to provide "hints" to help optimize the processing of each transaction occurred in the target memory space.]]></Help>
              </Information>
            </Setting>
            <Setting name="Prioritize TPH" selectedOption="Disable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="1">Enable</Option>
                  <Option value="0">Disable</Option>
                </AvailableOptions>
                <DefaultOption>Disable</DefaultOption>
                <Help><![CDATA[Select Yes to prioritize the TLP requests that will allow the "hints" to be sent to help facilitate and optimize the processing of certain transactions in the system memory.]]></Help>
                <WorkIf><![CDATA[  1 != Disable TPH  ]]></WorkIf>
              </Information>
            </Setting>
            <Setting name="Relaxed Ordering" selectedOption="Disable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="0">Disable</Option>
                  <Option value="1">Enable</Option>
                </AvailableOptions>
                <DefaultOption>Disable</DefaultOption>
                <Help><![CDATA[Select Enable to enable Relaxed Ordering support which will allow certain transactions to violate the strict-ordering rules of PCI and to be completed prior to other transactions that have already been enqueued.]]></Help>
              </Information>
            </Setting>
          </Menu>
          <Menu name="Intel® VT for Directed I/O (VT-d)" order="1">
            <Information>
              <Help><![CDATA[Press <Enter> to bring up the Intel® VT for Directed I/O (VT-d) Configuration menu.]]></Help>
            </Information>
            <Subtitle>Intel® VT for Directed I/O (VT-d)</Subtitle>
            <Subtitle>--------------------------------------------------</Subtitle>
            <Subtitle></Subtitle>
            <Setting name="Intel® VT for Directed I/O (VT-d)" order="2" selectedOption="Enable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="1">Enable</Option>
                  <Option value="0">Disable</Option>
                </AvailableOptions>
                <DefaultOption>Enable</DefaultOption>
                <Help><![CDATA[Enable/Disable Intel® Virtualization Technology for Directed I/O (VT-d) by reporting the I/O device assignment to VMM through DMAR ACPI Tables.]]></Help>
              </Information>
            </Setting>
            <Setting name="ACS Control" selectedOption="Enable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="1">Enable</Option>
                  <Option value="0">Disable</Option>
                </AvailableOptions>
                <DefaultOption>Enable</DefaultOption>
                <Help><![CDATA[Enable: Programs ACS only to Chipset Pcie Root Ports Bridges; Disable: Programs ACS to all Pcie bridges]]></Help>
                <WorkIf><![CDATA[  0 != Intel® VT for Directed I/O (VT-d)$2  ]]></WorkIf>
              </Information>
            </Setting>
            <Setting name="Interrupt Remapping" selectedOption="Enable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="1">Enable</Option>
                  <Option value="0">Disable</Option>
                </AvailableOptions>
                <DefaultOption>Enable</DefaultOption>
                <Help><![CDATA[Interrupt remapping allows VMM to route device interrupts to the VM that controls the device.]]></Help>
                <WorkIf><![CDATA[  0 != Intel® VT for Directed I/O (VT-d)$2  ]]></WorkIf>
              </Information>
            </Setting>
            <Setting name="PassThrough DMA" selectedOption="Enable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="1">Enable</Option>
                  <Option value="0">Disable</Option>
                </AvailableOptions>
                <DefaultOption>Enable</DefaultOption>
                <Help><![CDATA[Select Enable for the Non-Iscoh VT-d engine to pass through DMA (Direct Memory Access) to enhance system performance.]]></Help>
                <WorkIf><![CDATA[  0 != Intel® VT for Directed I/O (VT-d)$2  ]]></WorkIf>
              </Information>
            </Setting>
            <Setting name="ATS" selectedOption="Enable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="1">Enable</Option>
                  <Option value="0">Disable</Option>
                </AvailableOptions>
                <DefaultOption>Enable</DefaultOption>
                <Help><![CDATA[Select Enable to enable ATS (Address Translation Services) support for the Non-Iscoh VT-d engine to enhance system performance.]]></Help>
                <WorkIf><![CDATA[  0 != Intel® VT for Directed I/O (VT-d)$2  ]]></WorkIf>
              </Information>
            </Setting>
            <Setting name="Posted Interrupt" selectedOption="Enable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="1">Enable</Option>
                  <Option value="0">Disable</Option>
                </AvailableOptions>
                <DefaultOption>Enable</DefaultOption>
                <Help><![CDATA[Select Enable to support VT_D Posted Interrupt which will allow external interrupts to be sent directly from a direct-assigned device to a client machine in non-root mode to improve virtualization efficiency by simplifying interrupt migration and lessening the need of physical interrupts.]]></Help>
                <WorkIf><![CDATA[  0 != Intel® VT for Directed I/O (VT-d)$2  ]]></WorkIf>
              </Information>
            </Setting>
            <Setting name="Coherency Support (Non-Isoch)" selectedOption="Enable" type="Option">
              <Information>
                <AvailableOptions>
                  <Option value="1">Enable</Option>
                  <Option value="0">Disable</Option>
                </AvailableOptions>
                <DefaultOption>Enable</DefaultOption>
                <Help><![CDATA[Select Enable for the Non-Iscoh VT-d engine to pass through DMA (Direct Memory Access) to enhance system performance.]]></Help>
                <WorkIf><![CDATA[  0 != Intel® VT for Directed I/O (VT-d)$2  ]]></WorkIf>
              </Information>
            </Setting>
          </Menu>
          <Menu name="Intel® VMD Technology">
            <Information>
              <Help><![CDATA[Press <Enter> to bring up the Intel® VMD for Volume Management Device Configuration menu.]]></Help>
            </Information>
            <Subtitle>Intel® VMD Technology</Subtitle>
            <Subtitle>--------------------------------------------------</Subtitle>
            <Subtitle></Subtitle>
            <Menu name="Intel® VMD for Volume Management Device on CPU1">
              <Information>
                <Help><![CDATA[]]></Help>
              </Information>
              <Subtitle>VMD Config for PStack1</Subtitle>
              <Subtitle>--------------------------------------------------</Subtitle>
              <Setting name="Intel® VMD for Volume Management Device for PStack1" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology in this Stack.]]></Help>
                </Information>
              </Setting>
              <Setting name="CPU1 SLOT1 VMD port 2A" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology on specific root port]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for PStack1  ]]></WorkIf>
                </Information>
              </Setting>
              <Setting name="CPU1 SLOT1 VMD port 2B" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology on specific root port]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for PStack1  ]]></WorkIf>
                </Information>
              </Setting>
              <Setting name="CPU1 SLOT1 VMD port 2C" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology on specific root port]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for PStack1  ]]></WorkIf>
                </Information>
              </Setting>
              <Setting name="CPU1 SLOT1 VMD port 2D" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology on specific root port]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for PStack1  ]]></WorkIf>
                </Information>
              </Setting>
              <Setting name="Hot Plug Capable" order="1" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Hot Plug for PCIe Root Ports 2A-2D]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for PStack1  ]]></WorkIf>
                </Information>
              </Setting>
              <Subtitle></Subtitle>
              <!--Valid if:   0 != Intel® VMD for Volume Management Device for PStack1  -->
              <Subtitle>VMD Config for PStack2</Subtitle>
              <Subtitle>--------------------------------------------------</Subtitle>
              <Setting name="Intel® VMD for Volume Management Device for PStack2" order="1" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology in this Stack.]]></Help>
                </Information>
              </Setting>
              <Setting name="CPU1 SXB1 M.2 VMD port 3C" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology on specific root port]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for PStack2$1  ]]></WorkIf>
                </Information>
              </Setting>
              <Setting name="CPU1 JF2 M.2 VMD port 3D" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology on specific root port]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for PStack2$1  ]]></WorkIf>
                </Information>
              </Setting>
              <Setting name="Hot Plug Capable" order="2" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Hot Plug for PCIe Root Ports 3A-3D]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for PStack2$1  ]]></WorkIf>
                </Information>
              </Setting>
              <Subtitle></Subtitle>
              <!--Valid if:   0 != Intel® VMD for Volume Management Device for PStack2  -->
            </Menu>
            <Menu name="Intel® VMD for Volume Management Device on CPU2">
              <Information>
                <Help><![CDATA[]]></Help>
              </Information>
              <Subtitle>VMD Config for PStack0</Subtitle>
              <Subtitle>--------------------------------------------------</Subtitle>
              <Setting name="Intel® VMD for Volume Management Device for PStack0" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology in this Stack.]]></Help>
                </Information>
              </Setting>
              <Setting name="CPU2 SXB2 NVMe/SAS VMD port 1A" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology on specific root port]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for PStack0  ]]></WorkIf>
                </Information>
              </Setting>
              <Setting name="CPU2 SXB2 NVMe/SAS VMD port 1B" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology on specific root port]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for PStack0  ]]></WorkIf>
                </Information>
              </Setting>
              <Setting name="Hot Plug Capable" order="3" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Hot Plug for PCIe Root Ports 1A-1D]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for PStack0  ]]></WorkIf>
                </Information>
              </Setting>
              <Subtitle></Subtitle>
              <!--Valid if:   0 != Intel® VMD for Volume Management Device for PStack0  -->
              <Subtitle>VMD Config for PStack1</Subtitle>
              <Subtitle>--------------------------------------------------</Subtitle>
              <Setting name="Intel® VMD for Volume Management Device for Pstack1" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology in this Stack.]]></Help>
                </Information>
              </Setting>
              <Setting name="CPU2 SLOT2 VMD port 2A" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology on specific root port]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for Pstack1  ]]></WorkIf>
                </Information>
              </Setting>
              <Setting name="CPU2 SLOT2 VMD port 2B" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology on specific root port]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for Pstack1  ]]></WorkIf>
                </Information>
              </Setting>
              <Setting name="CPU2 SLOT2 VMD port 2C" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology on specific root port]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for Pstack1  ]]></WorkIf>
                </Information>
              </Setting>
              <Setting name="CPU2 SLOT2 VMD port 2D" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology on specific root port]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for Pstack1  ]]></WorkIf>
                </Information>
              </Setting>
              <Setting name="Hot Plug Capable" order="4" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Hot Plug for PCIe Root Ports 2A-2D]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for Pstack1  ]]></WorkIf>
                </Information>
              </Setting>
              <Subtitle></Subtitle>
              <!--Valid if:   0 != Intel® VMD for Volume Management Device for Pstack1  -->
              <Subtitle>VMD Config for PStack2</Subtitle>
              <Subtitle>--------------------------------------------------</Subtitle>
              <Setting name="Intel® VMD for Volume Management Device for PStack2" order="2" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology in this Stack.]]></Help>
                </Information>
              </Setting>
              <Setting name="CPU2 SXB2 NVMe VMD port 3A" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology on specific root port]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for PStack2$2  ]]></WorkIf>
                </Information>
              </Setting>
              <Setting name="CPU2 SXB2 NVMe VMD port 3B" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology on specific root port]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for PStack2$2  ]]></WorkIf>
                </Information>
              </Setting>
              <Setting name="CPU2 SXB2 NVMe VMD port 3C" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology on specific root port]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for PStack2$2  ]]></WorkIf>
                </Information>
              </Setting>
              <Setting name="CPU2 SXB2 NVMe VMD port 3D" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Intel® Volume Management Device Technology on specific root port]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for PStack2$2  ]]></WorkIf>
                </Information>
              </Setting>
              <Setting name="Hot Plug Capable" order="5" selectedOption="Disable" type="Option">
                <Information>
                  <AvailableOptions>
                    <Option value="0">Disable</Option>
                    <Option value="1">Enable</Option>
                  </AvailableOptions>
                  <DefaultOption>Disable</DefaultOption>
                  <Help><![CDATA[Enable/Disable Hot Plug for PCIe Root Ports 3A-3D]]></Help>
                  <WorkIf><![CDATA[  0 != Intel® VMD for Volume Management Device for PStack2$2  ]]></WorkIf>
                </Information>
              </Setting>
              <Subtitle></Subtitle>
              <!--Valid if:   0 != Intel® VMD for Volume Management Device for PStack2  -->
            </Menu>
          </Menu>
          <Subtitle></Subtitle>
          <Subtitle> IIO-PCIE Express Global Options</Subtitle>
          <Subtitle>========================================</Subtitle>
          <Setting name="PCI-E Completion Timeout Disable" selectedOption="No" type="Option">
            <Information>
              <AvailableOptions>
                <Option value="1">Yes</Option>
                <Option value="0">No</Option>
                <Option value="2">Per-Port</Option>
              </AvailableOptions>
              <DefaultOption>No</DefaultOption>
              <Help><![CDATA[Select Enable to enable PCI-E Completion Timeout support for electric tuning.]]></Help>
            </Information>
          </Setting>
        </Menu>
      </Menu>
      <Menu name="South Bridge">
        <Information>
          <Help><![CDATA[South Bridge Parameters]]></Help>
        </Information>
        <Subtitle></Subtitle>
        <Text>USB Module Version(21)</Text>
        <Subtitle></Subtitle>
        <Text>USB Devices:()</Text>
        <Subtitle>      1 Keyboard, 1 Mouse, 1 Hub</Subtitle>
        <Subtitle></Subtitle>
        <Setting name="Legacy USB Support" selectedOption="Enabled" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">Enabled</Option>
              <Option value="1">Disabled</Option>
              <Option value="2">Auto</Option>
            </AvailableOptions>
            <DefaultOption>Enabled</DefaultOption>
            <Help><![CDATA[Enables Legacy USB support. AUTO option disables legacy support if no USB devices are connected. DISABLE option will keep USB devices available only for EFI applications.]]></Help>
          </Information>
        </Setting>
        <Setting name="XHCI Hand-off" selectedOption="Enabled" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="1">Enabled</Option>
              <Option value="0">Disabled</Option>
            </AvailableOptions>
            <DefaultOption>Enabled</DefaultOption>
            <Help><![CDATA[This is a workaround for OSes without XHCI hand-off support. The XHCI ownership change should be claimed by XHCI driver.]]></Help>
          </Information>
        </Setting>
        <Setting name="Port 60/64 Emulation" selectedOption="Enabled" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">Disabled</Option>
              <Option value="1">Enabled</Option>
            </AvailableOptions>
            <DefaultOption>Enabled</DefaultOption>
            <Help><![CDATA[Enables I/O port 60h/64h emulation support. This should be enabled for the complete USB keyboard legacy support for non-USB aware OSes.]]></Help>
          </Information>
        </Setting>
        <Setting name="PCIe PLL SSC" selectedOption="Disable" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="255">Disable</Option>
              <Option value="5">Enable</Option>
            </AvailableOptions>
            <DefaultOption>Disable</DefaultOption>
            <Help><![CDATA[Enable or Disable PCIe PLL spread specturm clocking.]]></Help>
          </Information>
        </Setting>
      </Menu>
    </Menu>
    <Menu name="Server ME Information">
      <Information>
        <Help><![CDATA[Configure Server ME Technology Parameters]]></Help>
      </Information>
      <Subtitle>General ME Configuration</Subtitle>
      <Text>Oper. Firmware Version(4.1.3.239)</Text>
      <Text>Backup Firmware Version(N/A)</Text>
      <Text>Recovery Firmware Version(4.1.3.239)</Text>
      <Text>ME Firmware Status #1(0x000F0255)</Text>
      <Text>ME Firmware Status #2(0x88114026)</Text>
      <Text>  Current State(Operational)</Text>
      <Text>  Error Code(No Error)</Text>
    </Menu>
    <Menu name="PCH SATA Configuration">
      <Information>
        <Help><![CDATA[SATA devices and settings]]></Help>
      </Information>
      <Subtitle>PCH SATA Configuration</Subtitle>
      <Subtitle>--------------------------------------------------</Subtitle>
      <Subtitle></Subtitle>
      <Setting name="SATA Controller" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Enable or Disable SATA Controller]]></Help>
        </Information>
      </Setting>
      <Setting name="Configure SATA as" selectedOption="AHCI" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">AHCI</Option>
            <Option value="1">RAID</Option>
          </AvailableOptions>
          <DefaultOption>AHCI</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
          <WorkIf><![CDATA[  0 != SATA Controller  ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="SATA HDD Unlock" order="1" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Enable/Disable unlock HDD password ability in the OS while SATA is configured as RAID mode]]></Help>
          <WorkIf><![CDATA[  0 != SATA Controller  ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="SATA RSTe Boot Info" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Enable setting provides full int13h support for SATA controller attached devices. CSM storage OPROM policy should be set to legacy to make this selection effective.]]></Help>
          <WorkIf><![CDATA[ ( 0 != SATA Controller )  and  ( 1 == Configure SATA as ) ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="Aggressive Link Power Management" order="1" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[Enables/Disables SALP]]></Help>
        </Information>
      </Setting>
      <Setting name="SATA RAID Option ROM/UEFI Driver" selectedOption="Legacy" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">EFI</Option>
            <Option value="2">Legacy</Option>
          </AvailableOptions>
          <DefaultOption>Legacy</DefaultOption>
          <Help><![CDATA[In RAID mode load EFI driver. (If disabled loads LEGACY OPROM)]]></Help>
          <WorkIf><![CDATA[  1 == Configure SATA as  ]]></WorkIf>
        </Information>
      </Setting>
      <Subtitle></Subtitle>
      <Text>SATA Port 0([Not Installed])</Text>
      <Setting name="  Hot Plug" order="1" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="1" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  SATA Device Type" order="1" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>SATA Port 1([Not Installed])</Text>
      <Setting name="  Hot Plug" order="2" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="2" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  SATA Device Type" order="2" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>SATA Port 2([Not Installed])</Text>
      <Setting name="  Hot Plug" order="3" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="3" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  SATA Device Type" order="3" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>SATA Port 3([Not Installed])</Text>
      <Setting name="  Hot Plug" order="4" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="4" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  SATA Device Type" order="4" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>SATA Port 4([Not Installed])</Text>
      <Setting name="  Hot Plug" order="5" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="5" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  SATA Device Type" order="5" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>SATA Port 5([Not Installed])</Text>
      <Setting name="  Hot Plug" order="6" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="6" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  SATA Device Type" order="6" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>SATA Port 6(SuperMicro SSD - 63.3 GB)</Text>
      <Setting name="  Hot Plug" order="7" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="7" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  SATA Device Type" order="7" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>SATA Port 7([Not Installed])</Text>
      <Setting name="  Hot Plug" order="8" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="8" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  SATA Device Type" order="8" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
    </Menu>
    <Menu name="PCH sSATA Configuration">
      <Information>
        <Help><![CDATA[sSATA devices and settings]]></Help>
      </Information>
      <Subtitle>PCH sSATA Configuration</Subtitle>
      <Subtitle>--------------------------------------------------</Subtitle>
      <Subtitle></Subtitle>
      <Setting name="sSATA Controller" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Enable</Option>
            <Option value="0">Disable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Enable or Disable SATA Controller]]></Help>
        </Information>
      </Setting>
      <Setting name="Configure sSATA as" selectedOption="AHCI" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">AHCI</Option>
            <Option value="1">RAID</Option>
          </AvailableOptions>
          <DefaultOption>AHCI</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
          <WorkIf><![CDATA[  0 != sSATA Controller  ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="SATA HDD Unlock" order="2" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Enable/Disable unlock HDD password ability in the OS while SATA is configured as RAID mode]]></Help>
          <WorkIf><![CDATA[  0 != sSATA Controller  ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="sSATA RSTe Boot Info" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Enable setting provides full int13h support for SATA controller attached devices. CSM storage OPROM policy should be set to legacy to make this selection effective.]]></Help>
          <WorkIf><![CDATA[ ( 0 != sSATA Controller )  and  ( 1 == Configure sSATA as ) ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="Aggressive Link Power Management" order="2" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[Enables/Disables SALP]]></Help>
        </Information>
      </Setting>
      <Setting name="sSATA RAID Option ROM/UEFI Driver" selectedOption="Legacy" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">EFI</Option>
            <Option value="2">Legacy</Option>
          </AvailableOptions>
          <DefaultOption>Legacy</DefaultOption>
          <Help><![CDATA[In RAID mode load EFI driver. (If disabled loads LEGACY OPROM)]]></Help>
          <WorkIf><![CDATA[  1 == Configure sSATA as  ]]></WorkIf>
        </Information>
      </Setting>
      <Subtitle></Subtitle>
      <Text>sSATA Port 0([Not Installed])</Text>
      <Setting name="  Hot Plug" order="9" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="9" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  sSATA Device Type" order="1" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>sSATA Port 1([Not Installed])</Text>
      <Setting name="  Hot Plug" order="10" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="10" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  sSATA Device Type" order="2" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>sSATA Port 2([Not Installed])</Text>
      <Setting name="  Hot Plug" order="11" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="11" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  sSATA Device Type" order="3" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>sSATA Port 3([Not Installed])</Text>
      <Setting name="  Hot Plug" order="12" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="12" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  sSATA Device Type" order="4" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>sSATA Port 4([Not Installed])</Text>
      <Setting name="  Hot Plug" order="13" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="13" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  sSATA Device Type" order="5" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
      <Text>sSATA Port 5([Not Installed])</Text>
      <Setting name="  Hot Plug" order="14" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Designates this port as Hot Pluggable.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Spin Up Device" order="14" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[If enabled for any of ports Staggerred Spin Up will be performed and only the drives witch have this option enabled will spin up at boot. Otherwise all drives spin up at boot.]]></Help>
        </Information>
      </Setting>
      <Setting name="  sSATA Device Type" order="6" selectedOption="Hard Disk Drive" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Hard Disk Drive</Option>
            <Option value="1">Solid State Drive</Option>
          </AvailableOptions>
          <DefaultOption>Hard Disk Drive</DefaultOption>
          <Help><![CDATA[Identify the SATA port is connected to Solid State Drive or Hard Disk Drive]]></Help>
        </Information>
      </Setting>
    </Menu>
    <Menu name="PCIe/PCI/PnP Configuration">
      <Information>
        <Help><![CDATA[PCI, PCI-X and PCI Express Settings.]]></Help>
      </Information>
      <Text>PCI Bus Driver Version(A5.01.18)</Text>
      <Subtitle></Subtitle>
      <Subtitle>PCI Devices Common Settings:</Subtitle>
      <Setting name="Above 4G Decoding" selectedOption="Enabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Enabled</DefaultOption>
          <Help><![CDATA[Enables or Disables 64bit capable Devices to be Decoded in Above 4G Address Space (Only if System Supports 64 bit PCI Decoding).]]></Help>
        </Information>
      </Setting>
      <Setting name="SR-IOV Support" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[If system has SR-IOV capable PCIe Devices, this option Enables or Disables Single Root IO Virtualization Support.]]></Help>
        </Information>
      </Setting>
      <Setting name="MMIO High Base" selectedOption="56T" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">56T</Option>
            <Option value="1">40T</Option>
            <Option value="2">24T</Option>
            <Option value="3">16T</Option>
            <Option value="4">4T</Option>
            <Option value="6">2T</Option>
            <Option value="5">1T</Option>
          </AvailableOptions>
          <DefaultOption>56T</DefaultOption>
          <Help><![CDATA[Select MMIO High Base]]></Help>
        </Information>
      </Setting>
      <Setting name="MMIO High Granularity Size" selectedOption="256G" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">1G</Option>
            <Option value="1">4G</Option>
            <Option value="2">16G</Option>
            <Option value="3">64G</Option>
            <Option value="4">256G</Option>
            <Option value="5">1024G</Option>
          </AvailableOptions>
          <DefaultOption>256G</DefaultOption>
          <Help><![CDATA[Selects the allocation size used to assign mmioh resources.
Total mmioh space can be up to 32xgranularity.
Per stack mmioh resource assignments are multiples of the granularity where 1 unit per stack is the default allocation.]]></Help>
        </Information>
      </Setting>
      <Setting name="Maximum Read Request" selectedOption="Auto" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="55">Auto</Option>
            <Option value="0">128 Bytes</Option>
            <Option value="1">256 Bytes</Option>
            <Option value="2">512 Bytes</Option>
            <Option value="3">1024 Bytes</Option>
            <Option value="4">2048 Bytes</Option>
            <Option value="5">4096 Bytes</Option>
          </AvailableOptions>
          <DefaultOption>Auto</DefaultOption>
          <Help><![CDATA[Set Maximum Read Request Size of PCI Express Device or allow System BIOS to select the value.]]></Help>
        </Information>
      </Setting>
      <Setting name="MMCFG Base" selectedOption="2G" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">1G</Option>
            <Option value="1">1.5G</Option>
            <Option value="2">1.75G</Option>
            <Option value="3">2G</Option>
            <Option value="4">2.25G</Option>
            <Option value="5">3G</Option>
          </AvailableOptions>
          <DefaultOption>2G</DefaultOption>
          <Help><![CDATA[Select MMCFG Base]]></Help>
        </Information>
      </Setting>
      <Setting name="NVMe Firmware Source" selectedOption="Vendor Defined Firmware" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Vendor Defined Firmware</Option>
            <Option value="1">AMI Native Support</Option>
          </AvailableOptions>
          <DefaultOption>Vendor Defined Firmware</DefaultOption>
          <Help><![CDATA[AMI Native FW Support or Device Vendor Defined FW Support]]></Help>
        </Information>
      </Setting>
      <Setting name="VGA Priority" selectedOption="Onboard" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Onboard</Option>
            <Option value="2">Offboard</Option>
          </AvailableOptions>
          <DefaultOption>Onboard</DefaultOption>
          <Help><![CDATA[Select active Video type]]></Help>
        </Information>
      </Setting>
      <Setting name="CPU1 RSC-R1UTP-E16R PCI-E 3.0 X16 OPROM" selectedOption="EFI" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Legacy</Option>
            <Option value="2">EFI</Option>
          </AvailableOptions>
          <DefaultOption>Legacy</DefaultOption>
          <Help><![CDATA[Enables or disables CPU1 RSC-R1UTP-E16R PCI-E 3.0 X16 OPROM option.]]></Help>
        </Information>
      </Setting>
      <Setting name="CPU2 RSC-P-6 PCI-E 3.0 X16 OPROM" selectedOption="EFI" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Legacy</Option>
            <Option value="2">EFI</Option>
          </AvailableOptions>
          <DefaultOption>Legacy</DefaultOption>
          <Help><![CDATA[Enables or disables CPU2 RSC-P-6 PCI-E 3.0 X16 OPROM option.]]></Help>
        </Information>
      </Setting>
      <Setting name="CPU1 SXB1 PCI-E 3.0 X4 OPROM" selectedOption="EFI" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Legacy</Option>
            <Option value="2">EFI</Option>
          </AvailableOptions>
          <DefaultOption>Legacy</DefaultOption>
          <Help><![CDATA[Enables or disables CPU1 PCI-E 3.0 X4 OPROM option.]]></Help>
        </Information>
      </Setting>
      <Setting name="CPU1 JF2 PCI-E 3.0 X4 OPROM" selectedOption="EFI" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Legacy</Option>
            <Option value="2">EFI</Option>
          </AvailableOptions>
          <DefaultOption>Legacy</DefaultOption>
          <Help><![CDATA[Enables or disables CPU1 JF2 PCI-E 3.0 X4 OPROM option.]]></Help>
        </Information>
      </Setting>
      <Setting name="SIOM CPU1 PCI-E 3.0 X16 OPROM" selectedOption="EFI" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Legacy</Option>
            <Option value="2">EFI</Option>
          </AvailableOptions>
          <DefaultOption>Legacy</DefaultOption>
          <Help><![CDATA[Enables or disables SIOM CPU1 PCI-E 3.0 X16 OPROM option.]]></Help>
        </Information>
      </Setting>
      <Setting name="Bus Master Enable" selectedOption="Enabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Enabled</Option>
            <Option value="0">Disabled</Option>
          </AvailableOptions>
          <DefaultOption>Enabled</DefaultOption>
          <Help><![CDATA[Enabled - PCI Bus Driver enables the Bus Master Attribute for DMA transactions.
Disabled - PCI Bus Driver Disables the Bus Master Attribute for Pre-Boot DMA Protection.]]></Help>
        </Information>
      </Setting>
      <Setting name="Onboard NVMe1 Option ROM" selectedOption="EFI" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <!--Option ValidIf:  ( 1 == 0 ) -->
            <Option value="1">Legacy</Option>
            <Option value="2">EFI</Option>
          </AvailableOptions>
          <DefaultOption>EFI</DefaultOption>
          <Help><![CDATA[Select which firmware function to be loaded for onboard NVMe1.]]></Help>
        </Information>
      </Setting>
      <Setting name="Onboard NVMe2 Option ROM" selectedOption="EFI" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <!--Option ValidIf:  ( 1 == 0 ) -->
            <Option value="1">Legacy</Option>
            <Option value="2">EFI</Option>
          </AvailableOptions>
          <DefaultOption>EFI</DefaultOption>
          <Help><![CDATA[Select which firmware function to be loaded for onboard NVMe2]]></Help>
        </Information>
      </Setting>
      <Setting name="Onboard NVMe3 Option ROM" selectedOption="EFI" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <!--Option ValidIf:  ( 1 == 0 ) -->
            <Option value="1">Legacy</Option>
            <Option value="2">EFI</Option>
          </AvailableOptions>
          <DefaultOption>EFI</DefaultOption>
          <Help><![CDATA[Select which firmware function to be loaded for onboard NVMe3]]></Help>
        </Information>
      </Setting>
      <Setting name="Onboard NVMe4 Option ROM" selectedOption="EFI" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <!--Option ValidIf:  ( 1 == 0 ) -->
            <Option value="1">Legacy</Option>
            <Option value="2">EFI</Option>
          </AvailableOptions>
          <DefaultOption>EFI</DefaultOption>
          <Help><![CDATA[Select which firmware function to be loaded for onboard NVMe4]]></Help>
        </Information>
      </Setting>
      <Setting name="Onboard Video Option ROM" selectedOption="EFI" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Legacy</Option>
            <Option value="2">EFI</Option>
          </AvailableOptions>
          <DefaultOption>Legacy</DefaultOption>
          <Help><![CDATA[Select which onboard video firmware type to be loaded.]]></Help>
        </Information>
      </Setting>
      <Subtitle></Subtitle>
      <Subtitle></Subtitle>
      <Subtitle></Subtitle>
    </Menu>
    <Menu name="Super IO Configuration">
      <Information>
        <Help><![CDATA[System Super IO Chip Parameters.]]></Help>
      </Information>
      <Subtitle>Super IO Configuration</Subtitle>
      <Subtitle></Subtitle>
      <Text>Super IO Chip(AST2500)</Text>
      <Menu name="Serial Port 1 Configuration">
        <Information>
          <Help><![CDATA[Set Parameters of Serial Port 1 (COMA)]]></Help>
        </Information>
        <Subtitle>Serial Port 1 Configuration</Subtitle>
        <Subtitle></Subtitle>
        <Setting name="Serial Port 1" checkedStatus="Checked" type="CheckBox">
          <!--Checked/Unchecked-->
          <Information>
            <DefaultStatus>Checked</DefaultStatus>
            <Help><![CDATA[Enable or Disable Serial Port (COM)]]></Help>
          </Information>
        </Setting>
        <Text>Device Settings(IO=3F8h; IRQ=4;)</Text>
        <!--Valid if:   0 != Serial Port 1  -->
        <Subtitle></Subtitle>
        <Setting name="Change Settings" order="1" selectedOption="Auto" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">Auto</Option>
              <Option value="1">IO=3F8h; IRQ=4;</Option>
              <Option value="3">IO=2F8h; IRQ=4;</Option>
              <Option value="4">IO=3E8h; IRQ=4;</Option>
              <Option value="5">IO=2E8h; IRQ=4;</Option>
            </AvailableOptions>
            <DefaultOption>Auto</DefaultOption>
            <Help><![CDATA[Select an optimal settings for Super IO Device]]></Help>
            <WorkIf><![CDATA[  0 != Serial Port 1  ]]></WorkIf>
          </Information>
        </Setting>
      </Menu>
      <Menu name="Serial Port 2 Configuration">
        <Information>
          <Help><![CDATA[Set Parameters of Serial Port 2 (COMB)]]></Help>
        </Information>
        <Subtitle>Serial Port 2 Configuration</Subtitle>
        <Subtitle></Subtitle>
        <Setting name="Serial Port 2" checkedStatus="Checked" type="CheckBox">
          <!--Checked/Unchecked-->
          <Information>
            <DefaultStatus>Checked</DefaultStatus>
            <Help><![CDATA[Enable or Disable Serial Port (COM)]]></Help>
          </Information>
        </Setting>
        <Text>Device Settings(IO=2F8h; IRQ=3;)</Text>
        <!--Valid if:   0 != Serial Port 2  -->
        <Subtitle></Subtitle>
        <Setting name="Change Settings" order="2" selectedOption="Auto" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">Auto</Option>
              <Option value="1">IO=2F8h; IRQ=3;</Option>
              <Option value="2">IO=3F8h; IRQ=3;</Option>
              <Option value="4">IO=3E8h; IRQ=3;</Option>
              <Option value="5">IO=2E8h; IRQ=3;</Option>
            </AvailableOptions>
            <DefaultOption>Auto</DefaultOption>
            <Help><![CDATA[Select an optimal settings for Super IO Device]]></Help>
            <WorkIf><![CDATA[  0 != Serial Port 2  ]]></WorkIf>
          </Information>
        </Setting>
        <Setting name="Serial Port 2 Attribute" selectedOption="SOL" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="3">SOL</Option>
              <Option value="1">COM</Option>
            </AvailableOptions>
            <DefaultOption>SOL</DefaultOption>
            <Help><![CDATA[Select serial port 2 mode]]></Help>
          </Information>
        </Setting>
      </Menu>
    </Menu>
    <Menu name="Serial Port Console Redirection">
      <Information>
        <Help><![CDATA[Serial Port Console Redirection]]></Help>
      </Information>
      <Subtitle></Subtitle>
      <Subtitle>COM1</Subtitle>
      <Setting name="Console Redirection" order="1" checkedStatus="Unchecked" type="CheckBox">
        <!--Checked/Unchecked-->
        <Information>
          <DefaultStatus>Unchecked</DefaultStatus>
          <Help><![CDATA[Console Redirection Enable or Disable.]]></Help>
        </Information>
      </Setting>
      <Menu name="Console Redirection Settings" order="1">
        <Information>
          <Help><![CDATA[The settings specify how the host computer and the remote computer (which the user is using) will exchange data. Both computers should have the same or compatible settings.]]></Help>
          <WorkIf><![CDATA[ ( 0 != Console Redirection$1 ) ]]></WorkIf>
        </Information>
        <Subtitle>COM1</Subtitle>
        <Subtitle>Console Redirection Settings</Subtitle>
        <Subtitle></Subtitle>
        <Setting name="Terminal Type" order="1" selectedOption="VT100+" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">VT100</Option>
              <Option value="1">VT100+</Option>
              <Option value="2">VT-UTF8</Option>
              <Option value="3">ANSI</Option>
            </AvailableOptions>
            <DefaultOption>VT100+</DefaultOption>
            <Help><![CDATA[Emulation: ANSI: Extended ASCII char set. VT100: ASCII char set. VT100+: Extends VT100 to support color, function keys, etc. VT-UTF8: Uses UTF8 encoding to map Unicode chars onto 1 or more bytes.]]></Help>
          </Information>
        </Setting>
        <Setting name="Bits per second" order="1" selectedOption="115200" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="3">9600</Option>
              <Option value="4">19200</Option>
              <Option value="5">38400</Option>
              <Option value="6">57600</Option>
              <Option value="7">115200</Option>
            </AvailableOptions>
            <DefaultOption>115200</DefaultOption>
            <Help><![CDATA[Selects serial port transmission speed. The speed must be matched on the other side. Long or noisy lines may require lower speeds.]]></Help>
          </Information>
        </Setting>
        <Setting name="Data Bits" order="1" selectedOption="8" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="7">7</Option>
              <Option value="8">8</Option>
            </AvailableOptions>
            <DefaultOption>8</DefaultOption>
            <Help><![CDATA[Data Bits]]></Help>
          </Information>
        </Setting>
        <Setting name="Parity" order="1" selectedOption="None" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="1">None</Option>
              <Option value="2">Even</Option>
              <Option value="3">Odd</Option>
              <Option value="4">Mark</Option>
              <Option value="5">Space</Option>
            </AvailableOptions>
            <DefaultOption>None</DefaultOption>
            <Help><![CDATA[A parity bit can be sent with the data bits to detect some transmission errors. Even: parity bit is 0 if the num of 1's in the data bits is even. Odd: parity bit is 0 if num of 1's in the data bits is odd.  Mark: parity bit is always 1. Space: Parity bit is always 0. Mark and Space Parity do not allow for error detection. They can be used as an additional data bit.]]></Help>
          </Information>
        </Setting>
        <Setting name="Stop Bits" order="1" selectedOption="1" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="1">1</Option>
              <Option value="3">2</Option>
            </AvailableOptions>
            <DefaultOption>1</DefaultOption>
            <Help><![CDATA[Stop bits indicate the end of a serial data packet. (A start bit indicates the beginning). The standard setting is 1 stop bit. Communication with slow devices may require more than 1 stop bit.]]></Help>
          </Information>
        </Setting>
        <Setting name="Flow Control" order="1" selectedOption="None" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">None</Option>
              <Option value="1">Hardware RTS/CTS</Option>
            </AvailableOptions>
            <DefaultOption>None</DefaultOption>
            <Help><![CDATA[Flow control can prevent data loss from buffer overflow. When sending data, if the receiving buffers are full, a 'stop' signal can be sent to stop the data flow. Once the buffers are empty, a 'start' signal can be sent to re-start the flow. Hardware flow control uses two wires to send start/stop signals.]]></Help>
          </Information>
        </Setting>
        <Setting name="VT-UTF8 Combo Key Support" order="1" checkedStatus="Checked" type="CheckBox">
          <!--Checked/Unchecked-->
          <Information>
            <DefaultStatus>Checked</DefaultStatus>
            <Help><![CDATA[Enable VT-UTF8 Combination Key Support for ANSI/VT100 terminals]]></Help>
          </Information>
        </Setting>
        <Setting name="Recorder Mode" order="1" checkedStatus="Unchecked" type="CheckBox">
          <!--Checked/Unchecked-->
          <Information>
            <DefaultStatus>Unchecked</DefaultStatus>
            <Help><![CDATA[With this mode enabled only text will be sent. This is to capture Terminal data.]]></Help>
          </Information>
        </Setting>
        <Setting name="Resolution 100x31" order="1" checkedStatus="Checked" type="CheckBox">
          <!--Checked/Unchecked-->
          <Information>
            <DefaultStatus>Checked</DefaultStatus>
            <Help><![CDATA[Enables or disables extended terminal resolution]]></Help>
          </Information>
        </Setting>
        <Setting name="Legacy OS Redirection Resolution" order="1" selectedOption="80x25" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">80x24</Option>
              <Option value="1">80x25</Option>
            </AvailableOptions>
            <DefaultOption>80x25</DefaultOption>
            <Help><![CDATA[On Legacy OS, the Number of Rows and Columns supported redirection]]></Help>
          </Information>
        </Setting>
        <Setting name="Putty KeyPad" order="1" selectedOption="VT100" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="1">VT100</Option>
              <Option value="2">LINUX</Option>
              <Option value="4">XTERMR6</Option>
              <Option value="8">SCO</Option>
              <Option value="16">ESCN</Option>
              <Option value="32">VT400</Option>
            </AvailableOptions>
            <DefaultOption>VT100</DefaultOption>
            <Help><![CDATA[Select FunctionKey and KeyPad on Putty.]]></Help>
          </Information>
        </Setting>
        <Setting name="Redirection After BIOS POST" order="1" selectedOption="Always Enable" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">Always Enable</Option>
              <Option value="1">BootLoader</Option>
            </AvailableOptions>
            <DefaultOption>Always Enable</DefaultOption>
            <Help><![CDATA[When Bootloader is selected, then Legacy Console Redirection is disabled before booting to legacy OS. When Always Enable is selected, then Legacy Console Redirection is enabled for legacy OS. Default setting for this option is set to Always Enable.]]></Help>
          </Information>
        </Setting>
        <Subtitle></Subtitle>
      </Menu>
      <Subtitle></Subtitle>
      <Subtitle>SOL/COM2</Subtitle>
      <Setting name="Console Redirection" order="2" checkedStatus="Checked" type="CheckBox">
        <!--Checked/Unchecked-->
        <Information>
          <DefaultStatus>Checked</DefaultStatus>
          <Help><![CDATA[Console Redirection Enable or Disable.]]></Help>
        </Information>
      </Setting>
      <Menu name="Console Redirection Settings" order="2">
        <Information>
          <Help><![CDATA[The settings specify how the host computer and the remote computer (which the user is using) will exchange data. Both computers should have the same or compatible settings.]]></Help>
          <WorkIf><![CDATA[ ( 0 != Console Redirection$2 ) ]]></WorkIf>
        </Information>
        <Subtitle>SOL/COM2</Subtitle>
        <Subtitle>Console Redirection Settings</Subtitle>
        <Subtitle></Subtitle>
        <Setting name="Terminal Type" order="2" selectedOption="VT100+" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">VT100</Option>
              <Option value="1">VT100+</Option>
              <Option value="2">VT-UTF8</Option>
              <Option value="3">ANSI</Option>
            </AvailableOptions>
            <DefaultOption>VT100+</DefaultOption>
            <Help><![CDATA[Emulation: ANSI: Extended ASCII char set. VT100: ASCII char set. VT100+: Extends VT100 to support color, function keys, etc. VT-UTF8: Uses UTF8 encoding to map Unicode chars onto 1 or more bytes.]]></Help>
          </Information>
        </Setting>
        <Setting name="Bits per second" order="2" selectedOption="115200" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="3">9600</Option>
              <Option value="4">19200</Option>
              <Option value="5">38400</Option>
              <Option value="6">57600</Option>
              <Option value="7">115200</Option>
            </AvailableOptions>
            <DefaultOption>115200</DefaultOption>
            <Help><![CDATA[Selects serial port transmission speed. The speed must be matched on the other side. Long or noisy lines may require lower speeds.]]></Help>
          </Information>
        </Setting>
        <Setting name="Data Bits" order="2" selectedOption="8" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="7">7</Option>
              <Option value="8">8</Option>
            </AvailableOptions>
            <DefaultOption>8</DefaultOption>
            <Help><![CDATA[Data Bits]]></Help>
          </Information>
        </Setting>
        <Setting name="Parity" order="2" selectedOption="None" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="1">None</Option>
              <Option value="2">Even</Option>
              <Option value="3">Odd</Option>
              <Option value="4">Mark</Option>
              <Option value="5">Space</Option>
            </AvailableOptions>
            <DefaultOption>None</DefaultOption>
            <Help><![CDATA[A parity bit can be sent with the data bits to detect some transmission errors. Even: parity bit is 0 if the num of 1's in the data bits is even. Odd: parity bit is 0 if num of 1's in the data bits is odd.  Mark: parity bit is always 1. Space: Parity bit is always 0. Mark and Space Parity do not allow for error detection. They can be used as an additional data bit.]]></Help>
          </Information>
        </Setting>
        <Setting name="Stop Bits" order="2" selectedOption="1" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="1">1</Option>
              <Option value="3">2</Option>
            </AvailableOptions>
            <DefaultOption>1</DefaultOption>
            <Help><![CDATA[Stop bits indicate the end of a serial data packet. (A start bit indicates the beginning). The standard setting is 1 stop bit. Communication with slow devices may require more than 1 stop bit.]]></Help>
          </Information>
        </Setting>
        <Setting name="Flow Control" order="2" selectedOption="None" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">None</Option>
              <Option value="1">Hardware RTS/CTS</Option>
            </AvailableOptions>
            <DefaultOption>None</DefaultOption>
            <Help><![CDATA[Flow control can prevent data loss from buffer overflow. When sending data, if the receiving buffers are full, a 'stop' signal can be sent to stop the data flow. Once the buffers are empty, a 'start' signal can be sent to re-start the flow. Hardware flow control uses two wires to send start/stop signals.]]></Help>
          </Information>
        </Setting>
        <Setting name="VT-UTF8 Combo Key Support" order="2" checkedStatus="Checked" type="CheckBox">
          <!--Checked/Unchecked-->
          <Information>
            <DefaultStatus>Checked</DefaultStatus>
            <Help><![CDATA[Enable VT-UTF8 Combination Key Support for ANSI/VT100 terminals]]></Help>
          </Information>
        </Setting>
        <Setting name="Recorder Mode" order="2" checkedStatus="Unchecked" type="CheckBox">
          <!--Checked/Unchecked-->
          <Information>
            <DefaultStatus>Unchecked</DefaultStatus>
            <Help><![CDATA[With this mode enabled only text will be sent. This is to capture Terminal data.]]></Help>
          </Information>
        </Setting>
        <Setting name="Resolution 100x31" order="2" checkedStatus="Checked" type="CheckBox">
          <!--Checked/Unchecked-->
          <Information>
            <DefaultStatus>Checked</DefaultStatus>
            <Help><![CDATA[Enables or disables extended terminal resolution]]></Help>
          </Information>
        </Setting>
        <Setting name="Legacy OS Redirection Resolution" order="2" selectedOption="80x25" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">80x24</Option>
              <Option value="1">80x25</Option>
            </AvailableOptions>
            <DefaultOption>80x25</DefaultOption>
            <Help><![CDATA[On Legacy OS, the Number of Rows and Columns supported redirection]]></Help>
          </Information>
        </Setting>
        <Setting name="Putty KeyPad" order="2" selectedOption="VT100" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="1">VT100</Option>
              <Option value="2">LINUX</Option>
              <Option value="4">XTERMR6</Option>
              <Option value="8">SCO</Option>
              <Option value="16">ESCN</Option>
              <Option value="32">VT400</Option>
            </AvailableOptions>
            <DefaultOption>VT100</DefaultOption>
            <Help><![CDATA[Select FunctionKey and KeyPad on Putty.]]></Help>
          </Information>
        </Setting>
        <Setting name="Redirection After BIOS POST" order="2" selectedOption="Always Enable" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">Always Enable</Option>
              <Option value="1">BootLoader</Option>
            </AvailableOptions>
            <DefaultOption>Always Enable</DefaultOption>
            <Help><![CDATA[When Bootloader is selected, then Legacy Console Redirection is disabled before booting to legacy OS. When Always Enable is selected, then Legacy Console Redirection is enabled for legacy OS. Default setting for this option is set to Always Enable.]]></Help>
          </Information>
        </Setting>
        <Subtitle></Subtitle>
      </Menu>
      <Subtitle></Subtitle>
      <Subtitle>Legacy Console Redirection</Subtitle>
      <Setting name="Legacy Serial Redirection Port" selectedOption="COM1" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">COM1</Option>
            <Option value="1">SOL/COM2</Option>
          </AvailableOptions>
          <DefaultOption>COM1</DefaultOption>
          <Help><![CDATA[Select a COM port to display redirection of Legacy OS and Legacy OPROM Messages]]></Help>
        </Information>
      </Setting>
      <Subtitle></Subtitle>
      <Subtitle>Serial Port for Out-of-Band Management/</Subtitle>
      <Subtitle>Windows Emergency Management Services (EMS)</Subtitle>
      <Setting name="Console Redirection" order="3" checkedStatus="Unchecked" type="CheckBox">
        <!--Checked/Unchecked-->
        <Information>
          <DefaultStatus>Unchecked</DefaultStatus>
          <Help><![CDATA[Console Redirection Enable or Disable.]]></Help>
        </Information>
      </Setting>
      <Menu name="Console Redirection Settings" order="3">
        <Information>
          <Help><![CDATA[The settings specify how the host computer and the remote computer (which the user is using) will exchange data. Both computers should have the same or compatible settings.]]></Help>
          <WorkIf><![CDATA[  0 != Console Redirection$3  ]]></WorkIf>
        </Information>
        <Setting name="Out-of-Band Mgmt Port" selectedOption="COM1" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">COM1</Option>
              <Option value="1">SOL/COM2</Option>
            </AvailableOptions>
            <DefaultOption>COM1</DefaultOption>
            <Help><![CDATA[Microsoft Windows Emergency Management Services (EMS) allows for remote management of a Windows Server OS through a serial port.]]></Help>
          </Information>
        </Setting>
        <Setting name="Terminal Type" order="3" selectedOption="VT-UTF8" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">VT100</Option>
              <Option value="1">VT100+</Option>
              <Option value="2">VT-UTF8</Option>
              <Option value="3">ANSI</Option>
            </AvailableOptions>
            <DefaultOption>VT-UTF8</DefaultOption>
            <Help><![CDATA[VT-UTF8 is the preferred terminal type for out-of-band management. The next best choice is VT100+ and then VT100. See above, in Console Redirection Settings page, for more Help with Terminal Type/Emulation.]]></Help>
            <WorkIf><![CDATA[  0 != Console Redirection$3  ]]></WorkIf>
          </Information>
        </Setting>
        <Setting name="Bits per second" order="3" selectedOption="115200" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="3">9600</Option>
              <Option value="4">19200</Option>
              <Option value="6">57600</Option>
              <Option value="7">115200</Option>
            </AvailableOptions>
            <DefaultOption>115200</DefaultOption>
            <Help><![CDATA[Selects serial port transmission speed. The speed must be matched on the other side. Long or noisy lines may require lower speeds.]]></Help>
            <WorkIf><![CDATA[  0 != Console Redirection$3  ]]></WorkIf>
          </Information>
        </Setting>
        <Setting name="Flow Control" order="3" selectedOption="None" type="Option">
          <Information>
            <AvailableOptions>
              <Option value="0">None</Option>
              <Option value="1">Hardware RTS/CTS</Option>
              <Option value="2">Software Xon/Xoff</Option>
            </AvailableOptions>
            <DefaultOption>None</DefaultOption>
            <Help><![CDATA[Flow control can prevent data loss from buffer overflow. When sending data, if the receiving buffers are full, a 'stop' signal can be sent to stop the data flow. Once the buffers are empty, a 'start' signal can be sent to re-start the flow. Hardware flow control uses two wires to send start/stop signals.]]></Help>
            <WorkIf><![CDATA[  0 != Console Redirection$3  ]]></WorkIf>
          </Information>
        </Setting>
        <Text>Data Bits(8)</Text>
        <!--Valid if:   0 != Console Redirection  -->
        <Text>Parity(None)</Text>
        <!--Valid if:   0 != Console Redirection  -->
        <Text>Stop Bits(1)</Text>
        <!--Valid if:   0 != Console Redirection  -->
      </Menu>
    </Menu>
    <Menu name="ACPI Settings">
      <Information>
        <Help><![CDATA[System ACPI Parameters.]]></Help>
      </Information>
      <Subtitle>ACPI Settings</Subtitle>
      <Subtitle></Subtitle>
      <Setting name="PCI AER Support" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[$SMCUNHIDE$Enable/Disable ACPI OS to natively manage PCI Advanced Error Reporting.]]></Help>
        </Information>
      </Setting>
      <Setting name="Memory Corrected Error Enabling" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[$SMCUNHIDE$Enable/Disable Memory Corrected Error]]></Help>
        </Information>
      </Setting>
      <Setting name="NUMA" selectedOption="Enabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Enabled</DefaultOption>
          <Help><![CDATA[Enable or Disable Non uniform Memory Access (NUMA).]]></Help>
        </Information>
      </Setting>
      <Setting name="WHEA Support" selectedOption="Enabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Enabled</DefaultOption>
          <Help><![CDATA[Enable/Disable WHEA support]]></Help>
        </Information>
      </Setting>
      <Setting name="High Precision Event Timer" selectedOption="Enabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Enabled</DefaultOption>
          <Help><![CDATA[Enable or Disable the High Precision Event Timer.]]></Help>
        </Information>
      </Setting>
    </Menu>
    <Menu name="Trusted Computing">
      <Information>
        <Help><![CDATA[Trusted Computing Settings]]></Help>
      </Information>
      <Subtitle>Configuration</Subtitle>
      <Setting name="  Security Device Support" selectedOption="Enable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Enable</Option>
          </AvailableOptions>
          <DefaultOption>Enable</DefaultOption>
          <Help><![CDATA[Enables or Disables BIOS support for security device. O.S. will not show Security Device. TCG EFI protocol and INT1A interface will not be available.]]></Help>
        </Information>
      </Setting>
      <Setting name="  Disable Block Sid" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Enabled</Option>
            <Option value="0">Disabled</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[  Override to allow SID authentication in TCG Storage device]]></Help>
        </Information>
      </Setting>
      <Text>  NO Security Device Found()</Text>
    </Menu>
    <Menu name="HTTP BOOT Configuration">
      <Information>
        <Help><![CDATA[HTTP BOOT Settings]]></Help>
      </Information>
      <Subtitle>HTTP BOOT Configuration</Subtitle>
      <Setting name="Http Boot One Time" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[After create Http Boot Option, it will auto boot into HttpBoot in first time]]></Help>
        </Information>
      </Setting>
      <Setting name="Input the description" type="String">
        <Information>
          <MinSize>0</MinSize>
          <MaxSize>75</MaxSize>
          <DefaultString></DefaultString>
          <Help><![CDATA[]]></Help>
          <AllowingMultipleLine>False</AllowingMultipleLine>
        </Information>
        <StringValue><![CDATA[]]></StringValue>
      </Setting>
      <Setting name="Boot URI" type="String">
        <Information>
          <MinSize>0</MinSize>
          <MaxSize>80</MaxSize>
          <DefaultString></DefaultString>
          <Help><![CDATA[A new Boot Option will be created according to this Boot URI.]]></Help>
          <AllowingMultipleLine>False</AllowingMultipleLine>
        </Information>
        <StringValue><![CDATA[]]></StringValue>
      </Setting>
    </Menu>
    <Subtitle></Subtitle>
    <Subtitle></Subtitle>
    <Menu name="Driver Health">
      <Information>
        <Help><![CDATA[Provides Health Status for the Drivers/Controllers]]></Help>
      </Information>
      <Menu name="">
        <Information>
          <Help><![CDATA[Provides Health Status for the Drivers/Controllers]]></Help>
        </Information>
      </Menu>
    </Menu>
  </Menu>
  <Menu name="Event Logs">
    <Information />
    <Menu name="Change SMBIOS Event Log Settings">
      <Information>
        <Help><![CDATA[Press <Enter> to change the SMBIOS Event Log configuration.]]></Help>
      </Information>
      <Subtitle>Enabling/Disabling Options</Subtitle>
      <Setting name="SMBIOS Event Log" selectedOption="Enabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Enabled</DefaultOption>
          <Help><![CDATA[Change this to enable or disable all features of SMBIOS Event Logging during boot.]]></Help>
        </Information>
      </Setting>
      <Subtitle></Subtitle>
      <Subtitle>Erasing Settings</Subtitle>
      <Setting name="Erase Event Log" selectedOption="No" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">No</Option>
            <Option value="1">Yes, Next reset</Option>
            <Option value="2">Yes, Every reset</Option>
          </AvailableOptions>
          <DefaultOption>No</DefaultOption>
          <Help><![CDATA[Choose options for erasing SMBIOS Event Log.  Erasing is done prior to any logging activation during reset.]]></Help>
          <WorkIf><![CDATA[ ( 0 != SMBIOS Event Log ) ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="When Log is Full" selectedOption="Do Nothing" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Do Nothing</Option>
            <Option value="1">Erase Immediately</Option>
          </AvailableOptions>
          <DefaultOption>Do Nothing</DefaultOption>
          <Help><![CDATA[Choose options for reactions to a full SMBIOS Event Log.]]></Help>
          <WorkIf><![CDATA[ ( 0 != SMBIOS Event Log ) ]]></WorkIf>
        </Information>
      </Setting>
      <Subtitle></Subtitle>
      <Subtitle>SMBIOS Event Log Standard Settings</Subtitle>
      <Setting name="Log System Boot Event" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="1">Enabled</Option>
            <Option value="0">Disabled</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[Choose option to enable/disable logging of System boot event]]></Help>
          <WorkIf><![CDATA[ ( 0 != SMBIOS Event Log ) ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="MECI" numericValue="1" type="Numeric">
        <Information>
          <MaxValue>255</MaxValue>
          <MinValue>1</MinValue>
          <StepSize>1</StepSize>
          <DefaultValue>1</DefaultValue>
          <Help><![CDATA[Mutiple Event Count Increment:  The number of occurrences of a duplicate event that must pass before the multiple-event counter of log entry is updated.The value ranges from 1 to 255.]]></Help>
          <WorkIf><![CDATA[ ( 0 != SMBIOS Event Log ) ]]></WorkIf>
        </Information>
      </Setting>
      <Setting name="METW" numericValue="60" type="Numeric">
        <Information>
          <MaxValue>99</MaxValue>
          <MinValue>0</MinValue>
          <StepSize>1</StepSize>
          <DefaultValue>60</DefaultValue>
          <Help><![CDATA[Multiple Event Time Window:  The number of minutes which must pass between duplicate log entries which utilize a multiple-event counter. The value ranges from 0 to 99 minutes.]]></Help>
          <WorkIf><![CDATA[ ( 0 != SMBIOS Event Log ) ]]></WorkIf>
        </Information>
      </Setting>
      <Subtitle></Subtitle>
      <Subtitle>NOTE: All values changed here do not take effect</Subtitle>
      <Subtitle>      until computer is restarted.</Subtitle>
    </Menu>
    <Menu name="View SMBIOS Event Log">
      <Information>
        <Help><![CDATA[Press <Enter> to view the SMBIOS Event Log records.]]></Help>
      </Information>
      <Subtitle>DATE      TIME       ERROR CODE     SEVERITY</Subtitle>
      <Subtitle></Subtitle>
    </Menu>
    <Subtitle></Subtitle>
  </Menu>
  <Menu name="Security">
    <Information />
    <Subtitle></Subtitle>
    <Subtitle></Subtitle>
    <Subtitle></Subtitle>
    <Text>Administrator Password(Not Installed)</Text>
    <!--Valid if:   0 == Administrator Password  -->
    <Text>Administrator Password(Installed)</Text>
    <!--Valid if:   0 != Administrator Password  -->
    <Text>User Password(Not Installed)</Text>
    <!--Valid if:   0 == User Password  -->
    <Text>User Password(Installed)</Text>
    <!--Valid if:   0 != User Password  -->
    <Subtitle></Subtitle>
    <Subtitle>Password Description</Subtitle>
    <Subtitle></Subtitle>
    <Subtitle>If the Administrator's / User's password is set, </Subtitle>
    <Subtitle>then this only limits access to Setup and is</Subtitle>
    <Subtitle>asked for when entering Setup.</Subtitle>
    <Subtitle>Please set Administrator's password first in order </Subtitle>
    <Subtitle>to set User's password, if clear Administrator's </Subtitle>
    <Subtitle>password, the User's password will be cleared as well.</Subtitle>
    <Subtitle></Subtitle>
    <Subtitle>The password length must be</Subtitle>
    <Subtitle>in the following range:</Subtitle>
    <Text>Minimum length(3)</Text>
    <Text>Maximum length(20)</Text>
    <Subtitle></Subtitle>
    <Setting name="Administrator Password" type="Password">
      <Information>
        <Help>Set Administrator Password</Help>
        <MinSize>3</MinSize>
        <MaxSize>20</MaxSize>
        <HasPassword>False</HasPassword>
      </Information>
      <NewPassword><![CDATA[]]></NewPassword>
      <ConfirmNewPassword><![CDATA[]]></ConfirmNewPassword>
    </Setting>
    <Setting name="User Password" type="Password">
      <Information>
        <Help>Set User Password</Help>
        <WorkIf><![CDATA[  0 != Administrator Password  ]]></WorkIf>
        <MinSize>3</MinSize>
        <MaxSize>20</MaxSize>
        <HasPassword>False</HasPassword>
      </Information>
      <NewPassword><![CDATA[]]></NewPassword>
      <ConfirmNewPassword><![CDATA[]]></ConfirmNewPassword>
    </Setting>
    <Setting name="Password Check" selectedOption="Setup" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Setup</Option>
          <Option value="1">Always</Option>
        </AvailableOptions>
        <DefaultOption>Setup</DefaultOption>
        <Help><![CDATA[Setup: Check password while invoking setup. Always: Check password while invoking setup as well as on each boot.]]></Help>
      </Information>
    </Setting>
    <Subtitle></Subtitle>
    <Menu name="SMC Security Erase Configuration">
      <Information>
        <Help><![CDATA[]]></Help>
      </Information>
      <Text>HDD Name(SuperMicro SSD)</Text>
      <Text>HDD Serial Number(SMC0515D93017CHJ5126)</Text>
      <Text>Security Erase Mode(SAT3 Supported)</Text>
      <Text>Estimated Time(4 Minutes)</Text>
      <Text>HDD User Pwd Status:(NOT INSTALLED)</Text>
      <Setting name="Security Function" selectedOption="Disable" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disable</Option>
            <Option value="1">Security Erase</Option>
            <Option value="2">Set Password</Option>
          </AvailableOptions>
          <DefaultOption>Disable</DefaultOption>
          <Help><![CDATA[Security Erase:
SATA Devices (ATA/Non TCG) please input User Password.
SATA/Nvme Devices (TCG-Enterprise) please input EraseMaster Password.

Nvme Devices (TCG-Opal or Pyrite) please input Admin Password.
If Devices don't be set any password before, please input a Random Password.
]]></Help>
        </Information>
      </Setting>
      <Setting name="Password" type="String">
        <Information>
          <MinSize>0</MinSize>
          <MaxSize>32</MaxSize>
          <DefaultString></DefaultString>
          <Help><![CDATA[SMC HDD Security function]]></Help>
          <AllowingMultipleLine>False</AllowingMultipleLine>
        </Information>
        <StringValue><![CDATA[]]></StringValue>
      </Setting>
      <Subtitle></Subtitle>
    </Menu>
    <Subtitle></Subtitle>
    <Menu name="SMC Secure Boot Configuration">
      <Information>
        <Help><![CDATA[SecureBoot Option]]></Help>
      </Information>
      <Setting name="Secure Boot" selectedOption="Enabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Enabled</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[Secure Boot feature is Active if Secure Boot is Enabled,
Platform Key(PK) is enrolled and the System is in User mode.
The mode change requires platform reset]]></Help>
        </Information>
      </Setting>
      <Setting name="Reset Keys Type" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">Disabled</Option>
            <Option value="1">Reset all keys to default</Option>
            <Option value="2">Delete all keys</Option>
            <Option value="3">Delete PK key</Option>
          </AvailableOptions>
          <DefaultOption>Disabled</DefaultOption>
          <Help><![CDATA[Reset Keys Type]]></Help>
        </Information>
      </Setting>
      <Text>System Mode(Setup)</Text>
      <Text>Secure Boot(Not Active)</Text>
      <Setting name="Secure Boot Mode" selectedOption="Setup" type="Option">
        <Information>
          <AvailableOptions>
            <!--Option ValidIf:  (  (  ( 1 == 0 )  ||  ( 1 == 0 )  )  ||  ( 0 == 1 )  ) -->
            <Option value="0">Setup</Option>
            <!--Option ValidIf:  (  ( 1 == 0 )  ||  ( 1 == 1 )  ) -->
            <Option value="1">User</Option>
            <!--Option ValidIf:  ( 1 == 0 ) -->
            <Option value="2">Audit</Option>
            <!--Option ValidIf:  (  ( 1 == 0 )  ||  ( 1 == 1 )  ) -->
            <Option value="3">Deployed</Option>
          </AvailableOptions>
          <DefaultOption>Setup</DefaultOption>
          <Help><![CDATA[Secure Boot Mode]]></Help>
        </Information>
      </Setting>
    </Menu>
    <Subtitle>HDD Security Configuration:</Subtitle>
    <Menu name="P6:SuperMicro SSD">
      <Information>
        <Help><![CDATA[HDD Security Configuration for selected drive]]></Help>
      </Information>
      <Subtitle>HDD Password Description :</Subtitle>
      <Subtitle></Subtitle>
      <Subtitle>Allow Access to Set, Modify and Clear HardDisk User and Master Passwords.</Subtitle>
      <Subtitle>User Password need to be installed for Enabling Security.</Subtitle>
      <Subtitle>Master Password can be Modified only when successfully unlocked with Master Password in POST.</Subtitle>
      <Subtitle>If the 'Set HDD Password' option is grayed out, do power cycle to enablethe option again.</Subtitle>
      <Subtitle></Subtitle>
      <Subtitle></Subtitle>
      <Subtitle>HDD PASSWORD CONFIGURATION:</Subtitle>
      <Subtitle></Subtitle>
      <Text>Security Supported :(Yes)</Text>
      <Text>Security Supported :(No)</Text>
      <Text>Security Enabled   :(Yes)</Text>
      <Text>Security Enabled   :(No)</Text>
      <Text>Security Locked    :(Yes)</Text>
      <Text>Security Locked    :(No)</Text>
      <Text>Security Frozen    :(Yes)</Text>
      <Text>Security Frozen    :(No)</Text>
      <Text>HDD User Pwd Status:(INSTALLED)</Text>
      <Text>HDD User Pwd Status:(NOT INSTALLED)</Text>
      <Text>HDD Master Pwd Status :(INSTALLED)</Text>
      <Text>HDD Master Pwd Status :(NOT INSTALLED)</Text>
      <Subtitle></Subtitle>
      <Subtitle></Subtitle>
      <Subtitle></Subtitle>
    </Menu>
    <Subtitle></Subtitle>
    <Subtitle></Subtitle>
  </Menu>
  <Menu name="Boot">
    <Information />
    <Subtitle>Boot Configuration</Subtitle>
    <Subtitle></Subtitle>
    <Setting name="Boot mode select" selectedOption="UEFI" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">LEGACY</Option>
          <Option value="1">UEFI</Option>
          <Option value="2">DUAL</Option>
        </AvailableOptions>
        <DefaultOption>DUAL</DefaultOption>
        <Help><![CDATA[Select boot mode LEGACY/UEFI]]></Help>
      </Information>
    </Setting>
    <Setting name="LEGACY to EFI support" selectedOption="Disabled" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Disabled</Option>
          <Option value="1">Enabled</Option>
        </AvailableOptions>
        <DefaultOption>Disabled</DefaultOption>
        <Help><![CDATA[Enabled: System is able to boot to EFI OS after boot failed from Legacy boot order.]]></Help>
      </Information>
    </Setting>
    <Subtitle></Subtitle>
    <Subtitle>FIXED BOOT ORDER Priorities</Subtitle>
    <Setting name="Boot Option #1" order="1" selectedOption="UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">UEFI Hard Disk</Option>
          <Option value="1">UEFI AP</Option>
          <Option value="2">UEFI CD/DVD</Option>
          <Option value="3">UEFI USB Hard Disk</Option>
          <Option value="4">UEFI USB CD/DVD</Option>
          <Option value="5">UEFI USB Key</Option>
          <Option value="6">UEFI USB Floppy</Option>
          <Option value="7">UEFI USB Lan</Option>
          <Option value="8">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="9">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI Hard Disk</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #2" order="1" selectedOption="UEFI Hard Disk" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">UEFI Hard Disk</Option>
          <Option value="1">UEFI AP</Option>
          <Option value="2">UEFI CD/DVD</Option>
          <Option value="3">UEFI USB Hard Disk</Option>
          <Option value="4">UEFI USB CD/DVD</Option>
          <Option value="5">UEFI USB Key</Option>
          <Option value="6">UEFI USB Floppy</Option>
          <Option value="7">UEFI USB Lan</Option>
          <Option value="8">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="9">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI AP</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #3" order="1" selectedOption="UEFI CD/DVD" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">UEFI Hard Disk</Option>
          <Option value="1">UEFI AP</Option>
          <Option value="2">UEFI CD/DVD</Option>
          <Option value="3">UEFI USB Hard Disk</Option>
          <Option value="4">UEFI USB CD/DVD</Option>
          <Option value="5">UEFI USB Key</Option>
          <Option value="6">UEFI USB Floppy</Option>
          <Option value="7">UEFI USB Lan</Option>
          <Option value="8">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="9">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI CD/DVD</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #4" order="1" selectedOption="UEFI USB Hard Disk" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">UEFI Hard Disk</Option>
          <Option value="1">UEFI AP</Option>
          <Option value="2">UEFI CD/DVD</Option>
          <Option value="3">UEFI USB Hard Disk</Option>
          <Option value="4">UEFI USB CD/DVD</Option>
          <Option value="5">UEFI USB Key</Option>
          <Option value="6">UEFI USB Floppy</Option>
          <Option value="7">UEFI USB Lan</Option>
          <Option value="8">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="9">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI USB Hard Disk</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #5" order="1" selectedOption="UEFI USB CD/DVD" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">UEFI Hard Disk</Option>
          <Option value="1">UEFI AP</Option>
          <Option value="2">UEFI CD/DVD</Option>
          <Option value="3">UEFI USB Hard Disk</Option>
          <Option value="4">UEFI USB CD/DVD</Option>
          <Option value="5">UEFI USB Key</Option>
          <Option value="6">UEFI USB Floppy</Option>
          <Option value="7">UEFI USB Lan</Option>
          <Option value="8">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="9">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI USB CD/DVD</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #6" order="1" selectedOption="UEFI USB Key" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">UEFI Hard Disk</Option>
          <Option value="1">UEFI AP</Option>
          <Option value="2">UEFI CD/DVD</Option>
          <Option value="3">UEFI USB Hard Disk</Option>
          <Option value="4">UEFI USB CD/DVD</Option>
          <Option value="5">UEFI USB Key</Option>
          <Option value="6">UEFI USB Floppy</Option>
          <Option value="7">UEFI USB Lan</Option>
          <Option value="8">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="9">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI USB Key</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #7" order="1" selectedOption="UEFI USB Floppy" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">UEFI Hard Disk</Option>
          <Option value="1">UEFI AP</Option>
          <Option value="2">UEFI CD/DVD</Option>
          <Option value="3">UEFI USB Hard Disk</Option>
          <Option value="4">UEFI USB CD/DVD</Option>
          <Option value="5">UEFI USB Key</Option>
          <Option value="6">UEFI USB Floppy</Option>
          <Option value="7">UEFI USB Lan</Option>
          <Option value="8">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="9">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI USB Floppy</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #8" order="1" selectedOption="UEFI USB Lan" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">UEFI Hard Disk</Option>
          <Option value="1">UEFI AP</Option>
          <Option value="2">UEFI CD/DVD</Option>
          <Option value="3">UEFI USB Hard Disk</Option>
          <Option value="4">UEFI USB CD/DVD</Option>
          <Option value="5">UEFI USB Key</Option>
          <Option value="6">UEFI USB Floppy</Option>
          <Option value="7">UEFI USB Lan</Option>
          <Option value="8">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="9">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI USB Lan</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #9" order="1" selectedOption="UEFI AP" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">UEFI Hard Disk</Option>
          <Option value="1">UEFI AP</Option>
          <Option value="2">UEFI CD/DVD</Option>
          <Option value="3">UEFI USB Hard Disk</Option>
          <Option value="4">UEFI USB CD/DVD</Option>
          <Option value="5">UEFI USB Key</Option>
          <Option value="6">UEFI USB Floppy</Option>
          <Option value="7">UEFI USB Lan</Option>
          <Option value="8">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="9">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #1" order="2" selectedOption="Hard Disk: SuperMicro SSD      " type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>Hard Disk: SuperMicro SSD      </DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 1 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #2" order="2" selectedOption="CD/DVD" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>CD/DVD</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 1 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #3" order="2" selectedOption="USB Hard Disk" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>USB Hard Disk</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 1 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #4" order="2" selectedOption="USB CD/DVD" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>USB CD/DVD</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 1 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #5" order="2" selectedOption="USB Key" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>USB Key</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 1 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #6" order="2" selectedOption="USB Floppy" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>USB Floppy</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 1 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #7" order="2" selectedOption="USB Lan" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>USB Lan</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 1 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #8" order="2" selectedOption="Network" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>Network</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 1 2 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #1" order="3" selectedOption="Hard Disk: SuperMicro SSD      " type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="16">UEFI AP</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>Hard Disk: SuperMicro SSD      </DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #2" order="3" selectedOption="CD/DVD" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="16">UEFI AP</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>CD/DVD</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #3" order="3" selectedOption="USB Hard Disk" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="16">UEFI AP</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>USB Hard Disk</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #4" order="3" selectedOption="USB CD/DVD" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="16">UEFI AP</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>USB CD/DVD</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #5" order="3" selectedOption="USB Key" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="16">UEFI AP</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>USB Key</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #6" order="3" selectedOption="USB Floppy" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="16">UEFI AP</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>USB Floppy</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #7" order="3" selectedOption="USB Lan" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="16">UEFI AP</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>USB Lan</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #8" order="3" selectedOption="Network" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="16">UEFI AP</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>Network</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #9" order="2" selectedOption="UEFI Hard Disk" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="16">UEFI AP</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI Hard Disk</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #10" selectedOption="UEFI CD/DVD" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="16">UEFI AP</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI CD/DVD</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #11" selectedOption="UEFI USB Hard Disk" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="16">UEFI AP</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI USB Hard Disk</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #12" selectedOption="UEFI USB CD/DVD" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="16">UEFI AP</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI USB CD/DVD</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #13" selectedOption="UEFI USB Key" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="16">UEFI AP</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI USB Key</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #14" selectedOption="UEFI USB Floppy" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="16">UEFI AP</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI USB Floppy</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #15" selectedOption="UEFI USB Lan" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="16">UEFI AP</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI USB Lan</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #16" selectedOption="UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="16">UEFI AP</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Setting name="Boot Option #17" selectedOption="UEFI AP" type="Option">
      <Information>
        <AvailableOptions>
          <Option value="0">Hard Disk: SuperMicro SSD      </Option>
          <Option value="1">CD/DVD</Option>
          <Option value="2">USB Hard Disk</Option>
          <Option value="3">USB CD/DVD</Option>
          <Option value="4">USB Key</Option>
          <Option value="5">USB Floppy</Option>
          <Option value="6">USB Lan</Option>
          <Option value="7">Network</Option>
          <Option value="8">UEFI Hard Disk</Option>
          <Option value="9">UEFI CD/DVD</Option>
          <Option value="10">UEFI USB Hard Disk</Option>
          <Option value="11">UEFI USB CD/DVD</Option>
          <Option value="12">UEFI USB Key</Option>
          <Option value="13">UEFI USB Floppy</Option>
          <Option value="14">UEFI USB Lan</Option>
          <Option value="15">UEFI Network:UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28</Option>
          <Option value="16">UEFI AP</Option>
          <Option value="17">Disabled</Option>
        </AvailableOptions>
        <DefaultOption>UEFI AP</DefaultOption>
        <Help><![CDATA[Sets the system boot order]]></Help>
        <WorkIf><![CDATA[  Boot mode select is not in 0 1 3 4 5   ]]></WorkIf>
      </Information>
    </Setting>
    <Subtitle></Subtitle>
    <Subtitle></Subtitle>
    <Menu name="UEFI Application Boot Priorities">
      <Information>
        <Help><![CDATA[Specifies the Boot Device Priority sequence from available UEFI Application.]]></Help>
        <WorkIf><![CDATA[ ( Boot mode select is not in 0  ) ]]></WorkIf>
      </Information>
      <Setting name="Boot Option #1" order="4" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">UEFI: Built-in EFI Shell</Option>
            <Option value="18">Disabled</Option>
          </AvailableOptions>
          <DefaultOption>UEFI: Built-in EFI Shell</DefaultOption>
          <Help><![CDATA[Sets the system boot order]]></Help>
        </Information>
      </Setting>
    </Menu>
    <Menu name="UEFI NETWORK Drive BBS Priorities">
      <Information>
        <Help><![CDATA[Specifies the Boot Device Priority sequence from available UEFI NETWORK Drives.]]></Help>
        <WorkIf><![CDATA[ ( Boot mode select is not in 0  ) ]]></WorkIf>
      </Information>
      <Setting name="Boot Option #1" order="5" selectedOption="UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28(MAC,Address:ac1f6b7aeb76)" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28(MAC,Address:ac1f6b7aeb76)</Option>
            <Option value="1">UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28(MAC,Address:ac1f6b7aeb77)</Option>
            <Option value="2">UEFI: HTTP IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28(MAC,Address:ac1f6b7aeb76)</Option>
            <Option value="3">UEFI: HTTP IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28(MAC,Address:ac1f6b7aeb77)</Option>
            <Option value="18">Disabled</Option>
          </AvailableOptions>
          <DefaultOption>UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28(MAC,Address:ac1f6b7aeb76)</DefaultOption>
          <Help><![CDATA[Sets the system boot order]]></Help>
        </Information>
      </Setting>
      <Setting name="Boot Option #2" order="4" selectedOption="UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28(MAC,Address:ac1f6b7aeb77)" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28(MAC,Address:ac1f6b7aeb76)</Option>
            <Option value="1">UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28(MAC,Address:ac1f6b7aeb77)</Option>
            <Option value="2">UEFI: HTTP IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28(MAC,Address:ac1f6b7aeb76)</Option>
            <Option value="3">UEFI: HTTP IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28(MAC,Address:ac1f6b7aeb77)</Option>
            <Option value="18">Disabled</Option>
          </AvailableOptions>
          <DefaultOption>UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28(MAC,Address:ac1f6b7aeb77)</DefaultOption>
          <Help><![CDATA[Sets the system boot order]]></Help>
        </Information>
      </Setting>
      <Setting name="Boot Option #3" order="4" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28(MAC,Address:ac1f6b7aeb76)</Option>
            <Option value="1">UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28(MAC,Address:ac1f6b7aeb77)</Option>
            <Option value="2">UEFI: HTTP IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28(MAC,Address:ac1f6b7aeb76)</Option>
            <Option value="3">UEFI: HTTP IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28(MAC,Address:ac1f6b7aeb77)</Option>
            <Option value="18">Disabled</Option>
          </AvailableOptions>
          <DefaultOption>UEFI: HTTP IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28(MAC,Address:ac1f6b7aeb76)</DefaultOption>
          <Help><![CDATA[Sets the system boot order]]></Help>
        </Information>
      </Setting>
      <Setting name="Boot Option #4" order="4" selectedOption="Disabled" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28(MAC,Address:ac1f6b7aeb76)</Option>
            <Option value="1">UEFI: PXE IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28(MAC,Address:ac1f6b7aeb77)</Option>
            <Option value="2">UEFI: HTTP IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28(MAC,Address:ac1f6b7aeb76)</Option>
            <Option value="3">UEFI: HTTP IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28(MAC,Address:ac1f6b7aeb77)</Option>
            <Option value="18">Disabled</Option>
          </AvailableOptions>
          <DefaultOption>UEFI: HTTP IP4 Intel(R) Ethernet Controller XXV710 for 25GbE SFP28(MAC,Address:ac1f6b7aeb77)</DefaultOption>
          <Help><![CDATA[Sets the system boot order]]></Help>
        </Information>
      </Setting>
    </Menu>
    <Menu name="Hard Disk Drive BBS Priorities">
      <Information>
        <Help><![CDATA[Specifies the Boot Device Priority sequence from available Hard Disk Drives.]]></Help>
        <WorkIf><![CDATA[ ( Boot mode select is not in 1  ) ]]></WorkIf>
      </Information>
      <Setting name="Boot Option #1" order="6" selectedOption="ISATA  P6: SuperMicro SSD      (SATA,Port:6)" type="Option">
        <Information>
          <AvailableOptions>
            <Option value="0">ISATA  P6: SuperMicro SSD      (SATA,Port:6)</Option>
            <Option value="18">Disabled</Option>
          </AvailableOptions>
          <DefaultOption>ISATA  P6: SuperMicro SSD      (SATA,Port:6)</DefaultOption>
          <Help><![CDATA[Sets the system boot order]]></Help>
        </Information>
      </Setting>
    </Menu>
  </Menu>
</BiosCfg>`
)
