package server

import (
	"ghivarra/afterglow/src/mapping/entity/dbentity"
	"ghivarra/afterglow/src/repository/db/afterglowdb/servertable"
)

func createServer(payload ServerCreateRequestDto) (dbentity.ServerEntity, error) {
	return servertable.CreateServer(dbentity.ServerEntity{
		Id:          payload.Id,
		Alias:       payload.Alias,
		Name:        payload.Name,
		Description: payload.Description,
		IpAddress:   payload.IpAddress,
		AccountId:   payload.AccountId,
	})
}

func updateServer(id int, payload ServerUpdateRequestDto) (*dbentity.ServerEntity, error) {
	data := map[string]any{}

	if payload.Alias != nil {
		data["alias"] = *payload.Alias
	}

	if payload.Name != nil {
		data["name"] = *payload.Name
	}

	if payload.Description != nil {
		data["description"] = *payload.Description
	}

	if payload.IpAddress != nil {
		data["ip_address"] = *payload.IpAddress
	}

	if len(data) == 0 {
		return servertable.FetchById(id)
	}

	return servertable.PartialUpdateServer(id, data)
}

func deleteServer(id int) error {
	return servertable.DeleteServer(id)
}
