module github.com/metal-stack/metal-hammer

go 1.23

require (
	github.com/beevik/ntp v1.4.3
	github.com/cheggaaa/pb/v3 v3.1.5
	github.com/google/gopacket v1.1.19
	github.com/google/uuid v1.6.0
	github.com/grafana/loki-client-go v0.0.0-20230116142646-e7494d0ef70c
	github.com/jaypipes/ghw v0.12.0
	github.com/metal-stack/go-hal v0.5.3
	github.com/metal-stack/go-lldpd v0.4.7
	github.com/metal-stack/metal-api v0.34.0
	github.com/metal-stack/metal-go v0.33.0
	github.com/metal-stack/pixie v0.3.4
	github.com/metal-stack/v v1.0.3
	// archiver must stay in version v2.1.0, see replace below
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/moby/sys/mountinfo v0.7.2
	github.com/pierrec/lz4/v4 v4.1.21
	github.com/prometheus/common v0.59.1
	github.com/samber/slog-loki/v3 v3.5.0
	github.com/samber/slog-multi v1.2.1
	github.com/u-root/u-root v0.14.0
	github.com/vishvananda/netlink v1.3.0
	golang.org/x/sync v0.8.0
	golang.org/x/sys v0.25.0
	google.golang.org/grpc v1.66.0
	google.golang.org/protobuf v1.34.2
	gopkg.in/yaml.v3 v3.0.1
)

replace (
	github.com/mholt/archiver => github.com/mholt/archiver v2.1.0+incompatible
	// This is main with my pci-ids PR merged, remove once v0.14.1 is released
	github.com/u-root/u-root => github.com/u-root/u-root v0.14.1-0.20240326203528-0d741dd84444
)

require (
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/anmitsu/go-shlex v0.0.0-20200514113438-38f4b401e2be // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/avast/retry-go/v4 v4.6.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/coreos/go-oidc/v3 v3.11.0 // indirect
	github.com/creack/pty v1.1.23 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/dennwc/varint v1.0.0 // indirect
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/frankban/quicktest v1.14.6 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/gliderlabs/ssh v0.3.7 // indirect
	github.com/go-jose/go-jose/v4 v4.0.4 // indirect
	github.com/go-kit/kit v0.13.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-openapi/analysis v0.23.0 // indirect
	github.com/go-openapi/errors v0.22.0 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/loads v0.22.0 // indirect
	github.com/go-openapi/runtime v0.28.0 // indirect
	github.com/go-openapi/spec v0.21.0 // indirect
	github.com/go-openapi/strfmt v0.23.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-openapi/validate v0.24.0 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/grafana/regexp v0.0.0-20240518133315-a468a5bfb3bc // indirect
	github.com/jaypipes/pcidb v1.0.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/lestrrat-go/blackmagic v1.0.2 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc v1.0.6 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/jwx/v2 v2.1.1 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mdlayher/ethernet v0.0.0-20220221185849-529eae5b6118 // indirect
	github.com/mdlayher/lldp v0.0.0-20150915211757-afd9f83164c5 // indirect
	github.com/metal-stack/metal-lib v0.18.1 // indirect
	github.com/metal-stack/security v0.8.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/nwaples/rardecode v1.1.3 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.20.3 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/prometheus/prometheus v0.54.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/samber/lo v1.47.0 // indirect
	github.com/samber/slog-common v0.17.1 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/sethvargo/go-password v0.3.1 // indirect
	github.com/stmcginnis/gofish v0.19.0 // indirect
	github.com/u-root/uio v0.0.0-20240224005618-d2acac8f3701 // indirect
	github.com/ulikunitz/xz v0.5.12 // indirect
	github.com/vishvananda/netns v0.0.4 // indirect
	github.com/vmware/goipmi v0.0.0-20181114221114-2333cd82d702 // indirect
	go.mongodb.org/mongo-driver v1.16.1 // indirect
	go.opentelemetry.io/otel v1.29.0 // indirect
	go.opentelemetry.io/otel/metric v1.29.0 // indirect
	go.opentelemetry.io/otel/trace v1.29.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/exp v0.0.0-20240904232852-e7e105dedf7e // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/oauth2 v0.22.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	howett.net/plist v1.0.1 // indirect
)
