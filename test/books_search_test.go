package test

import (
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http/httptest"
	"testing"
)

func TestBookSearch(t *testing.T) {
	req := NewRequestAPI("GET", "/books/?search=foo", nil)
	resp := httptest.NewRecorder()
	books := sqlmock.NewRows([]string{"id", "isbn", "title", "description",
		"price", "number", "serie", "authors"})
	books.AddRow(1, "1234567890", nil, "Some description", 9.99,
		1, "Foo", []byte(`{"(1,\"author 1\")","(2,\"author 2\")"}`))
	mock.ExpectQuery("SELECT").WithArgs(10, 0, "foo").WillReturnRows(books)

	rulz.ServeHTTP(resp, req)
	result := resp.Result()

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	if result.StatusCode != 200 {
		t.Errorf("expected status code 200 not met, got: %d", result.StatusCode)
	}
	helperBodyCompare(t, result, "book_search.json")
}

func TestBookSearchEmpty(t *testing.T) {
	req := NewRequestAPI("GET", "/books/?search=foo", nil)
	resp := httptest.NewRecorder()
	books := sqlmock.NewRows([]string{"id", "isbn", "title", "description",
		"price", "number", "serie", "authors"})
	mock.ExpectQuery("SELECT").WithArgs(10, 0, "foo").WillReturnRows(books)

	rulz.ServeHTTP(resp, req)
	result := resp.Result()

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	if result.StatusCode != 200 {
		t.Errorf("expected status code 200 not met, got: %d", result.StatusCode)
	}
	helperBodyCompare(t, result, "book_search_empty.json")
}
