package megaeth

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	log "github.com/sirupsen/logrus"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	accTypes "main/pkg/types"
	"main/internal"
	"main/pkg/global"
)

func MintFunNFT(accountData accTypes.AccountData) (bool, error) {
	log.Infof("[%d/%d] | %s | [MintNFT] | Start Minting NFT ...\n",
		global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
	)

	client, ok := internal.GetClient(&accountData)
	if !ok {
		msg := fmt.Sprintf("[%d/%d] | %s | [MintNFT] | Problem with client initialization\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
		)
		return false, errors.New(msg)
	}

	tx, err := client.BuildTransaction(
		"0xb8027dca96746f073896c45f65b720f9bd2afee7",
		common.FromHex("0x1249c58b"),
		big.NewInt(0),
	)

	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [MintNFT] | Problem with building transaction: %v\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, err,
		)
		return false, errors.New(msg)
	}

	chainID, err := client.GetChainID()
	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [MintNFT] | Problem with getting chainID: %v\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, err,
		)
		return false, errors.New(msg)
	}

	signedTx, err := types.SignTx(
		tx,
		types.NewLondonSigner(chainID),
		client.Account.AccountKey,
	)
	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [MintNFT] | Problem with signing tx: %v\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, err,
		)
		return false, errors.New(msg)
	}

	ctx := context.Background()

	err = client.Rpc.SendTransaction(ctx, signedTx)
	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [MintNFT] | Problem with signing tx: %v\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, err,
		)
		return false, errors.New(msg)
	} else {
		msg := fmt.Sprintf("[%d/%d] | %s | [MintNFT] | Successfully executed transaction\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
		)
		log.Info(msg)
	}

	return true, nil
}
