package service

import (
	"context"
	"mime/multipart"

	"github.com/RobertGabdullin/summarizer/internal/models"
)

type metadataStorage interface {
	StoreMetadata(context.Context, *models.Metadata) (int, error)
}

type meetingStorage interface {
	StoreMeeting(context.Context, int, multipart.File) error
}

type requestProducer interface {
	ProduceRequest(context.Context, int) error
}

type UserData struct {
	metadata metadataStorage
	meeting  meetingStorage
	producer requestProducer
}

func (u *UserData) SaveData(ctx context.Context, record *models.Record) error {
	id, err := u.metadata.StoreMetadata(ctx, record.Metadata)
	if err != nil {
		// to do
	}

	err = u.meeting.StoreMeeting(ctx, id, record.Meeting)
	if err != nil {
		// to do
	}

	err = u.producer.ProduceRequest(ctx, id)
	if err != nil {
		// to do
	}

	return nil
}
