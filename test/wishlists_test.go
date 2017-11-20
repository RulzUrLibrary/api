package test

import (
	"fmt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http/httptest"
	"testing"
)

func TestWishlists401(t *testing.T) {
	req := NewRequestAPI("GET", "/wishlists/", nil)
	req.SetBasicAuth("foo", "bar")
	resp := httptest.NewRecorder()

	users := sqlmock.NewRows([]string{"auth", "id"})
	users.AddRow(false, 0)

	mock.ExpectQuery("SELECT").WithArgs("foo", "bar").WillReturnRows(users)
	rulz.ServeHTTP(resp, req)
	result := resp.Result()

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	if result.StatusCode != 401 {
		t.Errorf("expected status code 401 not met, got: %d", result.StatusCode)
	}
	helperBodyCompare(t, result, "401.json")
}

func TestWishlists500(t *testing.T) {
	req := NewRequestAPI("GET", "/wishlists/", nil)
	req.SetBasicAuth("foo", "bar")
	resp := httptest.NewRecorder()
	err := fmt.Errorf("unexpected error")

	users := sqlmock.NewRows([]string{"auth", "id"})
	users.AddRow(true, 1)
	mock.ExpectQuery("SELECT").WithArgs("foo", "bar").WillReturnRows(users)
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnError(err)
	rulz.ServeHTTP(resp, req)
	result := resp.Result()

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	if result.StatusCode != 500 {
		t.Errorf("expected status code 500 not met, got: %d", result.StatusCode)
	}
	helperBodyCompare(t, result, "500.json")
}

func TestWishlistsGet(t *testing.T) {
	req := NewRequestAPI("GET", "/wishlists/", nil)
	req.SetBasicAuth("foo", "bar")
	resp := httptest.NewRecorder()

	users := sqlmock.NewRows([]string{"auth", "id"})
	users.AddRow(true, 1)
	mock.ExpectQuery("SELECT").WithArgs("foo", "bar").WillReturnRows(users)

	count := sqlmock.NewRows([]string{"count"})
	count.AddRow(20)
	mock.ExpectQuery("SELECT COUNT").WithArgs(1).WillReturnRows(count)

	wishlists := sqlmock.NewRows([]string{"id", "name", "uuid", "book.id", "book.isbn",
		"book.title", "book.price", "book.num", "serie", "book.authors"})

	wishlists.AddRow(1, "wishlist_1", "uuid_1", 1, "book_1_isbn", "book_1_title",
		1.11, 1, "book_1_serie", []byte(`{"(,)"}`))
	wishlists.AddRow(1, "wishlist_1", "uuid_1", 2, "book_2_isbn", nil,
		1.11, 1, "book_2_serie", []byte(`{"(,)"}`))
	wishlists.AddRow(2, "wishlist_2", "uuid_2", 3, "book_3_isbn", "book_3_title",
		1.11, nil, nil, []byte(`{"(,)"}`))

	mock.ExpectQuery("SELECT").WithArgs(10, 0, 1).WillReturnRows(wishlists)

	rulz.ServeHTTP(resp, req)
	result := resp.Result()

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	if result.StatusCode != 200 {
		t.Errorf("expected status code 200 not met, got: %d", result.StatusCode)
	}
	helperBodyCompare(t, result, "wishlists_get.json")
}
