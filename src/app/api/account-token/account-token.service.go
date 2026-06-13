package accounttoken

import (
	"context"
	"errors"
	"ghivarra/afterglow/src/mapping/entity/dbentity"
	"ghivarra/afterglow/src/repository/client/contaboclient"
	"ghivarra/afterglow/src/repository/db/afterglowdb/accounttokentable"
	"ghivarra/afterglow/src/service/encryptservice"
	"time"
)

func updateToken(payload UpdateTokenRequestDto) (string, error) {
	result := contaboclient.Authenticate(context.Background(), payload.Username)
	if !result.ResultStatus {
		return "", result.Error
	}

	if result.AccountId == nil || result.Token == nil || result.Result == nil {
		return "", errors.New("contabo auth result is incomplete")
	}

	expiredAt := time.Now().UTC().Add(time.Duration(result.Result.ExpiresIn) * time.Second)
	accessToken, err := encryptservice.Encrypt(*result.Token)
	if err != nil {
		return "", err
	}

	_, err = accounttokentable.UpsertAccountToken(dbentity.AccountTokenEntity{
		AccountId:   *result.AccountId,
		AccessToken: accessToken,
		ExpiredAt:   expiredAt,
	})
	if err != nil {
		return "", err
	}

	return "OK", nil
}
