package test

import (
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http/httptest"
	"testing"
)

func TestBookGet(t *testing.T) {
	req := NewRequestAPI("GET", "/books/9782266155489", nil)
	resp := httptest.NewRecorder()
	books := sqlmock.NewRows([]string{"", "", "", "", "", "", "", "", ""})
	books.AddRow(1, "9782266155489", nil, "Some description", 9.99,
		1, "Dune", []byte(`{"(,)"}`), []byte(`{"(,,)"}`),
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
