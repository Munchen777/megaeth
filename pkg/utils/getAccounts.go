package utils

import (
	"crypto/ecdsa"
	"math/rand"

	log "github.com/sirupsen/logrus"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	
	"main/pkg/global"
	"main/pkg/types"
)

func privateKeyToAddress(privateKey *ecdsa.PrivateKey) (*common.Address, error) {
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKey)
	return &address, nil
}

func GetAccounts(inputs []string, onlyKeys bool) ([]types.AccountData, error) {
	var accounts []types.AccountData

	for _, input := range inputs {
		input = RemoveHexPrefix(input)

		if common.IsHexAddress("0x" + input) {
			if onlyKeys {
				log.Printf("%s | Address, Not Private Key", input)
			} else {
				accounts = append(accounts, types.AccountData{
					AccountLogData: "0x" + input,
					AccountKeyHex:  "",
					AccountKey:     nil,
					AccountAddress: common.HexToAddress("0x" + input),
				})
			}

			continue
		}

		privateKey, err := crypto.HexToECDSA(input)
		if err != nil {
			log.Printf("%s | Invalid Private Key", input)
			continue
		}

		sweepedAddress, err := privateKeyToAddress(privateKey)
		if err != nil {
			log.Printf("%s | Failed To Derive Address", input)
			continue
		}
		
		accounts = append(accounts, types.AccountData{
			AccountLogData: input,
			AccountKeyHex:  "",
			AccountKey:     privateKey,
			AccountAddress: *sweepedAddress,
		})
	}

	shuffleFlag := global.Config.ShuffleAccs
	if shuffleFlag {
		rand.Shuffle(len(accounts), func(i, j int) { accounts[i], accounts[j] = accounts[j], accounts[i] })
	}

	return accounts, nil
}
