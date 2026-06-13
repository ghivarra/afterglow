package contaboclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"ghivarra/afterglow/environment"
	"ghivarra/afterglow/src/mapping/dto"
	"ghivarra/afterglow/src/mapping/entity/dbentity"
	"ghivarra/afterglow/src/service/curlservice"
	"ghivarra/afterglow/src/service/encryptservice"
	"ghivarra/afterglow/src/service/logservice"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
)

const CREATE_SNAPSHOT_PATH = "/compute/instances/{{SERVER_ID}}/snapshots"
const CREATE_SNAPSHOT_TITLE = "AUTO BACKUP {{DATETIME}}"
const CREATE_SNAPSHOT_DESC = "(automated using {{APP_NAME}})"

func CreateSnapshotBackup(ctx context.Context, serverId int, encryptedAccessToken string) dto.ContaboCreateSnapshotResult {
	var err error
	var errMessage string

	// generate request id
	uuid := uuid.New()
	requestId := uuid.String()

	// decrypt token
	accessToken, err := encryptservice.Decrypt(encryptedAccessToken)
	if err != nil {
		errMessage = fmt.Sprintf("failed to decrypt account access token. Error: %v", err)
		log.Errorf(errMessage)

		return dto.ContaboCreateSnapshotResult{
			ResultStatus: false,
			Message:      errMessage,
			Error:        err,
		}
	}

	// build url and payload
	now := time.Now().UTC()
	datetime := now.Format("2006-01-02 15-04")
	path := strings.Replace(CREATE_SNAPSHOT_PATH, "{{SERVER_ID}}", strconv.Itoa(serverId), 1)
	snapshotUrl := environment.API_CONTABO_GENERAL_HOST + path
	payload := dto.ContaboCreateSnapshotPayload{
		Name:        strings.Replace(CREATE_SNAPSHOT_TITLE, "{{DATETIME}}", datetime, 1),
		Description: strings.Replace(CREATE_SNAPSHOT_DESC, "{{APP_NAME}}", environment.APP_NAME, 1),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		errMessage = fmt.Sprintf("failed to json stringify curl payloads. Error: %v", err)
		log.Errorf(errMessage)

		return dto.ContaboCreateSnapshotResult{
			ResultStatus: false,
			Message:      errMessage,
			Error:        err,
		}
	}
	payloadStr := string(payloadBytes)

	// create client
	ctx, cancel := context.WithTimeout(
		ctx,
		TIMEOUT_TIME*time.Second,
	)
	defer cancel()

	req, err := curlservice.CreateNewRequest(ctx, "POST", snapshotUrl, bytes.NewReader(payloadBytes))
	if err != nil {
		return dto.ContaboCreateSnapshotResult{
			ResultStatus: false,
			Message:      err.Error(),
			Error:        err,
		}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-request-id", requestId)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// build client
	requestTime := time.Now().UTC().Format(time.RFC3339)
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
			Url:          snapshotUrl,
			Payload:      payloadStr,
			ResponseCode: 425,
			RequestedAt:  requestTime,
		})

		return dto.ContaboCreateSnapshotResult{
			ResultStatus: false,
			Message:      errMessage,
			Error:        err,
		}
	}
	defer response.Body.Close()

	// response time
	responseTime := time.Now().UTC().Format(time.RFC3339)

	// parse body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		errMessage = fmt.Sprintf("failed to read response body. Error: %v", err)
		log.Errorf(errMessage)

		logservice.StoreContaboLog(dbentity.LogEntity{
			Id:           requestId,
			Url:          snapshotUrl,
			Payload:      payloadStr,
			ResponseCode: 425,
			RequestedAt:  requestTime,
		})

		return dto.ContaboCreateSnapshotResult{
			ResultStatus: false,
			Message:      errMessage,
			Error:        err,
		}
	}

	// parse
	var result dto.ContaboCreateSnapshotResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		errMessage = fmt.Sprintf("failed to parse response body into struct. Error: %v", err)
		log.Errorf(errMessage)

		logservice.StoreContaboLog(dbentity.LogEntity{
			Id:           requestId,
			Url:          snapshotUrl,
			Payload:      payloadStr,
			ResponseBody: new(string(body)),
			ResponseCode: response.StatusCode,
			RequestedAt:  requestTime,
			RespondedAt:  &responseTime,
		})

		return dto.ContaboCreateSnapshotResult{
			ResultStatus: false,
			Message:      errMessage,
			Error:        err,
		}
	}

	// check status
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		logservice.StoreContaboLog(dbentity.LogEntity{
			Id:           requestId,
			Url:          snapshotUrl,
			Payload:      payloadStr,
			ResponseBody: new(string(body)),
			ResponseCode: response.StatusCode,
			RequestedAt:  requestTime,
			RespondedAt:  &responseTime,
		})

		return dto.ContaboCreateSnapshotResult{
			ResultStatus: false,
			Message:      fmt.Sprintf("error %d", response.StatusCode),
			Error:        fmt.Errorf("error %d", response.StatusCode),
		}
	}

	logservice.StoreContaboLog(dbentity.LogEntity{
		Id:           requestId,
		Url:          snapshotUrl,
		Payload:      payloadStr,
		ResponseBody: new(string(body)),
		ResponseCode: response.StatusCode,
		RequestedAt:  requestTime,
		RespondedAt:  &responseTime,
	})

	return dto.ContaboCreateSnapshotResult{
		ResultStatus: true,
		Result:       &result,
		Message:      "OK",
		RawData:      new(string(body)),
		Error:        nil,
	}
}
