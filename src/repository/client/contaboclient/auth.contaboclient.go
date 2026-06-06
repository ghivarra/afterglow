package contaboclient

import (
	"context"
	"encoding/json"
	"fmt"
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

const TIMEOUT_TIME = 20
const LOGIN_URL = "https://auth.contabo.com/auth/realms/contabo/protocol/openid-connect/token"

func auth(ctx context.Context, id string) dto.ContaboAuthResult {
	var err error
	var errMessage string

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
	account, err := accounttable.FetchById(id)
	if err != nil || account == nil {
		errMessage = fmt.Sprintf("account credentials is not found. Error: %v", err)
		log.Errorf(errMessage)

		return dto.ContaboAuthResult{
			ResultStatus: false,
			Message:      errMessage,
			Error:        err,
		}
	}

	// decrypt password
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

	// payloads
	payloads, err := json.Marshal(url.Values{
		"client_id":     {account.ApiClientId},
		"client_secret": {account.ApiClientKey},
		"grant_type":    {"password"},
		"username":      {account.Username},
		"password":      {password},
	})
	if err != nil {
		errMessage = fmt.Sprintf("failed to json stringify curl payloads. Error: %v", err)
		log.Errorf(errMessage)

		return dto.ContaboAuthResult{
			ResultStatus: false,
			Message:      errMessage,
			Error:        err,
		}
	}

	// create client
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	req, err := curlservice.CreateNewRequest(ctx, "POST", LOGIN_URL, strings.NewReader(string(payloads)))

	// build client
	requestTime := time.Now().UTC()
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// response
	response, err := client.Do(req)
	if err != nil {
		errMessage = fmt.Sprintf("curl request failed. Error: %v", err)
		log.Errorf(errMessage)

		logservice.StoreContaboLog(dbentity.LogEntity{
			Id:           requestId,
			Url:          LOGIN_URL,
			Payload:      string(payloads),
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
			Payload:      string(payloads),
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

	// validate
	status, err := curlservice.ValidateResponseBody(response, result)
	if err != nil {
		errMessage = fmt.Sprintf("failed to parse response body into struct. Error: %v", err)
		log.Errorf(errMessage)

		logservice.StoreContaboLog(dbentity.LogEntity{
			Id:           requestId,
			Url:          LOGIN_URL,
			Payload:      string(payloads),
			ResponseBody: new(string(body)),
			ResponseCode: response.StatusCode,
			RequestedAt:  requestTime,
			RespondedAt:  &responseTime,
		})
		return dto.ContaboAuthResult{
			ResultStatus: false,
			Message:      errMessage,
			Error:        fmt.Errorf("error %d", status),
		}
	}

	// parse
	json.Unmarshal(body, &result)

	// check status
	if status != http.StatusOK {
		logservice.StoreContaboLog(dbentity.LogEntity{
			Id:           requestId,
			Url:          LOGIN_URL,
			Payload:      string(payloads),
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
		Result:       &result,
		Token:        &result.AccessToken,
		Message:      *new("OK"),
		RawData:      new(string(body)),
		Error:        nil,
	}
}
