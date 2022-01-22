package data

import (
	"gorm.io/gorm"
)

type Document struct {
	gorm.Model

	ID             int              `json:"ID,omitempty"`
	Title          *string          `gorm:"not null" json:"title,omitempty"`
	Uri            *string          `gorm:"not null" json:"uri,omitempty"`
	DocumentKindID int              `json:"documentKindID,omitempty"`
	DocumentKind   DocumentKind     `json:"documentKind"`
	Authors        []DocumentAuthor `gorm:"many2many:document_document_authors;" json:"authors,omitempty"`
	Tags           []DocumentTag    `gorm:"many2many:document_document_tags;" json:"tags,omitempty"`
	CreateTime     int              `gorm:"autoCreateTime" json:"createTime,omitempty"`
}

type DocumentKind struct {
	gorm.Model
	Name string `json:"name,omitempty"`
}

type DocumentTag struct {
	gorm.Model
	Tag string `json:"tag,omitempty"`
}

type DocumentAuthor struct {
	gorm.Model
	Name    string `json:"name,omitempty"`
	Surname string `json:"surname,omitempty"`
}
