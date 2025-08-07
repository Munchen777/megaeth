package megaeth

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"

	"main/pkg/global"
	"main/pkg/types"
	"main/pkg/utils"
)

func checkBalance(accountData types.AccountData) bool {
	var err error
	var result bool
	var url = "https://api.capmonster.cloud/getBalance"

	for i:= 1; i < 6; i++ {
		payload := map[string]interface{}{
			"clientKey": global.Config.CapmonsterAPIKey,
		}
		payloadBytes, _ := json.Marshal(payload)

		client := utils.GetClient()

		req := fasthttp.AcquireRequest()

		defer fasthttp.ReleaseRequest(req)

		req.SetRequestURI(url)
		req.Header.SetMethod("POST")
		req.SetBody(payloadBytes)

		resp := fasthttp.AcquireResponse()

		defer fasthttp.ReleaseResponse(resp)

		if err = client.Do(req, resp); err != nil {
			log.Warnf("[%d/%d] | %s | [checkBalance] | Failed to do request: %v\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, err,
			)
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
			continue
		}

		respStatus := resp.StatusCode()
		respBody := resp.Body()
		json := string(respBody)

		errorId := gjson.Get(json, "errorId").Int()

		if respStatus != fasthttp.StatusOK {
			log.Warnf("[%d/%d] | %s | [checkBalance] | Wrong Response Status Code: %v\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, respStatus,
			)
			continue
		}

		if errorId == 0 {
			balance := gjson.Get(json, "balance").Float()
			if balance <= 0.01 {
				log.Warnf("[%d/%d] | %s | [checkBalance] | Insufficient balance on Capmonster: %v\n",
					global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, err,
				)
				result = false

			} else {
				log.Printf("[%d/%d] | %s | [checkBalance] | You have enough balance %f (>0.01)\n",
					global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, balance,
				)
				result = true
				break
			}

		} else {
			log.Warnf("[%d/%d] | %s | [checkBalance] | External error while checking balance: %v\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, err,
			)
			result = false
		}
	}
	return result
}

func createTask(accountData types.AccountData) (int, error) {
	var err error
	var errorCode int = -1
	var url = "https://api.capmonster.cloud/createTask"

	for i:= 1; i < 6; i++ {
		payload := map[string]interface{}{
			"clientKey": global.Config.CapmonsterAPIKey,
			"task": map[string]interface{}{
				"type": "TurnstileTask",
				"websiteURL": "https://testnet.megaeth.com/#1",
				"websiteKey": "0x4AAAAAABA4JXCaw9E2Py-9",
				"userAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36",
			},
		}
		payloadBytes, _ := json.Marshal(payload)

		client := utils.GetClient()

		req := fasthttp.AcquireRequest()

		defer fasthttp.ReleaseRequest(req)

		req.SetRequestURI(url)
		req.Header.SetMethod("POST")
		req.SetBody(payloadBytes)

		resp := fasthttp.AcquireResponse()

		defer fasthttp.ReleaseResponse(resp)

		if err = client.Do(req, resp); err != nil {
			log.Warnf("[%d/%d] | %s | [checkBalance] | Failed to do request: %v\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, err,
			)
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
			continue
		}

		respStatus := resp.StatusCode()
		respBody := resp.Body()
		json := string(respBody)

		errorId := gjson.Get(json, "errorId").Int()

		if respStatus != fasthttp.StatusOK {
			log.Warnf("[%d/%d] | %s | [checkBalance] | Wrong Response Status Code: %v\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, respStatus,
			)
			continue
		}

		if errorId == 0 {
			return int(gjson.Get(json, "taskId").Int()), nil
		} else {
			errorCode = int(gjson.Get(json, "errorCode").Int())
		}
	}
	return errorCode, err
}

func resultCaptcha(accountData types.AccountData, taskId int) (string, error) {
	var err error
	var url = "https://api.capmonster.cloud/getTaskResult"

	for i:= 1; i < 6; i++ {
		payload := map[string]interface{}{
			"clientKey": global.Config.CapmonsterAPIKey,
			"taskId": taskId,
		}
		payloadBytes, _ := json.Marshal(payload)

		client := utils.GetClient()

		req := fasthttp.AcquireRequest()

		defer fasthttp.ReleaseRequest(req)

		req.SetRequestURI(url)
		req.Header.SetMethod("POST")
		req.SetBody(payloadBytes)

		resp := fasthttp.AcquireResponse()

		defer fasthttp.ReleaseResponse(resp)

		if err = client.Do(req, resp); err != nil {
			log.Warnf("[%d/%d] | %s | [resultCaptcha] | Failed to do request: %v\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, err,
			)
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
			continue
		}

		respStatus := resp.StatusCode()
		respBody := resp.Body()
		json := string(respBody)

		if respStatus != fasthttp.StatusOK {
			log.Warnf("[%d/%d] | %s | [resultCaptcha] | Wrong Response Status Code: %v\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, respStatus,
			)
			continue
		}
		status := gjson.Get(json, "status").String()

		if status == "processing" {
			log.Infof("[%d/%d] | %s | [resultCaptcha] | Captcha result is still processing. Sleep 3 sec ...\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
			)
			time.Sleep(3 * time.Second)
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)

		} else if status == "ready" {
			return gjson.Get(json, "solution.token").String(), nil
		}
	}
	return "", err
}

func getTokenFaucet(accountData types.AccountData, token string) (bool, error) {
	var err error
	var url = "https://carrot.megaeth.com/claim"

	for i:= 1; i < 6; i++ {
		payload := map[string]interface{}{
			"addr": accountData.AccountAddress,
			"token": token,
		}
		payloadBytes, _ := json.Marshal(payload)

		client := utils.GetClient()

		req := fasthttp.AcquireRequest()

		defer fasthttp.ReleaseRequest(req)

		req.SetRequestURI(url)
		req.Header.SetMethod("POST")
		req.SetBody(payloadBytes)

		resp := fasthttp.AcquireResponse()

		defer fasthttp.ReleaseResponse(resp)

		if err = client.Do(req, resp); err != nil {
			log.Warnf("[%d/%d] | %s | [getTokenFaucet] | Failed to do request: %v\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, err,
			)
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
			continue
		}

		respStatus := resp.StatusCode()
		respBody := resp.Body()
		json := string(respBody)

		if respStatus != fasthttp.StatusOK {
			log.Warnf("[%d/%d] | %s | [getTokenFaucet] | Wrong Response Status Code: %v\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, respStatus,
			)
			continue
		}

		message, success := gjson.Get(json, "message").String(), gjson.Get(json, "success").Bool()
		if message == "" && success {
			msg := fmt.Sprintf("[%d/%d] | %s | [getTokenFaucet] | Successfully get tokens from Faucet\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
			)
			log.Info(msg)
			return true, nil

		} else if !success {
			msg := fmt.Sprintf("[%d/%d] | %s | [getTokenFaucet] | %v\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, message,
			)
			return false, errors.New(msg)
		}
	}
	return false, err
}

func FaucetTokens(accountData types.AccountData) (bool, error) {
	if global.Config.CapmonsterAPIKey == "" {
		msg := fmt.Sprintf("[%d/%d] | %s | [FaucetTokens] | You miss Capmonster API Key\n",
		global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
	)
	return false, errors.New(msg)
}

	// check balance
	result := checkBalance(accountData)

	if !result {
		msg := fmt.Sprintf("[%d/%d] | %s | [FaucetTokens] | Not enough money for Capmonster captcha service\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
		)
		return false, errors.New(msg)
	}

	// create task for solving captcha
	taskId, err := createTask(accountData)

	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [FaucetTokens] | Problem with creating task for captcha solving\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
		)
		return false, errors.New(msg)
	}

	// try to resolve captcha
	token, err := resultCaptcha(accountData, taskId)

	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [FaucetTokens] | Problem with getting captcha token\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
		)
		return false, errors.New(msg)
	}

	_, err = getTokenFaucet(accountData, token)

	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [FaucetTokens] | Problem with getting test tokens from Faucet\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
		)
		return false, errors.New(msg)
	}

	return true, nil
}
