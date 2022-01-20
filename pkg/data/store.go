package data

import "context"

type Store interface {
	Init() error
	InsertDocument(context.Context, InsertDocumentRequest) error
	Search(context.Context, SearchRequest) ([]Document, error)
}

type InsertDocumentRequest struct {
	Document Document
}

type SearchRequest struct {
	Title      string
	Tags       []string
	Pagination PaginationRequest
}

type PaginationRequest struct {
	Page     int
	PageSize int
}
