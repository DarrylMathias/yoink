package models

import "github.com/google/uuid"

type DocMeta struct{
	Id uuid.UUID
	DocLength int32
}