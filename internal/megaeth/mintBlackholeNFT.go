package megaeth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"

	"main/internal"
	"main/pkg/global"
	accTypes "main/pkg/types"
	"main/pkg/utils"
)

type DataRaw struct {
	Source string `json:"source"`
	Action string `json:"action"`
	ChainID *big.Int `json:"chainId"`
	ClientID string `json:"clientId"`
	ContractAddress string `json:"contractAddress"`
	TransactionHash string `json:"transactionHash"`
	WalletAddress common.Address `json:"walletAddress"`
	WalletType string `json:"walletType"`
}

func verifyTransaction(
	accountData *accTypes.AccountData,
	chainID **big.Int,
	contractAddress *string,
	txHash *string,
) bool {
	var url = "https://c.thirdweb.com/event"
	clientHTTP := utils.GetClient()

	if clientHTTP == nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [verifyTransaction] | Failed to get http client\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
		)
		log.Error(msg)
		return false
	}

	reqBody := DataRaw{
		Source: "sdk",
		Action: "transaction:sent",
		ChainID: *chainID,
		ClientID: "154af4b042b6a335e64ef7636462b86d",
		ContractAddress: *contractAddress,
		TransactionHash: *txHash,
		WalletAddress: accountData.AccountAddress,
		WalletType: "io.rabby",
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [verifyTransaction] | Failed to marshal request payload: %v\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, err,
		)
		log.Error(msg)
		return false
	}

	req := fasthttp.AcquireRequest()

	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.SetBody(bodyBytes)

	req.Header.Set("accept", "*/*")
	req.Header.Set("authoruty", "c.thirdweb.com")
	req.Header.Set("accept-language", "ru,en-US;q=0.9,en;q=0.8")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("content-type", "text/plain;charset=UTF-8")
	req.Header.Set("origin", "https://morkie.xyz")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("referer", "https://morkie.xyz/")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	req.Header.Set("x-client-id", "154af4b042b6a335e64ef7636462b86d")
	req.Header.Set("x-sdk-name", "unified-sdk")
	req.Header.Set("x-sdk-os", "win")
	req.Header.Set("x-sdk-platform", "browser")
	req.Header.Set("x-sdk-version", "5.105.16")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	if err := clientHTTP.Do(req, resp); err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [verifyTransaction] | Failed to do request: %v\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, err,
		)
		log.Error(msg)
		return false
	}

	respStatus := resp.StatusCode()

	if respStatus != 200 {
		msg := fmt.Sprintf("[%d/%d] | %s | [verifyTransaction] | Wrong Response Status Code: %v\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, respStatus,
		)
		log.Error(msg)
		return false
	}

	respBody := resp.Body()
	json := string(respBody)
	message := gjson.Get(json, "message").String()

	return message == "OK"
}

func MintBlackholeNFT(accountData accTypes.AccountData) (bool, error) {
	log.Infof("[%d/%d] | %s | [MintBlackholeNFT] | Start Minting Blackhole NFT ...\n",
		global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
	)

	var contractAddress = "0xcfD3dDe3A4B393a2a204ff16B112C2cA9B85abb7"

	quantity := big.NewInt(1)
	currency := common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE")
	pricePerToken := big.NewInt(0)

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
		msg := fmt.Sprintf("[%d/%d] | %s | [MintBlackholeNFT] | Problem with encoding parameters\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
		)
		log.Error(msg)
		return false, errors.New(msg)
	}

	client, ok := internal.GetClient(&accountData)
	if !ok {
		msg := fmt.Sprintf("[%d/%d] | %s | [MintBlackholeNFT] | Problem with client initialization\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
		)
		log.Error(msg)
		return false, errors.New(msg)
	}

	tx, err := client.BuildTransaction(
		contractAddress,
		data,
		big.NewInt(0),
	)

	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [MintBlackholeNFT] | Problem with building transaction: %v\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, err,
		)
		log.Error(msg)
		return false, errors.New(msg)
	}

	chainID, err := client.GetChainID()
	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [MintBlackholeNFT] | Problem with getting chainID: %v\n",
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
		msg := fmt.Sprintf("[%d/%d] | %s | [MintBlackholeNFT] | Problem with signing tx: %v\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, err,
		)
		log.Error(msg)
		return false, errors.New(msg)
	}

	ctx := context.Background()

	err = client.Rpc.SendTransaction(ctx, signedTx)
	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [MintBlackholeNFT] | Problem with signing tx: %v\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, err,
		)
		log.Error(err)
		return false, errors.New(msg)
	} else {
		msg := fmt.Sprintf("[%d/%d] | %s | [MintBlackholeNFT] | Successfully executed transaction\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
		)
		log.Info(msg)
	}

	receipt, err := bind.WaitMined(ctx, client.Rpc, signedTx)
	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [MintBlackholeNFT] | Problem while waiting mining transaction: %v\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, err,
		)
		log.Error(msg)
		return false, errors.New(msg)
	}

	receiptHash := receipt.TxHash.Hex()
	if receipt.Status == types.ReceiptStatusSuccessful {
		msg := fmt.Sprintf("[%d/%d] | %s | [MintBlackholeNFT] | Transaction successfully was mined | Tx Hash: %v\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, receiptHash,
		)
		log.Info(msg)
	} else {
		msg := fmt.Sprintf("[%d/%d] | %s | [MintBlackholeNFT] | Transaction was reverted: %v\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, receiptHash,
		)
		log.Error(msg)
	}

	// verify transaction
	result := verifyTransaction(&accountData, &chainID, &contractAddress, &receiptHash)

	// verification failed but tx hash exists -> transaction was made
	if !result {
		msg := fmt.Sprintf("[%d/%d] | %s | [MintBlackholeNFT] | Problem with transaction verification: %v\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, receiptHash,
		)
		log.Error(msg)
		return true, errors.New(msg)
	}

	return true, nil

}
