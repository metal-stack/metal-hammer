package ipmi

import (
	"reflect"
	"testing"
)

// Output of root@ipmitest:~# ipmitool lan print
const lanPrint = `
Set in Progress         : Set Complete
Auth Type Support       : NONE MD2 MD5 PASSWORD
Auth Type Enable        : Callback : MD2 MD5 PASSWORD
                        : User     : MD2 MD5 PASSWORD
                        : Operator : MD2 MD5 PASSWORD
                        : Admin    : MD2 MD5 PASSWORD
                        : OEM      : MD2 MD5 PASSWORD
IP Address Source       : Static Address
IP Address              : 10.248.36.246
Subnet Mask             : 255.255.252.0
MAC Address             : 0c:c4:7a:ed:3e:27
SNMP Community String   : public
IP Header               : TTL=0x00 Flags=0x00 Precedence=0x00 TOS=0x00
BMC ARP Control         : ARP Responses Enabled, Gratuitous ARP Disabled
Default Gateway IP      : 10.248.36.1
Default Gateway MAC     : 30:b6:4f:c3:a0:c1
Backup Gateway IP       : 0.0.0.0
Backup Gateway MAC      : 00:00:00:00:00:00
802.1q VLAN ID          : Disabled
802.1q VLAN Priority    : 0
RMCP+ Cipher Suites     : 1,2,3,6,7,8,11,12
Cipher Suite Priv Max   : XaaaXXaaaXXaaXX
                        :     X=Cipher Suite Unused
                        :     c=CALLBACK
                        :     u=USER
                        :     o=OPERATOR
                        :     a=ADMIN
                        :     O=OEM
Bad Password Threshold  : Not Available
`
const lanPrint2 = "Set in Progress         : Set Complete\nAuth Type Support       : NONE MD2 MD5 PASSWORD \nAuth Type Enable        : Callback : MD2 MD5 PASSWORD \n                        : User     : MD2 MD5 PASSWORD \n                        : Operator : MD2 MD5 PASSWORD \n                        : Admin    : MD2 MD5 PASSWORD \n                        : OEM      : MD2 MD5 PASSWORD \nIP Address Source       : DHCP Address\nIP Address              : 192.168.2.53\nSubnet Mask             : 255.255.255.0\nMAC Address             : ac:1f:6b:73:c9:f0\nSNMP Community String   : public\nIP Header               : TTL=0x00 Flags=0x00 Precedence=0x00 TOS=0x00\nBMC ARP Control         : ARP Responses Enabled, Gratuitous ARP Disabled\nDefault Gateway IP      : 192.168.2.1\nDefault Gateway MAC     : 00:00:00:00:00:00\nBackup Gateway IP       : 0.0.0.0\nBackup Gateway MAC      : 00:00:00:00:00:00\n802.1q VLAN ID          : Disabled\n802.1q VLAN Priority    : 0\nRMCP+ Cipher Suites     : 1,2,3,6,7,8,11,12\nCipher Suite Priv Max   : XaaaXXaaaXXaaXX\n                        :     X=Cipher Suite Unused\n                        :     c=CALLBACK\n                        :     u=USER\n                        :     o=OPERATOR\n                        :     a=ADMIN\n                        :     O=OEM\nBad Password Threshold  : 0\nInvalid password disable: no\nAttempt Count Reset Int.: 0\nUser Lockout Interval   : 0\n"

func Test_getLanConfig(t *testing.T) {
	tests := []struct {
		name      string
		cmdOutput string
		want      map[string]string
	}{
		{
			name:      "simple",
			cmdOutput: lanPrint,
			want: map[string]string{
				"IP Address":  "10.248.36.246",
				"Subnet Mask": "255.255.252.0",
				"MAC Address": "0c:c4:7a:ed:3e:27",
			},
		},
		{
			name:      "from real",
			cmdOutput: lanPrint2,
			want: map[string]string{
				"IP Address":  "192.168.2.53",
				"Subnet Mask": "255.255.255.0",
				"MAC Address": "ac:1f:6b:73:c9:f0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := output2Map(tt.cmdOutput)
			for key, value := range tt.want {
				if got[key] != value {
					t.Errorf("getLanConfig() = %v, want %v", got[key], value)
				}
			}
		})
	}
}

func TestGetLanConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    LanConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Ipmitool{Command: "/bin/true"}
			got, err := i.GetLanConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLanConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLanConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLanConfig_From(t *testing.T) {
	type fields struct {
		IP  string
		Mac string
	}
	tests := []struct {
		name   string
		fields fields
		input  map[string]string
	}{
		{
			name: "simple",
			fields: fields{
				IP:  "192.168.2.53",
				Mac: "ac:1f:6b:73:c9:f0",
			},
			input: output2Map(lanPrint2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LanConfig{}
			from(c, tt.input)
			if c.IP != tt.fields.IP {
				t.Errorf("IP is not as expected")
			}
			if c.Mac != tt.fields.Mac {
				t.Errorf("Mac is not as expected")
			}
		})
	}
}
