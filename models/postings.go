package models

import "github.com/google/uuid"

type Posting struct {
    PageId uuid.UUID `json:"page_id" gorm:"primaryKey; index"`
    TermId int64     `json:"term_id" gorm:"primaryKey; index"`
    TF     int32     `json:"tf"`
}