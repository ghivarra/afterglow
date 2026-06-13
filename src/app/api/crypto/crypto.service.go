package crypto

import (
	"errors"
	"ghivarra/afterglow/environment"
	"ghivarra/afterglow/src/service/encryptservice"
)

func decryptText(payload DecryptRequestDto) (DecryptResponseDto, error) {
	if environment.APP_ENV != "development" {
		return DecryptResponseDto{}, errors.New("crypto decrypt endpoint is only active on development environment")
	}

	text, err := encryptservice.Decrypt(payload.Text)
	if err != nil {
		return DecryptResponseDto{}, err
	}

	return DecryptResponseDto{
		Text: text,
	}, nil
}
