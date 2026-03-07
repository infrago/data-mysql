module github.com/infrago/data-mysql

go 1.25.3

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/infrago/data v0.12.0
	github.com/infrago/infra v0.12.0
)

replace github.com/infrago/data => ../data

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/infrago/base v0.12.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
