package megaeth

import (
	"context"
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

func checkBalance(ctx context.Context, accountData types.AccountData) bool {
	var err error
	var result bool
	var url = "https://api.capmonster.cloud/getBalance"

	maxAttempts := 5
	retryDelay := 3 * time.Second

	payload := map[string]interface{}{
		"clientKey": global.Config.CapmonsterAPIKey,
	}

	ctx, cancel := context.WithTimeout(ctx, 30 * time.Second)
	defer cancel()

	for i := 1; i < maxAttempts + 1; i++ {
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
			log.Warnf("[%d/%d] | %s | [checkBalance] | Attempt: [%d/%d] | Failed to do request: %v\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, i, maxAttempts, err,
			)
		} else {
			respStatus := resp.StatusCode()
			respBody := resp.Body()
			json := string(respBody)
	
			errorId := gjson.Get(json, "errorId").Int()
	
			if respStatus != fasthttp.StatusOK {
				log.Warnf("[%d/%d] | %s | [checkBalance] | Attempt: [%d/%d] | Wrong Response Status Code: %v\n",
					global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, i, maxAttempts, respStatus,
				)
			}

			if errorId == 0 {
				balance := gjson.Get(json, "balance").Float()
				if balance <= 0.01 {
					log.Warnf("[%d/%d] | %s | [checkBalance] | Attempt: [%d/%d] | Insufficient balance on Capmonster: %v\n",
						global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, i, maxAttempts, balance,
					)
				} else {
					log.Printf("[%d/%d] | %s | [checkBalance] |Attempt: [%d/%d] | You have enough balance %0.2f (>0.01)\n",
						global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, i, maxAttempts, balance,
					)
					result = true
					break
				}
			} else {
				log.Warnf("[%d/%d] | %s | [checkBalance] | Attempt: [%d/%d] | External error while checking balance: %v\n",
					global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, i, maxAttempts, err,
				)
			}
		}

		select {
		case <-ctx.Done():
			log.Warnf("[%d/%d] | %s | [checkBalance] | Attempt: [%d/%d] | Cancelled (timeout reached)\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, i, maxAttempts,
			)
			return false
		case <-time.After(retryDelay):
			log.Infof("[%d/%d] | %s | [checkBalance] | Attempt: [%d/%d] | Retry after %f sec ...\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, i, maxAttempts, retryDelay.Seconds(),
			)
		}
	}
	return result
}

func createTask(ctx context.Context, accountData types.AccountData) (int, error) {
	var err error
	var errorCode int = -1
	var url = "https://api.capmonster.cloud/createTask"

	maxAttempts := 5
	retryDelay := 3 * time.Second

	ctx, cancel := context.WithTimeout(ctx, 30 * time.Second)
	defer cancel()

	payload := map[string]interface{}{
		"clientKey": global.Config.CapmonsterAPIKey,
		"task": map[string]interface{}{
			"type":       "TurnstileTask",
			"websiteURL": "https://testnet.megaeth.com/#1",
			"websiteKey": "0x4AAAAAABA4JXCaw9E2Py-9",
			"userAgent":  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36",
		},
	}

	for i := 1; i < maxAttempts + 1; i++ {
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
			log.Warnf("[%d/%d] | %s | [createTask] | Attempt: [%d/%d] | Failed to do request: %v\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, i, maxAttempts, err,
			)
		} else {
			respStatus := resp.StatusCode()
			respBody := resp.Body()
			json := string(respBody)
	
			errorId := gjson.Get(json, "errorId").Int()
	
			if respStatus != fasthttp.StatusOK {
				log.Warnf("[%d/%d] | %s | [createTask] | Attempt: [%d/%d] | Wrong Response Status Code: %v\n",
					global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, i, maxAttempts, respStatus,
				)
			}
	
			if errorId == 0 {
				return int(gjson.Get(json, "taskId").Int()), nil
			} else {
				errorCode = int(gjson.Get(json, "errorCode").Int())
			}
		}

		select {
		case <-ctx.Done():
			log.Warnf("[%d/%d] | %s | [createTask] | Attempt: [%d/%d] | Cancelled (timeout reached)\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, i, maxAttempts,
			)
			return errorCode, err
		case <-time.After(retryDelay):
			log.Infof("[%d/%d] | %s | [createTask] | Attempt: [%d/%d] | Retry after %f sec ...\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, i, maxAttempts, retryDelay.Seconds(),
			)
		}
	}
	return errorCode, err
}

func resultCaptcha(ctx context.Context, accountData types.AccountData, taskId int) (string, error) {
	var err error
	var url = "https://api.capmonster.cloud/getTaskResult"

	maxAttempts := 5
	retryDelay := 3 * time.Second

	ctx, cancel := context.WithTimeout(ctx, 30 * time.Second)
	defer cancel()

	payload := map[string]interface{}{
		"clientKey": global.Config.CapmonsterAPIKey,
		"taskId":    taskId,
	}

	for i := 0; i < maxAttempts; i++ {
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
			log.Warnf("[%d/%d] | %s | [resultCaptcha] | Attempt: [%d/%d] | Failed to do request: %v\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, i, maxAttempts, err,
			)
		} else {
			respStatus := resp.StatusCode()
			respBody := resp.Body()
			json := string(respBody)
	
			if respStatus != fasthttp.StatusOK {
				log.Warnf("[%d/%d] | %s | [resultCaptcha] | Attempt: [%d/%d] | Wrong Response Status Code: %v\n",
					global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, i, maxAttempts, respStatus,
				)
			}
	
			status := gjson.Get(json, "status").String()
	
			if status == "processing" {
				log.Infof("[%d/%d] | %s | [resultCaptcha] | Attempt: [%d/%d] | Captcha result is still processing ...\n",
					global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, i, maxAttempts,
				)
				log.Infof("[%d/%d] | %s | [resultCaptcha] | Attempt: [%d/%d] | Retry after %v seconds ...\n",
					global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, i, maxAttempts, retryDelay.Seconds(),
				)
			} else if status == "ready" {
				return gjson.Get(json, "solution.token").String(), nil
			}
		}

		select {
		case <-ctx.Done():
			log.Warnf("[%d/%d] | %s | [resultCaptcha] | Cancelled while waiting for retry\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
			)
			return "", err
		case <-time.After(retryDelay):
			log.Infof("[%d/%d] | %s | [resultCaptcha] | Attempt: [%d/%d] | Retry after %v seconds ...\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, i, maxAttempts, retryDelay.Seconds(),
			)
		}
	}
	return "", err
}

func getTokenFaucet(ctx context.Context, accountData types.AccountData, token string) (bool, error) {
	var err error
	var url = "https://carrot.megaeth.com/claim"

	maxAttempts := 5
	retryDelay := 3 * time.Second

	ctx, cancel := context.WithTimeout(ctx, 30 * time.Second)
	defer cancel()

	payload := map[string]interface{}{
		"addr":  accountData.AccountAddress,
		"token": token,
	}

	for i := 0; i < maxAttempts; i++ {
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
			log.Warnf("[%d/%d] | %s | [getTokenFaucet] | Attempt: [%d/%d] | Failed to do request: %v\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, i, maxAttempts, err,
			)
		} else {
			respStatus := resp.StatusCode()
			respBody := resp.Body()
			json := string(respBody)
	
			if respStatus != fasthttp.StatusOK {
				log.Warnf("[%d/%d] | %s | [getTokenFaucet] | Attempt: [%d/%d] | Wrong Response Status Code: %v\n",
					global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, i, maxAttempts, respStatus,
				)
			}
	
			message, success := gjson.Get(json, "message").String(), gjson.Get(json, "success").Bool()
			if message == "" && success {
				msg := fmt.Sprintf("[%d/%d] | %s | [getTokenFaucet] | Attempt: [%d/%d] | Successfully get tokens from Faucet\n",
					global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, i, maxAttempts,
				)
				log.Info(msg)
				return true, nil
	
			} else if !success {
				msg := fmt.Sprintf("[%d/%d] | %s | [getTokenFaucet] | Attempt: [%d/%d] | Failed to get tokens from Faucet: %v\n",
					global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, i, maxAttempts, message,
				)
				return false, errors.New(msg)
			}
		}

		select {
		case <-ctx.Done():
			log.Warnf("[%d/%d] | %s | [getTokenFaucet] | Attempt: [%d/%d] | Cancelled while waiting for retry\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, i, maxAttempts,
			)
			return false, err
		case <-time.After(retryDelay):
			log.Infof("[%d/%d] | %s | [getTokenFaucet] | Attempt: [%d/%d] | Retry after %v seconds ...\n",
				global.CurrentProgress, global.TargetProgress, accountData.AccountAddress, i, maxAttempts, retryDelay.Seconds(),
			)
		}
	}
	return false, err
}

func FaucetTokens(ctx context.Context, accountData types.AccountData) (bool, error) {
	if global.Config.CapmonsterAPIKey == "" {
		msg := fmt.Sprintf("[%d/%d] | %s | [FaucetTokens] | You miss Capmonster API Key\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
		)
		return false, errors.New(msg)
	}

	// check balance
	result := checkBalance(ctx, accountData)

	if !result {
		msg := fmt.Sprintf("[%d/%d] | %s | [FaucetTokens] | Not enough money for Capmonster captcha service\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
		)
		return false, errors.New(msg)
	}

	// create task for solving captcha
	taskId, err := createTask(ctx, accountData)

	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [FaucetTokens] | Problem with creating task for captcha solving\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
		)
		return false, errors.New(msg)
	}

	// try to resolve captcha
	token, err := resultCaptcha(ctx, accountData, taskId)

	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [FaucetTokens] | Problem with getting captcha token\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
		)
		return false, errors.New(msg)
	}

	_, err = getTokenFaucet(ctx, accountData, token)

	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | %s | [FaucetTokens] | Problem with getting test tokens from Faucet\n",
			global.CurrentProgress, global.TargetProgress, accountData.AccountAddress,
		)
		return false, errors.New(msg)
	}

	return true, nil
}
