package data

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"path"
	"testing"
	"time"
)

func TestDBStore_Init(t *testing.T) {
	tmpDb := path.Join(t.TempDir(), t.Name())
	db, err := gorm.Open(sqlite.Open(tmpDb), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	require.NoError(t, err)

	store := NewDBStore(db)
	err = store.Init()
	require.NoError(t, err)
}

func TestDBStore_InsertDocument(t *testing.T) {
	tmpDb := path.Join(t.TempDir(), t.Name())
	db, err := gorm.Open(sqlite.Open(tmpDb), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	require.NoError(t, err)

	store := NewDBStore(db)
	err = store.Init()
	require.NoError(t, err)

	err = store.InsertDocument(context.TODO(), InsertDocumentRequest{
		Document: Document{
			Title: "Test Title",
			Uri:   "file://test.txt",
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
			PublishDate: time.Date(1996, 7, 10, 19, 0, 0, 0, time.UTC),
		},
	})

	require.NoError(t, err)
}

func TestDBStore_Search(t *testing.T) {
	tmpDb := path.Join(t.TempDir(), t.Name())
	db, err := gorm.Open(sqlite.Open(tmpDb), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	require.NoError(t, err)

	store := NewDBStore(db)
	err = store.Init()
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		err = store.InsertDocument(context.TODO(), InsertDocumentRequest{
			Document: Document{
				Title: fmt.Sprintf("Test Title %d", i),
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

	documents, err := store.Search(context.TODO(), SearchRequest{
		Title: "Test",
		Tags:  []string{"book"},
		Pagination: PaginationRequest{
			Page:     1,
			PageSize: 2,
		},
	})

	require.NoError(t, err)
	require.Len(t, documents, 2)

	require.Equal(t, "Test Title 0", documents[0].Title)
	require.Equal(t, "Test Title 1", documents[1].Title)

	for _, doc := range documents {
		require.NotEmpty(t, doc.Tags)
		require.NotEmpty(t, doc.Authors)
	}

	documents, err = store.Search(context.TODO(), SearchRequest{
		Title: "Book",
		Tags:  []string{"book"},
		Pagination: PaginationRequest{
			Page:     1,
			PageSize: 2,
		},
	})

	require.NoError(t, err)
	require.Empty(t, documents, "title should not match")

	documents, err = store.Search(context.TODO(), SearchRequest{
		Tags: []string{"tag_0"},
		Pagination: PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
	})

	require.NoError(t, err)
	require.Len(t, documents, 1, "tag only search must yield only 1 result")
}
