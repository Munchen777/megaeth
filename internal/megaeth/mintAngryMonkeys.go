package megaeth

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	accTypes "main/pkg/types"
	"main/internal"
	"main/pkg/global"
	"main/pkg/utils"
)

func MintAngryMonkeys(accountData accTypes.AccountData) (bool, error) {
	log.Infof("[%d/%d] | %s | [MintAngryMonkeys] | Start Minting Angry Monkeys NFT ...\n",
		global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
	)

	var contractAddress = "0x8ac06714c0d417569bcc642cd74e48a64fe99504"

	quantity := big.NewInt(1)
	currency := common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE")
	pricePerToken := big.NewInt(1440000000000000)

	leafIndex := big.NewInt(0)
	leafAmount := new(big.Int)
	leafAmount.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
	leafAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")

	claim := accTypes.ClaimStruct{
		Proof:       []common.Hash{},
		LeafIndex:   leafIndex,
		LeafAmount:  leafAmount,
		LeafAddress: leafAddress,
	}
	signature := []byte{}

	contractABI, err := utils.LoadABI(filepath.Join("abi", "claim.json"))
	if err != nil {
		return false, err
	}

	data, err := contractABI.Pack(
		"claim",
		accountData.AccountAddress,
		quantity,
		currency,
		pricePerToken,
		claim,
		signature,
	)
	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [MintAngryMonkeys] | Problem with encoding parameters\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
		)
		log.Error(msg)
		return false, errors.New(msg)
	}

	client, ok := internal.GetClient(&accountData)
	if !ok {
		msg := fmt.Sprintf("[%d/%d] | %s | [MintAngryMonkeys] | Problem with client initialization\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
		)
		log.Error(msg)
		return false, errors.New(msg)
	}

	tx, err := client.BuildTransaction(
		contractAddress,
		data,
		big.NewInt(1440000000000000),
	)

	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [MintAngryMonkeys] | Problem with building transaction: %v\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, err,
		)
		log.Error(msg)
		return false, errors.New(msg)
	}

	chainID, err := client.GetChainID()
	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [MintAngryMonkeys] | Problem with getting chainID: %v\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, err,
		)
		log.Error(msg)
		return false, errors.New(msg)
	}

	signedTx, err := types.SignTx(
		tx,
		types.NewLondonSigner(chainID),
		client.Account.AccountKey,
	)
	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [MintAngryMonkeys] | Problem with signing tx: %v\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, err,
		)
		log.Error(msg)
		return false, errors.New(msg)
	}

	ctx := context.Background()

	err = client.Rpc.SendTransaction(ctx, signedTx)
	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [MintAngryMonkeys] | Problem with signing tx: %v\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, err,
		)
		log.Error(err)
		return false, errors.New(msg)
	} else {
		msg := fmt.Sprintf("[%d/%d] | %s | [MintAngryMonkeys] | Successfully executed transaction\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
		)
		log.Info(msg)
	}

	return true, nil
}
