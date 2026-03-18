module github.com/bruli/raspberryWaterSystem

go 1.25.5

require github.com/davecgh/go-spew v1.1.1

require gopkg.in/yaml.v3 v3.0.1

require (
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/go-chi/chi/v5 v5.2.3
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/google/uuid v1.6.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/prometheus/client_golang v1.23.2
	github.com/rs/cors v1.11.1
	github.com/stianeikeland/go-rpio/v4 v4.6.0
	periph.io/x/conn/v3 v3.7.2
	periph.io/x/devices/v3 v3.7.4
	periph.io/x/host/v3 v3.8.5
)

require (
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/charmbracelet/lipgloss v1.0.0 // indirect
	github.com/charmbracelet/x/ansi v0.8.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/matryer/moq v0.6.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mfridman/tparse v0.17.0 // indirect
	github.com/muesli/termenv v0.15.3-0.20240618155329-98d742f6907a // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.66.1 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/santhosh-tekuri/jsonschema/cmd/jv v0.7.0 // indirect
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	golang.org/x/mod v0.26.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/telemetry v0.0.0-20250710130107-8d8967aff50b // indirect
	golang.org/x/text v0.28.0 // indirect
	golang.org/x/tools v0.35.0 // indirect
	golang.org/x/tools/go/expect v0.1.1-deprecated // indirect
	golang.org/x/tools/go/packages/packagestest v0.1.1-deprecated // indirect
	golang.org/x/vuln v1.1.4 // indirect
	google.golang.org/protobuf v1.36.8 // indirect
	mvdan.cc/gofumpt v0.8.0 // indirect
)

require (
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/rs/zerolog v1.32.0
	github.com/stretchr/testify v1.11.1
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
)

tool (
	github.com/matryer/moq
	github.com/mfridman/tparse
	github.com/santhosh-tekuri/jsonschema/cmd/jv
	golang.org/x/vuln/cmd/govulncheck
	mvdan.cc/gofumpt
)
