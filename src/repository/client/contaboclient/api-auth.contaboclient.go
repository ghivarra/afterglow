package contaboclient

import (
	"context"
	"encoding/json"
	"fmt"
	"ghivarra/afterglow/environment"
	"ghivarra/afterglow/src/mapping/dto"
	"ghivarra/afterglow/src/mapping/entity/dbentity"
	"ghivarra/afterglow/src/repository/db/afterglowdb/accounttable"
	"ghivarra/afterglow/src/service/curlservice"
	"ghivarra/afterglow/src/service/encryptservice"
	"ghivarra/afterglow/src/service/logservice"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
)

const LOGIN_PATH = "/auth/realms/contabo/protocol/openid-connect/token"

func Authenticate(ctx context.Context, username string) dto.ContaboAuthResult {
	var err error
	var errMessage string
	var LOGIN_URL = environment.API_CONTABO_HOST + LOGIN_PATH

	// generate request id
	uuid, err := uuid.NewV7()
	if err != nil {
		errMessage = fmt.Sprintf("failed to generate uuid. Error: %v", err)
		log.Errorf(errMessage)

		return dto.ContaboAuthResult{
			ResultStatus: false,
			Message:      errMessage,
			Error:        err,
		}
	}
	requestId := uuid.String()

	// find account
	account, err := accounttable.FetchByUsername(username)
	if err != nil || account == nil {
		errMessage = fmt.Sprintf("account credentials is not found. Error: %v", err)
		log.Errorf(errMessage)

		return dto.ContaboAuthResult{
			ResultStatus: false,
			Message:      errMessage,
			Error:        err,
		}
	}

	// decrypt credentials
	password, err := encryptservice.Decrypt(account.Password)
	if err != nil {
		errMessage = fmt.Sprintf("failed to decrypt account password credential. Error: %v", err)
		log.Errorf(errMessage)

		return dto.ContaboAuthResult{
			ResultStatus: false,
			Message:      errMessage,
			Error:        err,
		}
	}

	apiClientId, err := encryptservice.Decrypt(account.ApiClientId)
	if err != nil {
		errMessage = fmt.Sprintf("failed to decrypt account api client id credential. Error: %v", err)
		log.Errorf(errMessage)

		return dto.ContaboAuthResult{
			ResultStatus: false,
			Message:      errMessage,
			Error:        err,
		}
	}

	apiClientKey, err := encryptservice.Decrypt(account.ApiClientKey)
	if err != nil {
		errMessage = fmt.Sprintf("failed to decrypt account api client key credential. Error: %v", err)
		log.Errorf(errMessage)

		return dto.ContaboAuthResult{
			ResultStatus: false,
			Message:      errMessage,
			Error:        err,
		}
	}

	// payloads
	payloads := url.Values{
		"client_id":     {apiClientId},
		"client_secret": {apiClientKey},
		"grant_type":    {"password"},
		"username":      {account.Username},
		"password":      {password},
	}
	payloadStr := payloads.Encode()

	// create client
	ctx, cancel := context.WithTimeout(
		ctx,
		TIMEOUT_TIME*time.Second,
	)
	defer cancel()

	req, err := curlservice.CreateNewRequest(ctx, "POST", LOGIN_URL, strings.NewReader(payloadStr))
	if err != nil {
		return dto.ContaboAuthResult{
			ResultStatus: false,
			Message:      err.Error(),
			Error:        err,
		}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// build client
	requestTime := time.Now().UTC()
	client := &http.Client{
		Timeout: TIMEOUT_TIME * time.Second,
	}

	// response
	response, err := client.Do(req)
	if err != nil {
		errMessage = fmt.Sprintf("curl request failed. Error: %v", err)
		log.Errorf(errMessage)

		logservice.StoreContaboLog(dbentity.LogEntity{
			Id:           requestId,
			Url:          LOGIN_URL,
			Payload:      payloadStr,
			ResponseCode: 425,
			RequestedAt:  requestTime,
		})

		return dto.ContaboAuthResult{
			ResultStatus: false,
			Message:      errMessage,
			Error:        err,
		}
	}
	defer response.Body.Close()

	// response time
	responseTime := time.Now().UTC()

	// parse body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		errMessage = fmt.Sprintf("failed to read response body. Error: %v", err)
		log.Errorf(errMessage)

		logservice.StoreContaboLog(dbentity.LogEntity{
			Id:           requestId,
			Url:          LOGIN_URL,
			Payload:      payloadStr,
			ResponseCode: 425,
			RequestedAt:  requestTime,
		})

		return dto.ContaboAuthResult{
			ResultStatus: false,
			Message:      errMessage,
			Error:        err,
		}
	}

	// create result variable
	var result dto.ContaboAuthResponse

	// parse
	err = json.Unmarshal(body, &result)
	if err != nil {
		errMessage = fmt.Sprintf("failed to parse response body into struct. Error: %v", err)
		log.Errorf(errMessage)

		logservice.StoreContaboLog(dbentity.LogEntity{
			Id:           requestId,
			Url:          LOGIN_URL,
			Payload:      payloadStr,
			ResponseBody: new(string(body)),
			ResponseCode: response.StatusCode,
			RequestedAt:  requestTime,
			RespondedAt:  &responseTime,
		})
		return dto.ContaboAuthResult{
			ResultStatus: false,
			Message:      errMessage,
			Error:        err,
		}
	}

	// check status
	if response.StatusCode != http.StatusOK {
		logservice.StoreContaboLog(dbentity.LogEntity{
			Id:           requestId,
			Url:          LOGIN_URL,
			Payload:      payloadStr,
			ResponseBody: new(string(body)),
			ResponseCode: response.StatusCode,
			RequestedAt:  requestTime,
			RespondedAt:  &responseTime,
		})

		return dto.ContaboAuthResult{
			ResultStatus: false,
			Message:      fmt.Sprintf("error %d", response.StatusCode),
			Error:        fmt.Errorf("error %d", response.StatusCode),
		}
	}

	return dto.ContaboAuthResult{
		ResultStatus: true,
		AccountId:    &account.Id,
		Result:       &result,
		Token:        &result.AccessToken,
		Message:      *new("OK"),
		RawData:      new(string(body)),
		Error:        nil,
	}
}
