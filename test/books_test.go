package test

import (
	"database/sql"
	"fmt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http/httptest"
	"testing"
)

func TestBookGet(t *testing.T) {
	req := NewRequestAPI("GET", "/books/9782266155489", nil)
	resp := httptest.NewRecorder()
	books := sqlmock.NewRows([]string{"", "", "", "", "", "", "", "", ""})
	books.AddRow(1, "9782266155489", nil, "Some description", 9.99,
		1, "Dune", []byte(`{"(1,\"Frank Herbert\")","(2,\"Michel Demuth\")"}`), []byte(`{"(,,)"}`),
	)
	mock.ExpectQuery("SELECT").WithArgs("9782266155489").WillReturnRows(books)

	rulz.ServeHTTP(resp, req)
	result := resp.Result()

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	if result.StatusCode != 200 {
		t.Errorf("expected status code 200 not met, got: %d", result.StatusCode)
	}
	helperBodyCompare(t, result, "book_get.json")
}

func TestBookGet404(t *testing.T) {
	req := NewRequestAPI("GET", "/books/9782266155489", nil)
	resp := httptest.NewRecorder()

	mock.ExpectQuery("SELECT").WithArgs("9782266155489").WillReturnError(sql.ErrNoRows)
	rulz.ServeHTTP(resp, req)
	result := resp.Result()

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	if result.StatusCode != 404 {
		t.Errorf("expected status code 404 not met, got: %d", result.StatusCode)
	}
	helperBodyCompare(t, result, "book_get_404.json")
}

func TestBookGet500(t *testing.T) {
	req := NewRequestAPI("GET", "/books/9782266155489", nil)
	resp := httptest.NewRecorder()
	err := fmt.Errorf("unexpected error")

	mock.ExpectQuery("SELECT").WithArgs("9782266155489").WillReturnError(err)
	rulz.ServeHTTP(resp, req)
	result := resp.Result()

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	if result.StatusCode != 500 {
		t.Errorf("expected status code 500 not met, got: %d", result.StatusCode)
	}
	helperBodyCompare(t, result, "book_get_500.json")
}
