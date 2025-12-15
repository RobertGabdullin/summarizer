package service

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/RobertGabdullin/summarizer/internal/models"
)

type metadataStorage interface {
	StoreMetadata(context.Context, *models.Metadata) (int, error)
}

type meetingStorage interface {
	StoreMeeting(context.Context, string, io.Reader) error
}

type UserData struct {
	Metadata metadataStorage
	Meeting  meetingStorage
}

func NewUserData(metadataStorage metadataStorage, meeting meetingStorage) *UserData {
	return &UserData{
		Metadata: metadataStorage,
		Meeting:  meeting,
	}
}

func (u *UserData) SaveData(ctx context.Context, record *models.Record) error {
	id, err := u.Metadata.StoreMetadata(ctx, record.Metadata)
	if err != nil {
		return fmt.Errorf("Не удалось сохранить метаданные: %w", err)
	}

	idStr := strconv.Itoa(id)

	err = u.Meeting.StoreMeeting(ctx, idStr, record.Meeting)
	if err != nil {
		return fmt.Errorf("Не удалось сохранить запись встречи: %w", err)
	}

	return nil
}
