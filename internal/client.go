package internal

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"

	"main/pkg/types"
	"main/pkg/global"
)

type Client struct {
	Rpc *ethclient.Client
	Account *types.AccountData
}

func (c *Client) GetNonce() (uint64, error) {
	ctx := context.Background()
	return c.Rpc.PendingNonceAt(ctx, c.Account.AccountAddress)
}

func (c *Client) GetBalance() (*big.Int, error) {
	ctx := context.Background()
	return c.Rpc.BalanceAt(ctx, c.Account.AccountAddress, nil)
}

func (c *Client) GetChainID() (*big.Int, error) {
	ctx := context.Background()
	return c.Rpc.ChainID(ctx)
}

func (c *Client) BuildTransaction(
	to string,
	data []byte,
	value *big.Int,
) (*ethTypes.Transaction, error) {
	ctx := context.Background()
	nonce, err := c.GetNonce()

	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [BuildTransaction] | Problem with getting nonce\n",
			global.CurrentProgress, global.TargetProgress, c.Account.AccountAddress,
		)
		log.Error(msg)
		return nil, errors.New(msg)
	}
	
	chainID, err := c.GetChainID()

	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [BuildTransaction] | Problem with getting chainID\n",
			global.CurrentProgress, global.TargetProgress, c.Account.AccountAddress,
		)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	toAddress := common.HexToAddress(to)
	msg := ethereum.CallMsg{
		From:  c.Account.AccountAddress,
		To:    &toAddress,
		Value: value,
		Data:  data,
	}

	header, err := c.Rpc.HeaderByNumber(ctx, nil)
	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [BuildTransaction] | Problem with suggesting priority fee\n",
			global.CurrentProgress, global.TargetProgress, c.Account.AccountAddress,
		)
		log.Error(msg)
		return nil, errors.New(msg)
	}
	baseFee := header.BaseFee

	priorityFee, err := c.Rpc.SuggestGasTipCap(ctx)
	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [BuildTransaction] | Problem with suggesting priority fee\n",
			global.CurrentProgress, global.TargetProgress, c.Account.AccountAddress,
		)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	maxFee := new(big.Int).Add(baseFee, priorityFee)

	gasLimit, err := c.Rpc.EstimateGas(ctx, msg)
	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [BuildTransaction] | Problem with estimating gas\n",
			global.CurrentProgress, global.TargetProgress, c.Account.AccountAddress,
		)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	tx := ethTypes.NewTx(&ethTypes.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: maxFee,
		GasTipCap: priorityFee,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     value,
		Data:      data,
	})

	return tx, nil
}

func GetClient(accountData *types.AccountData) (*Client, bool) {
	ctx := context.Background()
	rpc, err := ethclient.DialContext(ctx, "https://carrot.megaeth.com/rpc")

	if err != nil {
		return nil, false
	}

	return &Client{Rpc: rpc, Account: accountData}, true
}
