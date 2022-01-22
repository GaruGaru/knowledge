package data

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"path"
	"testing"
)

func TestDBCatalog_Init(t *testing.T) {
	tmpDb := path.Join(t.TempDir(), t.Name())
	db, err := gorm.Open(sqlite.Open(tmpDb), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	require.NoError(t, err)

	catalog := NewDBCatalog(db)
	err = catalog.Init()
	require.NoError(t, err)
}

func TestDBCatalog_InsertDocument(t *testing.T) {
	tmpDb := path.Join(t.TempDir(), t.Name())
	db, err := gorm.Open(sqlite.Open(tmpDb), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	require.NoError(t, err)

	catalog := NewDBCatalog(db)
	err = catalog.Init()
	require.NoError(t, err)

	requests := []InsertDocumentRequest{
		{
			Document: Document{
				Title: strptr("Test Title"),
				Uri:   strptr("file://test.txt"),
				DocumentKind: DocumentKind{
					Name: "file",
				},
				Authors: []DocumentAuthor{
					{Name: "Me", Surname: "Me"},
					{Name: "Myself", Surname: "Myself"},
				},
				Tags: []DocumentTag{
					{Tag: "test"},
					{Tag: "book"},
				},
			},
		},
		{
			Document: Document{
				Title: strptr("Test Title"),
				Uri:   strptr("file://test.txt"),
				DocumentKind: DocumentKind{
					Name: "file",
				},
				Tags: []DocumentTag{
					{Tag: "test"},
					{Tag: "book"},
				},
			},
		},
		{
			Document: Document{
				Title: strptr("Test Title"),
				Uri:   strptr("file://test.txt"),
				DocumentKind: DocumentKind{
					Name: "file",
				},
			},
		},
	}

	for _, request := range requests {
		err = catalog.InsertDocument(context.TODO(), request)
		require.NoError(t, err)
	}
}

func TestDBCatalog_InsertDocument_Invalid(t *testing.T) {
	tmpDb := path.Join(t.TempDir(), t.Name())
	db, err := gorm.Open(sqlite.Open(tmpDb), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	require.NoError(t, err)

	catalog := NewDBCatalog(db)
	err = catalog.Init()
	require.NoError(t, err)

	requests := []InsertDocumentRequest{
		{
			Document: Document{},
		},
	}

	for _, request := range requests {
		err = catalog.InsertDocument(context.TODO(), request)
		require.Error(t, err)
	}
}

func TestDBCatalog_ListDocuments(t *testing.T) {
	tmpDb := path.Join(t.TempDir(), t.Name())
	db, err := gorm.Open(sqlite.Open(tmpDb), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	require.NoError(t, err)

	catalog := NewDBCatalog(db)
	err = catalog.Init()
	require.NoError(t, err)

	const documentsCount = 10
	for i := 0; i < documentsCount; i++ {
		err = catalog.InsertDocument(context.TODO(), InsertDocumentRequest{
			Document: Document{
				Title: strptr(fmt.Sprintf("Test Title %d", i)),
				Uri:   strptr(""),
				Authors: []DocumentAuthor{
					{Name: "Author"},
				},
				Tags: []DocumentTag{
					{Tag: "test"},
					{Tag: "book"},
					{Tag: fmt.Sprintf("tag_%d", i)},
				},
			},
		})
		require.NoError(t, err)
	}

	documents, err := catalog.ListDocuments(context.TODO(), ListDocumentsRequest{
		Title: "Test",
		Tags:  []string{"book"},
		Pagination: PaginationRequest{
			Page:     1,
			PageSize: 2,
		},
	})

	require.NoError(t, err)
	require.Len(t, documents.Items, 2)

	require.Equal(t, "Test Title 0", *documents.Items[0].Title)
	require.Equal(t, "Test Title 1", *documents.Items[1].Title)

	require.Equal(t, PaginationResponse{
		TotalElements: documentsCount,
		Page:          1,
		Pages:         documentsCount / 2,
	}, documents.Pagination)

	for _, doc := range documents.Items {
		require.NotEmpty(t, doc.Tags)
		require.NotEmpty(t, doc.Authors)
	}

	documents, err = catalog.ListDocuments(context.TODO(), ListDocumentsRequest{
		Title: "Book",
		Tags:  []string{"book"},
		Pagination: PaginationRequest{
			Page:     1,
			PageSize: 2,
		},
	})

	require.NoError(t, err)
	require.Empty(t, documents.Items, "title should not match")

	documents, err = catalog.ListDocuments(context.TODO(), ListDocumentsRequest{
		Tags: []string{"tag_0"},
		Pagination: PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
	})

	require.NoError(t, err)
	require.Len(t, documents.Items, 1, "tag only search must yield only 1 result")
}

func TestDBCatalog_GetDocument(t *testing.T) {
	tmpDb := path.Join(t.TempDir(), t.Name())
	db, err := gorm.Open(sqlite.Open(tmpDb), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	require.NoError(t, err)

	catalog := NewDBCatalog(db)
	err = catalog.Init()
	require.NoError(t, err)

	insertedDocument := Document{
		Title: strptr("Test Title"),
		Uri:   strptr("file://test.txt"),
		DocumentKind: DocumentKind{
			Name: "file",
		},
		Authors: []DocumentAuthor{
			{Name: "Me", Surname: "Me"},
			{Name: "Myself", Surname: "Myself"},
		},
		Tags: []DocumentTag{
			{Tag: "test"},
			{Tag: "book"},
		},
	}

	err = catalog.InsertDocument(context.TODO(), InsertDocumentRequest{
		Document: insertedDocument,
	})

	require.NoError(t, err)

	document, err := catalog.GetDocument(context.TODO(), GetDocumentRequest{
		DocumentID: 1,
	})

	require.NoError(t, err)
	require.Equal(t, 1, document.ID)
	require.Equal(t, *insertedDocument.Title, *document.Title)
	require.Equal(t, *insertedDocument.Uri, *document.Uri)
	require.Equal(t, insertedDocument.Tags[0].Tag, document.Tags[0].Tag)
	require.Equal(t, insertedDocument.Authors[0].Name, document.Authors[0].Name)

	document, err = catalog.GetDocument(context.TODO(), GetDocumentRequest{
		DocumentID: 9999,
	})
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
