package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garugaru/knowledge/pkg/data"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCatalogApi_InsertDocument(t *testing.T) {
	var inserted data.Document
	catalog := mockCatalog{
		insertDocument: func(ctx context.Context, request data.InsertDocumentRequest) error {
			inserted = request.Document
			return nil
		},
	}

	api := New(Config{}, catalog)

	var request = data.Document{
		Title:        strptr("title"),
		Uri:          strptr("uri"),
		DocumentKind: data.DocumentKind{Name: "uri"},
		Authors: []data.DocumentAuthor{
			{Name: "name", Surname: "surname"},
		},
		Tags: []data.DocumentTag{
			{Tag: "tag_0"},
		},
	}

	body, err := json.Marshal(request)
	require.NoError(t, err)

	router := api.catalogRouter()
	r := httptest.NewRecorder()
	router.ServeHTTP(r, httptest.NewRequest(http.MethodPost, "/catalog/documents", bytes.NewBuffer(body)))

	require.Equal(t, http.StatusOK, r.Code)
	require.Equal(t, request, inserted)
}

func TestCatalogApi_GetDocument(t *testing.T) {
	var documentID = 1
	var inserted = data.Document{
		Title: strptr("title"),
		Uri:   strptr("uri"),
		Authors: []data.DocumentAuthor{
			{Name: "name", Surname: "surname"},
		},
	}
	catalog := mockCatalog{
		getDocument: func(ctx context.Context, request data.GetDocumentRequest) (data.Document, error) {
			if request.DocumentID != documentID {
				return data.Document{}, errors.New("document not found")
			}
			return inserted, nil
		},
	}

	api := New(Config{}, catalog)

	router := api.catalogRouter()
	r := httptest.NewRecorder()
	router.ServeHTTP(r, httptest.NewRequest(http.MethodGet, "/catalog/documents/1", nil))
	require.Equal(t, http.StatusOK, r.Code)

	var response data.Document
	err := json.NewDecoder(r.Body).Decode(&response)
	require.NoError(t, err)
	require.Equal(t, inserted, response)
}

func TestCatalogApi_ListDocuments(t *testing.T) {
	var documents = []data.Document{
		{
			Title: strptr("title_1"),
			Uri:   strptr("uri_1"),
		},
		{
			Title: strptr("title_2"),
			Uri:   strptr("uri_2"),
		},
	}

	var (
		page          = 2
		pageSize      = 3
		totalElements = int64(100)
	)
	catalog := mockCatalog{
		listDocuments: func(ctx context.Context, request data.ListDocumentsRequest) (data.ListDocumentsResponse, error) {
			if request.Pagination.Page != page {
				return data.ListDocumentsResponse{}, errors.New("invalid page, expected page 2")
			}

			if request.Pagination.PageSize != pageSize {
				return data.ListDocumentsResponse{}, errors.New("invalid page size, expected page 3")
			}

			return data.ListDocumentsResponse{
				Items: documents,
				Pagination: data.PaginationResponse{
					TotalElements: totalElements,
					Page:          page,
					Pages:         10,
				},
			}, nil
		},
	}

	api := New(Config{}, catalog)

	router := api.catalogRouter()
	r := httptest.NewRecorder()
	router.ServeHTTP(r, httptest.NewRequest(http.MethodGet, fmt.Sprintf("/catalog/documents?page=%d&page_size=%d", page, pageSize), nil))
	require.Equal(t, http.StatusOK, r.Code)

	var response data.ListDocumentsResponse
	err := json.NewDecoder(r.Body).Decode(&response)
	require.NoError(t, err)
	require.Equal(t, documents, response.Items)
	require.Equal(t, data.PaginationResponse{
		TotalElements: totalElements,
		Page:          page,
		Pages:         10,
	}, response.Pagination)
}

type mockCatalog struct {
	insertDocument func(ctx context.Context, request data.InsertDocumentRequest) error
	getDocument    func(ctx context.Context, request data.GetDocumentRequest) (data.Document, error)
	listDocuments  func(ctx context.Context, request data.ListDocumentsRequest) (data.ListDocumentsResponse, error)
}

func (m mockCatalog) Init() error {
	return nil
}

func (m mockCatalog) InsertDocument(ctx context.Context, request data.InsertDocumentRequest) error {
	return m.insertDocument(ctx, request)
}

func (m mockCatalog) GetDocument(ctx context.Context, request data.GetDocumentRequest) (data.Document, error) {
	return m.getDocument(ctx, request)
}

func (m mockCatalog) ListDocuments(ctx context.Context, request data.ListDocumentsRequest) (data.ListDocumentsResponse, error) {
	return m.listDocuments(ctx, request)
}
