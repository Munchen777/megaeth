package types

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type ModuleFunction func(AccountData) (bool, error)

type AccountData struct {
	AccountKeyHex  string
	AccountKey     *ecdsa.PrivateKey
	AccountAddress common.Address
	AccountLogData string
}

type ClaimStruct struct {
	Proof       []common.Hash
	LeafIndex   *big.Int
	LeafAmount  *big.Int
	LeafAddress common.Address
}

type Settings struct {
	DelayBeforeStart struct {
		Min int `yaml:"min"`
		Max int `yaml:"max"`
	} `yaml:"delay_before_start"`

	DelayBetweenAccs struct {
		Min int `yaml:"min"`
		Max int `yaml:"max"`
	} `yaml:"delay_between_accs"`

	CapmonsterAPIKey string `yaml:"capmonster_api_key"`
	ShuffleAccs bool `yaml:"shuffle_accs"`
}
