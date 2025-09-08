package megaeth

import (
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	accTypes "main/pkg/types"
	"main/internal"
	"main/pkg/global"
)

func (f FunNFT) MintFunNFT(ctx context.Context, accountData accTypes.AccountData) (bool, error) {
	log.Infof("[%d/%d] | %s | [%s] | Start minting ...\n",
		global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, f.DisplayName,
	)

	var contractAddress = f.ContractAddress

	client, ok := internal.GetClient(&accountData)
	if !ok {
		msg := fmt.Sprintf("[%d/%d] | %s | [%s] | Problem with client initialization\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, f.DisplayName,
		)
		log.Error(msg)
		return false, errors.New(msg)
	}

	tx, err := client.BuildTransaction(
		contractAddress,
		common.FromHex(f.Data),
		f.Value,
	)

	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [%s] | Problem with building transaction: %v\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, f.DisplayName, err,
		)
		log.Error(msg)
		return false, errors.New(msg)
	}

	chainID, err := client.GetChainID()
	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [%s] | Problem with getting chainID: %v\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, f.DisplayName, err,
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
		msg := fmt.Sprintf("[%d/%d] | %s | [%s] | Problem with signing tx: %v\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, f.DisplayName, err,
		)
		log.Error(msg)
		return false, errors.New(msg)
	}

	err = client.Rpc.SendTransaction(ctx, signedTx)
	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [%s] | Problem with signing tx: %v\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, f.DisplayName, err,
		)
		log.Error(msg)
		return false, errors.New(msg)
	} else {
		msg := fmt.Sprintf("[%d/%d] | %s | [%s] | Successfully executed transaction\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, f.DisplayName,
		)
		log.Info(msg)
	}

	return true, nil
}
