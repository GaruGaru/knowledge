package data

import "context"

type Catalog interface {
	Init() error
	InsertDocument(context.Context, InsertDocumentRequest) error
	GetDocument(context.Context, GetDocumentRequest) (Document, error)
	ListDocuments(context.Context, ListDocumentsRequest) (ListDocumentsResponse, error)
}

type InsertDocumentRequest struct {
	Document Document
}

type ListDocumentsRequest struct {
	Title      string
	Tags       []string
	Pagination PaginationRequest
}

type ListDocumentsResponse struct {
	Items      []Document         `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}

type PaginationRequest struct {
	Page     int
	PageSize int
}

func (p PaginationRequest) Offset() int {
	return (p.Page - 1) * p.PageSize
}

type GetDocumentRequest struct {
	DocumentID int
}

type PaginationResponse struct {
	TotalElements int64 `json:"total_elements"`
	Page          int   `json:"page"`
	Pages         int   `json:"pages"`
}
