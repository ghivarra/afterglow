package account

import (
	"fmt"
	"ghivarra/afterglow/src/mapping/entity/dbentity"
	"ghivarra/afterglow/src/repository/db/afterglowdb/accounttable"
	"ghivarra/afterglow/src/service/encryptservice"

	"github.com/google/uuid"
)

func createAccount(payload AccountCreateRequestDto) (AccountCreateResponseDto, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return AccountCreateResponseDto{}, fmt.Errorf("failed to generate uuid. Error: %v", err)
	}

	password, err := encryptservice.Encrypt(payload.Password)
	if err != nil {
		return AccountCreateResponseDto{}, fmt.Errorf("failed to encrypt account password credential. Error: %v", err)
	}

	apiClientId, err := encryptservice.Encrypt(payload.ApiClientId)
	if err != nil {
		return AccountCreateResponseDto{}, fmt.Errorf("failed to encrypt account api client id credential. Error: %v", err)
	}

	apiClientKey, err := encryptservice.Encrypt(payload.ApiClientKey)
	if err != nil {
		return AccountCreateResponseDto{}, fmt.Errorf("failed to encrypt account api client key credential. Error: %v", err)
	}

	entity, err := accounttable.CreateAccount(dbentity.AccountEntity{
		Id:           uuid.String(),
		Username:     payload.Username,
		Password:     password,
		ApiClientId:  apiClientId,
		ApiClientKey: apiClientKey,
		IsActive:     1,
	})
	if err != nil {
		return AccountCreateResponseDto{}, err
	}

	return AccountCreateResponseDto{
		Id:       entity.Id,
		Username: entity.Username,
	}, nil
}
