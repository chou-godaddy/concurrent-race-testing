module concurrent-race-testing

go 1.21

require (
	github.com/gdcorp-domains/fulfillment-generic-queue-client v1.0.4
	github.com/gdcorp-domains/fulfillment-goapimodels v1.0.89
	github.com/gdcorp-domains/fulfillment-golang-clients v1.0.121
	github.com/gdcorp-domains/fulfillment-golang-httpclient v1.0.79
	github.com/gdcorp-domains/fulfillment-golang-logging v1.0.21
	github.com/gdcorp-domains/fulfillment-golang-sso-auth v1.0.53
	github.com/gdcorp-domains/fulfillment-gosecrets v1.0.17
	github.com/gdcorp-domains/fulfillment-registry-domains v1.0.95
)

require (
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/aws/aws-sdk-go v1.53.18 // indirect
	github.com/aws/aws-xray-sdk-go v1.8.4 // indirect
	github.com/awslabs/kinesis-aggregation/go v0.0.0-20230808105340-e631fe742486 // indirect
	github.com/bradfitz/gomemcache v0.0.0-20230905024940-24af94b03874 // indirect
	github.com/bytedance/sonic v1.11.8 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/elastic/go-licenser v0.4.1 // indirect
	github.com/elastic/go-sysinfo v1.14.0 // indirect
	github.com/elastic/go-windows v1.0.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.4 // indirect
	github.com/gdcorp-domains/fulfillment-command-factory v0.14.43 // indirect
	github.com/gdcorp-domains/fulfillment-domain-eligibility v1.0.6 // indirect
	github.com/gdcorp-domains/fulfillment-go-filebuffer v1.0.0 // indirect
	github.com/gdcorp-domains/fulfillment-go2epp v1.1.30-0.20240509071517-b744b42981ec // indirect
	github.com/gdcorp-domains/fulfillment-rg-client v1.0.50 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.10.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.21.0 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/harlow/kinesis-consumer v0.3.5 // indirect
	github.com/jcchavezs/porto v0.6.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/juju/xml v0.0.0-20160224194805-b5bf18ebd8b8 // indirect
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/santhosh-tekuri/jsonschema v1.2.4 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.54.0 // indirect
	go.elastic.co/apm v1.15.0 // indirect
	go.elastic.co/apm/module/apmhttp v1.15.0 // indirect
	go.elastic.co/fastjson v1.3.0 // indirect
	golang.org/x/arch v0.8.0 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/mod v0.18.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	golang.org/x/tools v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240604185151-ef581f913117 // indirect
	google.golang.org/grpc v1.64.0 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	howett.net/plist v1.0.1 // indirect
)

// remove iam refresher lock using atomic.Value for jwt token
replace github.com/gdcorp-domains/fulfillment-golang-sso-auth v1.0.53 => /Users/chou/Desktop/godaddy/gdcorp-domains/fulfillment-golang-sso-auth

// remove httpclient wrapper lock
replace github.com/gdcorp-domains/fulfillment-golang-httpclient v1.0.79 => /Users/chou/Desktop/godaddy/gdcorp-domains/fulfillment-golang-httpclient
