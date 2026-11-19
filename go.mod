module pqv2xbench

go 1.25.6

require (
	github.com/lattice-safe/dilithium-go v0.3.0
	github.com/lattice-safe/falcon-go v0.1.1
	github.com/lattice-safe/lattice-slh-dsa-go v0.2.0
)

require (
	golang.org/x/crypto v0.51.0 // indirect
	golang.org/x/sys v0.44.0 // indirect
)

replace github.com/lattice-safe/dilithium-go => ./third_party/dilithium-go
