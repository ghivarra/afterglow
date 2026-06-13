package contaboclient

import (
	"bytes"
	"context"
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

const DELETE_SNAPSHOT_PATH = "/compute/instances/{{SERVER_ID}}/snapshots/{{SNAPSHOT_ID}}"

func DeleteSnapshotBackup(ctx context.Context, serverId int, snapshotId string, encryptedAccessToken string) dto.ContaboDeleteSnapshotResult {
	var err error
	var errMessage string

	// generate request id
	uuid, err := uuid.NewV7()
	if err != nil {
		errMessage = fmt.Sprintf("failed to generate uuid. Error: %v", err)
		log.Errorf(errMessage)

		return dto.ContaboDeleteSnapshotResult{
			ResultStatus: false,
			Message:      errMessage,
			Error:        err,
		}
	}
	requestId := uuid.String()

	// decrypt token
	accessToken, err := encryptservice.Decrypt(encryptedAccessToken)
	if err != nil {
		errMessage = fmt.Sprintf("failed to decrypt account access token. Error: %v", err)
		log.Errorf(errMessage)

		return dto.ContaboDeleteSnapshotResult{
			ResultStatus: false,
			Message:      errMessage,
			Error:        err,
		}
	}

	// build url and payload
	path := strings.Replace(DELETE_SNAPSHOT_PATH, "{{SERVER_ID}}", strconv.Itoa(serverId), 1)
	path = strings.Replace(path, "{{SNAPSHOT_ID}}", snapshotId, 1)
	snapshotUrl := environment.API_CONTABO_HOST + path
	payloadBytes := []byte("{}")
	payloadStr := string(payloadBytes)

	// create client
	ctx, cancel := context.WithTimeout(
		ctx,
		TIMEOUT_TIME*time.Second,
	)
	defer cancel()

	req, err := curlservice.CreateNewRequest(ctx, "DELETE", snapshotUrl, bytes.NewReader(payloadBytes))
	if err != nil {
		return dto.ContaboDeleteSnapshotResult{
			ResultStatus: false,
			Message:      err.Error(),
			Error:        err,
		}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-request-id", requestId)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

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
			Url:          snapshotUrl,
			Payload:      payloadStr,
			ResponseCode: 425,
			RequestedAt:  requestTime,
		})

		return dto.ContaboDeleteSnapshotResult{
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
			Url:          snapshotUrl,
			Payload:      payloadStr,
			ResponseCode: 425,
			RequestedAt:  requestTime,
		})

		return dto.ContaboDeleteSnapshotResult{
			ResultStatus: false,
			Message:      errMessage,
			Error:        err,
		}
	}

	// check status
	if response.StatusCode != http.StatusNoContent {
		logservice.StoreContaboLog(dbentity.LogEntity{
			Id:           requestId,
			Url:          snapshotUrl,
			Payload:      payloadStr,
			ResponseBody: new(string(body)),
			ResponseCode: response.StatusCode,
			RequestedAt:  requestTime,
			RespondedAt:  &responseTime,
		})

		return dto.ContaboDeleteSnapshotResult{
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

	return dto.ContaboDeleteSnapshotResult{
		ResultStatus: true,
		Message:      "OK",
		Error:        nil,
	}
}
