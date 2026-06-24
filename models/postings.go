package models

import "github.com/google/uuid"

type Posting struct {
    PageId uuid.UUID
    TF     int32
}