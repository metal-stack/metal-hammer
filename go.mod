module github.com/metal-stack/metal-hammer

go 1.27

replace (
	github.com/mholt/archiver => github.com/mholt/archiver v2.1.0+incompatible
	github.com/prometheus/prometheus => github.com/prometheus/prometheus v0.304.0
	// keep this until https://github.com/u-root/u-root/pull/3451 is merged and released
	github.com/u-root/u-root => github.com/majst01/u-root v0.0.0-20250910091544-306665b6f8e8
)

require (
	github.com/beevik/ntp v1.5.0
	github.com/cheggaaa/pb/v3 v3.2.1
	github.com/google/go-cmp v0.7.0
	github.com/google/gopacket v1.1.19
	github.com/grafana/loki-client-go v0.0.0-20260414011004-43add134e848
	github.com/jaypipes/ghw v0.25.0
	github.com/metal-stack/api v0.5.5
	github.com/metal-stack/go-hal v0.7.3
	github.com/metal-stack/go-lldpd v0.4.12
	github.com/metal-stack/metal-lib v0.26.3
	github.com/metal-stack/os-installer v0.3.1
	github.com/metal-stack/pixie v0.4.2-0.20260905145935-f1a73f077d13
	github.com/metal-stack/v v1.0.3
	// archiver must stay in version v2.1.0, see replace below
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/moby/sys/mountinfo v0.7.2
	github.com/pierrec/lz4/v4 v4.1.29
	github.com/prometheus/common v0.71.0
	github.com/samber/slog-loki/v3 v3.7.2
	github.com/samber/slog-multi v1.8.0
	github.com/u-root/u-root v0.15.0
	github.com/vishvananda/netlink v1.3.1
	golang.org/x/sync v0.22.0
	golang.org/x/sys v0.47.0
	google.golang.org/protobuf v1.36.12
	gopkg.in/yaml.v3 v3.0.1
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.12-20260825204119-511051f7f437.2 // indirect
	buf.build/go/protovalidate v1.4.0 // indirect
	buf.build/go/protoyaml v0.7.0 // indirect
	cel.dev/cel-go v0.32.0 // indirect
	cel.dev/expr v0.25.3 // indirect
	connectrpc.com/connect v1.20.0 // indirect
	dario.cat/mergo v1.0.2 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.5.0 // indirect
	github.com/Masterminds/sprig/v3 v3.3.0 // indirect
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/ajeddeloh/go-json v0.0.0-20160803184958-73d058cf8437 // indirect
	github.com/anmitsu/go-shlex v0.0.0-20200514113438-38f4b401e2be // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/avast/retry-go/v4 v4.7.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/containerd/platforms v1.0.0-rc.5 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/go-systemd/v22 v22.7.0 // indirect
	github.com/creack/pty v1.1.24 // indirect
	github.com/dennwc/varint v1.0.0 // indirect
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/fatih/color v1.19.0 // indirect
	github.com/flatcar/ignition v0.36.2 // indirect
	github.com/gliderlabs/ssh v0.3.8 // indirect
	github.com/go-kit/kit v0.13.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.1 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/godbus/dbus/v5 v5.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grafana/loki/pkg/push v0.0.0-20260907023342-c4926dbc041c // indirect
	github.com/grafana/regexp v0.0.0-20250905093917-f7b3be9d1853 // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/jaypipes/pcidb v1.1.1 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.20.0 // indirect
	github.com/klauspost/connect-compress/v2 v2.1.1 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/mattn/go-runewidth v0.0.29 // indirect
	github.com/mdlayher/ethernet v0.0.0-20220221185849-529eae5b6118 // indirect
	github.com/mdlayher/lldp v0.0.0-20150915211757-afd9f83164c5 // indirect
	github.com/minio/minlz v1.2.0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.3-0.20250322232337-35a7c28c31ee // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/nwaples/rardecode v1.1.3 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.24.1 // indirect
	github.com/prometheus/client_model v0.6.3 // indirect
	github.com/prometheus/procfs v0.22.0 // indirect
	github.com/prometheus/prometheus v0.314.0 // indirect
	github.com/rekby/gpt v0.0.0-20200614112001-7da10aec5566 // indirect
	github.com/samber/lo v1.53.0 // indirect
	github.com/samber/slog-common v0.22.0 // indirect
	github.com/sethvargo/go-password v0.4.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/stmcginnis/gofish v0.25.0 // indirect
	github.com/u-root/uio v0.0.0-20240224005618-d2acac8f3701 // indirect
	github.com/ulikunitz/xz v0.5.16 // indirect
	github.com/vincent-petithory/dataurl v1.0.0 // indirect
	github.com/vishvananda/netns v0.0.5 // indirect
	github.com/vmware/goipmi v0.0.0-20181114221114-2333cd82d702 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.yaml.in/yaml/v2 v2.4.4 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	go4.org v0.0.0-20260112195520-a5071408f32f // indirect
	golang.org/x/crypto v0.56.0 // indirect
	golang.org/x/exp v0.0.0-20260824195058-e88cd73687aa // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/oauth2 v0.36.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260904194346-d0f1323225a4 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260904194346-d0f1323225a4 // indirect
	google.golang.org/grpc v1.83.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	howett.net/plist v1.0.2-0.20250314012144-ee69052608d9 // indirect
)
