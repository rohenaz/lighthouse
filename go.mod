module github.com/yourusername/lighthouse

go 1.24.3

require (
	github.com/bsv-blockchain/go-sdk v0.0.0
	github.com/spf13/cobra v1.8.0
	github.com/stretchr/testify v1.9.0
	google.golang.org/protobuf v1.32.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/bsv-blockchain/go-sdk => ../go-sdk
