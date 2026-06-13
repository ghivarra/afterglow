package snapshot

import (
	"context"
	"errors"
	"fmt"
	"ghivarra/afterglow/src/mapping/entity/dbentity"
	"ghivarra/afterglow/src/repository/client/contaboclient"
	"ghivarra/afterglow/src/repository/db/afterglowdb/accounttokentable"
	"ghivarra/afterglow/src/repository/db/afterglowdb/servertable"
	"ghivarra/afterglow/src/repository/db/afterglowdb/snapshottable"
)

var errSnapshotResponseEmpty = errors.New("snapshot response is empty")

func createSnapshot(payload SnapshotCreateRequestDto) (dbentity.SnapshotEntity, error) {
	server, err := servertable.FetchByAlias(payload.ServerAlias)
	if err != nil {
		return dbentity.SnapshotEntity{}, err
	}

	accountToken, err := accounttokentable.FetchByAccountId(server.AccountId)
	if err != nil {
		return dbentity.SnapshotEntity{}, err
	}

	activeSnapshots, err := snapshottable.FetchActiveSnapshotsByServerId(server.Id)
	if err != nil {
		return dbentity.SnapshotEntity{}, err
	}

	for _, activeSnapshot := range activeSnapshots {
		deleteResult := contaboclient.DeleteSnapshotBackup(
			context.Background(),
			server.Id,
			activeSnapshot.Id,
			accountToken.AccessToken,
		)
		if !deleteResult.ResultStatus {
			if deleteResult.Error != nil {
				return dbentity.SnapshotEntity{}, fmt.Errorf("failed to delete snapshot backup %s. Error: %v", activeSnapshot.Id, deleteResult.Error)
			}

			return dbentity.SnapshotEntity{}, fmt.Errorf("failed to delete snapshot backup %s", activeSnapshot.Id)
		}

		err = snapshottable.SoftDeleteSnapshotById(activeSnapshot.Id)
		if err != nil {
			return dbentity.SnapshotEntity{}, err
		}
	}

	snapshotResult := contaboclient.CreateSnapshotBackup(context.Background(), server.Id, accountToken.AccessToken)
	if !snapshotResult.ResultStatus {
		return dbentity.SnapshotEntity{}, snapshotResult.Error
	}

	if snapshotResult.Result == nil || len(snapshotResult.Result.Data) == 0 {
		return dbentity.SnapshotEntity{}, errSnapshotResponseEmpty
	}

	contaboSnapshot := snapshotResult.Result.Data[0]

	description := contaboSnapshot.Description
	entity, err := snapshottable.CreateSnapshot(dbentity.SnapshotEntity{
		Id:          contaboSnapshot.SnapshotID,
		ServerId:    server.Id,
		Name:        contaboSnapshot.Name,
		Description: &description,
		CreatedAt:   contaboSnapshot.CreatedDate,
		DeletedAt:   nil,
	})
	if err != nil {
		return dbentity.SnapshotEntity{}, err
	}

	return entity, nil
}
