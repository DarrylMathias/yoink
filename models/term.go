package models

type Term struct {
    Id   int64  `json:"id" gorm:"primaryKey;autoIncrement"`
    Word string `json:"word" gorm:"uniqueIndex;not null; index"`
    DF   int32  `json:"df"`
}