package models

import (
	"mime/multipart"
)

type Record struct {
	Meeting  multipart.File
	Metadata *Metadata
}
