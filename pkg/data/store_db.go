package data

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"math"
)

type DBStore struct {
	db *gorm.DB
}

func NewDBStore(db *gorm.DB) *DBStore {
	return &DBStore{db: db}
}

func (d *DBStore) Init() error {
	return d.db.AutoMigrate(&Document{}, &DocumentKind{}, &DocumentTag{}, &DocumentAuthor{})
}

func (d *DBStore) InsertDocument(ctx context.Context, req InsertDocumentRequest) error {
	result := d.db.WithContext(ctx).Create(&req.Document)
	return result.Error
}

func (d *DBStore) ListDocuments(ctx context.Context, request ListDocumentsRequest) (ListDocumentsResponse, error) {
	query := d.db.WithContext(ctx)

	query = query.Preload("Tags").Preload("Authors")

	if len(request.Title) != 0 {
		query = query.Where("title LIKE ?", fmt.Sprintf("%%%s%%", request.Title))
	}

	if len(request.Tags) != 0 {
		query = query.Joins("LEFT JOIN document_document_tags on document_document_tags.document_id = documents.id")
		query = query.Joins("LEFT JOIN document_tags on document_tags.id = document_document_tags.document_tag_id", request.Tags)
		query = query.Where("document_tags.tag IN ?", request.Tags)
	}

	var totalElements int64
	cntQuery := query.Model(&Document{}).Count(&totalElements)
	if err := cntQuery.Error; err != nil {
		return ListDocumentsResponse{}, err
	}

	query = query.Order("title")
	query = query.Offset(request.Pagination.Offset()).Limit(request.Pagination.PageSize)

	var documents []Document
	query = query.Find(&documents)

	var response ListDocumentsResponse

	response.Items = documents
	response.Pagination = PaginationResponse{
		TotalElements: totalElements,
		Page:          request.Pagination.Page,
		Pages:         int(math.Ceil(float64(totalElements) / float64(request.Pagination.PageSize))),
	}

	return response, query.Error
}

func (d *DBStore) GetDocument(ctx context.Context, request GetDocumentRequest) (Document, error) {
	var document Document
	query := d.db.WithContext(ctx).Preload("Tags").Preload("Authors").First(&document, request.DocumentID)
	return document, query.Error
}
