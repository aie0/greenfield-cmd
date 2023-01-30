package keyring

import (
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	ethHd "github.com/evmos/ethermint/crypto/hd"
)

var (
	// SupportedAlgorithms defines the list of signing algorithms used on BFS:
	//  - eth_secp256k1 (Ethereum)
	SupportedAlgorithms = keyring.SigningAlgoList{ethHd.EthSecp256k1}
	// SupportedAlgorithmsLedger defines the list of signing algorithms used on BFS for the Ledger device:
	//  - eth_secp256k1 (Ethereum)
	SupportedAlgorithmsLedger = keyring.SigningAlgoList{ethHd.EthSecp256k1}
)

// ETHAlgoOption defines a function keys options for the ethereum Secp256k1 curve.
// It supports eth_secp256k1 keys for accounts.
func ETHAlgoOption() keyring.Option {
	return func(options *keyring.Options) {
		options.SupportedAlgos = SupportedAlgorithms
		options.SupportedAlgosLedger = SupportedAlgorithmsLedger
	}
}
