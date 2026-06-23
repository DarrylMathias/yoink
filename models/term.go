package models

type Term struct {
    Id   int64  `json:"id" gorm:"primaryKey;autoIncrement"`
    Word string `json:"word" gorm:"uniqueIndex;not null; index"`
    // no of documents in which term appears
    DF   int32  `json:"df"`
}