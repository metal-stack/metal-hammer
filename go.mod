module github.com/metal-stack/metal-hammer

go 1.19

require (
	github.com/beevik/ntp v0.3.0
	github.com/cheggaaa/pb/v3 v3.1.0
	github.com/google/gopacket v1.1.19
	github.com/google/uuid v1.3.0
	github.com/jaypipes/ghw v0.9.0
	github.com/metal-stack/go-hal v0.4.3
	github.com/metal-stack/go-lldpd v0.4.2
	github.com/metal-stack/metal-api v0.21.5-0.20221108201044-00cea8081c41
	github.com/metal-stack/metal-go v0.22.0
	github.com/metal-stack/pixie v0.2.2
	github.com/metal-stack/v v1.0.3
	// archiver must stay in version v2.1.0, see replace below
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/pierrec/lz4/v4 v4.1.17
	github.com/u-root/u-root v0.10.0
	github.com/vishvananda/netlink v1.2.1-beta.2
	go.uber.org/zap v1.24.0
	golang.org/x/sync v0.1.0
	golang.org/x/sys v0.4.0
	google.golang.org/grpc v1.52.3
	google.golang.org/protobuf v1.28.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/anmitsu/go-shlex v0.0.0-20200514113438-38f4b401e2be // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/avast/retry-go/v4 v4.3.2 // indirect
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/coreos/go-oidc/v3 v3.5.0 // indirect
	github.com/creack/pty v1.1.18 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.1.0 // indirect
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/fatih/color v1.14.1 // indirect
	github.com/frankban/quicktest v1.14.4 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/gliderlabs/ssh v0.3.5 // indirect
	github.com/go-jose/go-jose/v3 v3.0.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-openapi/analysis v0.21.4 // indirect
	github.com/go-openapi/errors v0.20.3 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/loads v0.21.2 // indirect
	github.com/go-openapi/runtime v0.25.0 // indirect
	github.com/go-openapi/spec v0.20.8 // indirect
	github.com/go-openapi/strfmt v0.21.3 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/go-openapi/validate v0.22.1 // indirect
	github.com/goccy/go-json v0.10.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.4.3 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/jaypipes/pcidb v1.0.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.8 // indirect
	github.com/lestrrat-go/blackmagic v1.0.1 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/jwx v1.2.25 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/mdlayher/ethernet v0.0.0-20220221185849-529eae5b6118 // indirect
	github.com/mdlayher/lldp v0.0.0-20150915211757-afd9f83164c5 // indirect
	github.com/metal-stack/metal-lib v0.11.3 // indirect
	github.com/metal-stack/security v0.6.6 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/nwaples/rardecode v1.1.3 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rivo/uniseg v0.4.3 // indirect
	github.com/sethvargo/go-password v0.2.0 // indirect
	github.com/stmcginnis/gofish v0.13.0 // indirect
	github.com/ulikunitz/xz v0.5.11 // indirect
	github.com/vishvananda/netns v0.0.4 // indirect
	github.com/vmware/goipmi v0.0.0-20181114221114-2333cd82d702 // indirect
	go.mongodb.org/mongo-driver v1.11.1 // indirect
	go.opentelemetry.io/otel v1.11.1 // indirect
	go.opentelemetry.io/otel/trace v1.11.1 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/oauth2 v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230127162408-596548ed4efa // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	howett.net/plist v1.0.0 // indirect
)

replace github.com/mholt/archiver => github.com/mholt/archiver v2.1.0+incompatible
