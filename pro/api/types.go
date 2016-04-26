package api

import (
    "github.com/jinzhu/gorm"
)

type Repository struct {
    gorm.Model
    RepoName string
    Group    string
    Tags     []Tag
    TPullCount uint
}

type Tag struct {
    gorm.Model
    RepositoryID uint `sql:"index"`
    Name         string
    Description  string `sql:"type:varchar(255)"`
    Digest       string
    PushTime     string
    PullCount    uint
    Size         int64
}

