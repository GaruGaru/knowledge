package data

import (
	"gorm.io/gorm"
)

type Document struct {
	gorm.Model

	ID             int
	Title          *string `gorm:"not null"`
	Uri            *string `gorm:"not null"`
	DocumentKindID int
	DocumentKind   DocumentKind
	Authors        []DocumentAuthor `gorm:"many2many:document_document_authors;"`
	Tags           []DocumentTag    `gorm:"many2many:document_document_tags;"`
	CreateTime     int              `gorm:"autoCreateTime"`
}

type DocumentKind struct {
	gorm.Model
	Name string
}

type DocumentTag struct {
	gorm.Model
	Tag string
}

type DocumentAuthor struct {
	gorm.Model
	Name    string
	Surname string
}
