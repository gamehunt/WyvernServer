module wyvern/server

go 1.25.6

require (
	github.com/bytemare/opaque v0.10.0
	github.com/c-bata/go-prompt v0.2.6
	github.com/valkey-io/valkey-go v1.0.71
	go.mongodb.org/mongo-driver/v2 v2.5.0
)

ignore (
	./init
	./storage
)

require (
	filippo.io/edwards25519 v1.0.0 // indirect
	filippo.io/nistec v0.0.2 // indirect
	github.com/bytemare/crypto v0.4.3 // indirect
	github.com/bytemare/hash v0.1.5 // indirect
	github.com/bytemare/hash2curve v0.1.3 // indirect
	github.com/bytemare/ksf v0.1.0 // indirect
	github.com/gtank/ristretto255 v0.1.2 // indirect
	github.com/klauspost/compress v1.17.6 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mattn/go-tty v0.0.3 // indirect
	github.com/pkg/term v1.2.0-beta.2 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.2.0 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
)
