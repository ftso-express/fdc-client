module local/fdc

go 1.21.5

require (
	flare-common v0.0.0-00010101000000-000000000000
	github.com/ethereum/go-ethereum v1.13.15
)

replace flare-common => ./flare-common

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/holiman/uint256 v1.2.4 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	gorm.io/driver/mysql v1.5.6 // indirect
	gorm.io/gorm v1.25.9 // indirect
)
