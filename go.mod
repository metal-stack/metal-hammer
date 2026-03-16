module github.com/metal-stack/metal-hammer

go 1.26

require (
	github.com/beevik/ntp v1.5.0
	github.com/cheggaaa/pb/v3 v3.1.7
	github.com/google/go-cmp v0.7.0
	github.com/google/gopacket v1.1.19
	github.com/google/uuid v1.6.0
	github.com/grafana/loki-client-go v0.0.0-20251015150631-c42bbddc310a
	github.com/jaypipes/ghw v0.21.2
	github.com/metal-stack/api v0.0.54-0.20260309104254-e1a94cd811ff
	github.com/metal-stack/go-hal v0.7.0
	github.com/metal-stack/go-lldpd v0.4.11
	github.com/metal-stack/metal-api v0.43.0
	github.com/metal-stack/metal-go v0.43.1-0.20260313140852-e1dafffd26a3
	github.com/metal-stack/metal-lib v0.24.0
	github.com/metal-stack/os-installer v0.2.1-0.20260316114451-0d0d563b7d9c
	github.com/metal-stack/pixie v0.4.1
	github.com/metal-stack/v v1.0.3
	// archiver must stay in version v2.1.0, see replace below
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/moby/sys/mountinfo v0.7.2
	github.com/pierrec/lz4/v4 v4.1.26
	github.com/prometheus/common v0.67.5
	github.com/samber/slog-loki/v3 v3.7.1
	github.com/samber/slog-multi v1.7.1
	github.com/u-root/u-root v0.15.0
	github.com/vishvananda/netlink v1.3.1
	golang.org/x/sync v0.20.0
	golang.org/x/sys v0.42.0
	google.golang.org/grpc v1.79.2
	google.golang.org/protobuf v1.36.11
	gopkg.in/yaml.v3 v3.0.1
)

replace (
	github.com/mholt/archiver => github.com/mholt/archiver v2.1.0+incompatible
	// keep this until https://github.com/u-root/u-root/pull/3451 is merged and released
	github.com/u-root/u-root => github.com/majst01/u-root v0.0.0-20250910091544-306665b6f8e8
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.11-20260209202127-80ab13bee0bf.1 // indirect
	buf.build/go/protovalidate v1.1.3 // indirect
	buf.build/go/protoyaml v0.6.0 // indirect
	cel.dev/expr v0.25.1 // indirect
	dario.cat/mergo v1.0.2 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.4.0 // indirect
	github.com/Masterminds/sprig/v3 v3.3.0 // indirect
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/ajeddeloh/go-json v0.0.0-20160803184958-73d058cf8437 // indirect
	github.com/anmitsu/go-shlex v0.0.0-20200514113438-38f4b401e2be // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/avast/retry-go/v4 v4.7.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/clipperhouse/stringish v0.1.1 // indirect
	github.com/clipperhouse/uax29/v2 v2.3.0 // indirect
	github.com/coreos/go-oidc/v3 v3.17.0 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/go-systemd/v22 v22.7.0 // indirect
	github.com/creack/pty v1.1.24 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.4.1 // indirect
	github.com/dennwc/varint v1.0.0 // indirect
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/flatcar/ignition v0.36.2 // indirect
	github.com/gliderlabs/ssh v0.3.8 // indirect
	github.com/go-jose/go-jose/v4 v4.1.3 // indirect
	github.com/go-kit/kit v0.13.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-openapi/analysis v0.24.3 // indirect
	github.com/go-openapi/errors v0.22.7 // indirect
	github.com/go-openapi/jsonpointer v0.22.5 // indirect
	github.com/go-openapi/jsonreference v0.21.5 // indirect
	github.com/go-openapi/loads v0.23.3 // indirect
	github.com/go-openapi/runtime v0.29.3 // indirect
	github.com/go-openapi/spec v0.22.4 // indirect
	github.com/go-openapi/strfmt v0.26.0 // indirect
	github.com/go-openapi/swag v0.25.5 // indirect
	github.com/go-openapi/swag/cmdutils v0.25.5 // indirect
	github.com/go-openapi/swag/conv v0.25.5 // indirect
	github.com/go-openapi/swag/fileutils v0.25.5 // indirect
	github.com/go-openapi/swag/jsonname v0.25.5 // indirect
	github.com/go-openapi/swag/jsonutils v0.25.5 // indirect
	github.com/go-openapi/swag/loading v0.25.5 // indirect
	github.com/go-openapi/swag/mangling v0.25.5 // indirect
	github.com/go-openapi/swag/netutils v0.25.5 // indirect
	github.com/go-openapi/swag/stringutils v0.25.5 // indirect
	github.com/go-openapi/swag/typeutils v0.25.5 // indirect
	github.com/go-openapi/swag/yamlutils v0.25.5 // indirect
	github.com/go-openapi/validate v0.25.2 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/goccy/go-json v0.10.6 // indirect
	github.com/godbus/dbus/v5 v5.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/google/cel-go v0.27.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/grafana/loki/pkg/push v0.0.0-20250903093132-95ced86f69c5 // indirect
	github.com/grafana/regexp v0.0.0-20240518133315-a468a5bfb3bc // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/jaypipes/pcidb v1.1.1 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/lestrrat-go/blackmagic v1.0.4 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc/v3 v3.0.4 // indirect
	github.com/lestrrat-go/jwx/v3 v3.0.13 // indirect
	github.com/lestrrat-go/option/v2 v2.0.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.19 // indirect
	github.com/mdlayher/ethernet v0.0.0-20220221185849-529eae5b6118 // indirect
	github.com/mdlayher/lldp v0.0.0-20150915211757-afd9f83164c5 // indirect
	github.com/metal-stack/security v0.9.6 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/nwaples/rardecode v1.1.3 // indirect
	github.com/oklog/ulid/v2 v2.1.1 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.23.2 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/procfs v0.19.1 // indirect
	github.com/prometheus/prometheus v0.304.0 // indirect
	github.com/rekby/gpt v0.0.0-20200614112001-7da10aec5566 // indirect
	github.com/samber/lo v1.53.0 // indirect
	github.com/samber/slog-common v0.20.0 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/sethvargo/go-password v0.3.1 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/stmcginnis/gofish v0.21.3 // indirect
	github.com/u-root/uio v0.0.0-20240224005618-d2acac8f3701 // indirect
	github.com/ulikunitz/xz v0.5.15 // indirect
	github.com/vincent-petithory/dataurl v1.0.0 // indirect
	github.com/vishvananda/netns v0.0.5 // indirect
	github.com/vmware/goipmi v0.0.0-20181114221114-2333cd82d702 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.42.0 // indirect
	go.opentelemetry.io/otel/metric v1.42.0 // indirect
	go.opentelemetry.io/otel/trace v1.42.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.yaml.in/yaml/v2 v2.4.4 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	go4.org v0.0.0-20260112195520-a5071408f32f // indirect
	golang.org/x/crypto v0.49.0 // indirect
	golang.org/x/exp v0.0.0-20260312153236-7ab1446f8b90 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/oauth2 v0.36.0 // indirect
	golang.org/x/text v0.35.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260311181403-84a4fc48630c // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260311181403-84a4fc48630c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	howett.net/plist v1.0.2-0.20250314012144-ee69052608d9 // indirect
)
