package nnsmetadata

type Network string

const (
	Mainnet Network = "mainnet"
	Goerli  Network = "goerli"
)

// ParseNetwork parses the given string into a Network.
func ParseNetwork(s string) (Network, bool) {
	switch s {
	case string(Mainnet):
		return Mainnet, true
	case string(Goerli):
		return Goerli, true
	default:
		return "", false
	}
}
