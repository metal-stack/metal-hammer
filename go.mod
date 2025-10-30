module github.com/metal-stack/metal-hammer

go 1.25

require (
	github.com/beevik/ntp v1.4.3
	github.com/cheggaaa/pb/v3 v3.1.7
	github.com/google/gopacket v1.1.19
	github.com/google/uuid v1.6.0
	github.com/grafana/loki-client-go v0.0.0-20240913122146-e119d400c3a5
	github.com/jaypipes/ghw v0.19.1
	github.com/metal-stack/go-hal v0.6.1-0.20251030141306-b0c9bd1d7973
	github.com/metal-stack/go-lldpd v0.4.10
	github.com/metal-stack/metal-api v0.42.2
	github.com/metal-stack/metal-go v0.42.2
	github.com/metal-stack/pixie v0.3.6
	github.com/metal-stack/v v1.0.3
	// archiver must stay in version v2.1.0, see replace below
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/moby/sys/mountinfo v0.7.2
	github.com/pierrec/lz4/v4 v4.1.22
	github.com/pkg/term v1.1.0
	github.com/prometheus/common v0.65.0
	github.com/samber/slog-loki/v3 v3.5.4
	github.com/samber/slog-multi v1.4.1
	github.com/u-root/u-root v0.15.0
	github.com/vishvananda/netlink v1.3.1
	golang.org/x/sync v0.16.0
	golang.org/x/sys v0.35.0
	google.golang.org/grpc v1.75.0
	google.golang.org/protobuf v1.36.8
	gopkg.in/yaml.v3 v3.0.1
)

replace (
	github.com/mholt/archiver => github.com/mholt/archiver v2.1.0+incompatible
	// keep this until https://github.com/u-root/u-root/pull/3451 is merged and released
	// github.com/u-root/u-root => github.com/majst01/u-root v0.0.0-20250910091544-306665b6f8e8
	github.com/u-root/u-root => github.com/u-root/u-root v0.14.1-0.20240326203528-0d741dd84444
)

require (
	github.com/Masterminds/semver/v3 v3.4.0 // indirect
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/anmitsu/go-shlex v0.0.0-20200514113438-38f4b401e2be // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/avast/retry-go/v4 v4.6.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/coreos/go-oidc/v3 v3.15.0 // indirect
	github.com/creack/pty v1.1.24 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.4.0 // indirect
	github.com/dennwc/varint v1.0.0 // indirect
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/frankban/quicktest v1.14.6 // indirect
	github.com/gliderlabs/ssh v0.3.8 // indirect
	github.com/go-jose/go-jose/v4 v4.1.2 // indirect
	github.com/go-kit/kit v0.13.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-openapi/analysis v0.23.0 // indirect
	github.com/go-openapi/errors v0.22.2 // indirect
	github.com/go-openapi/jsonpointer v0.22.0 // indirect
	github.com/go-openapi/jsonreference v0.21.1 // indirect
	github.com/go-openapi/loads v0.22.0 // indirect
	github.com/go-openapi/runtime v0.28.0 // indirect
	github.com/go-openapi/spec v0.21.0 // indirect
	github.com/go-openapi/strfmt v0.23.0 // indirect
	github.com/go-openapi/swag v0.24.1 // indirect
	github.com/go-openapi/swag/cmdutils v0.24.0 // indirect
	github.com/go-openapi/swag/conv v0.24.0 // indirect
	github.com/go-openapi/swag/fileutils v0.24.0 // indirect
	github.com/go-openapi/swag/jsonname v0.24.0 // indirect
	github.com/go-openapi/swag/jsonutils v0.24.0 // indirect
	github.com/go-openapi/swag/loading v0.24.0 // indirect
	github.com/go-openapi/swag/mangling v0.24.0 // indirect
	github.com/go-openapi/swag/netutils v0.24.0 // indirect
	github.com/go-openapi/swag/stringutils v0.24.0 // indirect
	github.com/go-openapi/swag/typeutils v0.24.0 // indirect
	github.com/go-openapi/swag/yamlutils v0.24.0 // indirect
	github.com/go-openapi/validate v0.24.0 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/grafana/loki/pkg/push v0.0.0-20250903093132-95ced86f69c5 // indirect
	github.com/grafana/regexp v0.0.0-20240518133315-a468a5bfb3bc // indirect
	github.com/jaypipes/pcidb v1.1.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/lestrrat-go/blackmagic v1.0.4 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc v1.0.6 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/jwx/v2 v2.1.6 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mdlayher/ethernet v0.0.0-20220221185849-529eae5b6118 // indirect
	github.com/mdlayher/lldp v0.0.0-20150915211757-afd9f83164c5 // indirect
	github.com/metal-stack/metal-lib v0.23.4 // indirect
	github.com/metal-stack/security v0.9.4 // indirect
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
	github.com/prometheus/client_golang v1.23.0 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/procfs v0.17.0 // indirect
	github.com/prometheus/prometheus v0.304.0 // indirect
	github.com/rekby/gpt v0.0.0-20200614112001-7da10aec5566 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/samber/lo v1.51.0 // indirect
	github.com/samber/slog-common v0.19.0 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/sethvargo/go-password v0.3.1 // indirect
	github.com/stmcginnis/gofish v0.20.0 // indirect
	github.com/u-root/uio v0.0.0-20240224005618-d2acac8f3701 // indirect
	github.com/ulikunitz/xz v0.5.15 // indirect
	github.com/vishvananda/netns v0.0.5 // indirect
	github.com/vmware/goipmi v0.0.0-20181114221114-2333cd82d702 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.mongodb.org/mongo-driver v1.17.4 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.38.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/exp v0.0.0-20250819193227-8b4c13bb791b // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/oauth2 v0.30.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250826171959-ef028d996bc1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	howett.net/plist v1.0.2-0.20250314012144-ee69052608d9 // indirect
)
