package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garugaru/knowledge/pkg/data"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (a Api) catalogRouter() *mux.Router {
	router := mux.NewRouter()
	router.Path("/catalog/documents").Methods(http.MethodPost).HandlerFunc(a.catalogInsertDocument)
	router.Path("/catalog/documents").Methods(http.MethodGet).HandlerFunc(a.catalogListDocument)
	router.Path("/catalog/documents/{id:[0-9]+}").Methods(http.MethodGet).HandlerFunc(a.catalogGetDocument)
	return router
}

func (a Api) catalogGetDocument(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, present := params["id"]
	if !present {
		httpErr(w, errors.New("'id' parameter must be present"), http.StatusBadRequest)
		return
	}

	documentID, err := strconv.Atoi(id)
	if err != nil {
		httpErr(w, err, http.StatusBadRequest)
		return
	}

	document, err := a.catalog.GetDocument(r.Context(), data.GetDocumentRequest{
		DocumentID: documentID,
	})

	if err != nil {
		httpErr(w, err, http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(document); err != nil {
		httpErr(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a Api) catalogListDocument(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	var request = data.ListDocumentsRequest{
		Pagination: data.PaginationRequest{
			Page:     1,
			PageSize: 100,
		},
	}

	title, present := params["title"]
	if present && len(title) > 0 {
		request.Title = title[0]
	}

	tags, present := params["tags"]
	if present {
		request.Tags = tags
	}

	pageParam, present := params["page"]
	if present {
		page, err := strconv.Atoi(pageParam[0])
		if err != nil {
			httpErr(w, fmt.Errorf("invalid 'page' parameter value: %s", pageParam), http.StatusBadRequest)
			return
		}
		request.Pagination.Page = page
	}

	pageSizeParam, present := params["page_size"]
	if present {
		pageSize, err := strconv.Atoi(pageSizeParam[0])
		if err != nil {
			httpErr(w, fmt.Errorf("invalid 'page_size' parameter value: %s", pageParam), http.StatusBadRequest)
			return
		}
		request.Pagination.PageSize = pageSize
	}

	document, err := a.catalog.ListDocuments(r.Context(), request)

	if err != nil {
		httpErr(w, err, http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(document); err != nil {
		httpErr(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a Api) catalogInsertDocument(w http.ResponseWriter, r *http.Request) {
	var document data.Document

	if err := json.NewDecoder(r.Body).Decode(&document); err != nil {
		httpErr(w, err, http.StatusBadRequest)
		return
	}

	err := a.catalog.InsertDocument(r.Context(), data.InsertDocumentRequest{
		Document: document,
	})

	if err != nil {
		httpErr(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
