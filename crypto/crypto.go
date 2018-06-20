package crypto

import (
	sha3 "github.com/bankex/go-plasma/crypto/sha3"
	"github.com/ethereum/go-ethereum/common"
)

var EmptyAddress = common.Address{}

func PubkeyToAddress(pubBytes []byte) common.Address {
	if len(pubBytes) != 64 {
		return EmptyAddress
	}
	return common.BytesToAddress(sha3.Keccak256(pubBytes[1:])[12:])
}