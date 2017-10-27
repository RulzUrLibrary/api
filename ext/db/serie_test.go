package db

//import (
//	"database/sql"
//	"fmt"
//	"testing"
//)
//
//type Result struct {
//	numbers []int
//	size    int
//}
//
//func (r Result) String() string {
//	return fmt.Sprintf("{numbers: %#v, size: %d}", r.numbers, r.size)
//}
//
//func (r1 Result) Equal(r2 Result) (ok bool) {
//	if len(r1.numbers) != len(r2.numbers) || r1.size != r2.size {
//		return
//	}
//	for i, _ := range r1.numbers {
//		if r1.numbers[i] != r2.numbers[i] {
//			return
//		}
//	}
//	return true
//}
//
//type testPairs struct {
//	numbers Numbers
//	result  Result
//}
//
//var tests = []testPairs{
//	{Numbers{Number{sql.NullInt64{Valid: false}}}, Result{[]int{}, 0}},
//	{Numbers{
//		Number{sql.NullInt64{Valid: true, Int64: int64(1)}},
//		Number{sql.NullInt64{Valid: true, Int64: int64(2)}},
//		Number{sql.NullInt64{Valid: true, Int64: int64(2)}},
//	}, Result{[]int{1, 2}, 2}},
//	{Numbers{
//		Number{sql.NullInt64{Valid: true, Int64: int64(1)}},
//		Number{sql.NullInt64{Valid: false, Int64: int64(2)}},
//		Number{sql.NullInt64{Valid: false, Int64: int64(2)}},
//	}, Result{[]int{1}, 2}},
//}
//
//func TestNumbersToStructs(t *testing.T) {
//	var result Result
//	for _, pair := range tests {
//		result.numbers, result.size = pair.numbers.ToStructs()
//		if !result.Equal(pair.result) {
//			t.Error(
//				"For", pair.numbers,
//				"expected", pair.result,
//				"got", result,
//			)
//		}
//	}
//}
