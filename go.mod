module github.com/infrago/data-mysql

go 1.25.3

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/infrago/data v0.10.0
	github.com/infrago/infra v0.10.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/infrago/base v0.10.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
)

replace github.com/infrago/infra => ../bamgoo

replace github.com/infrago/data => ../data
