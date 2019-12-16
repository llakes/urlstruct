package urlstruct_test

import (
	"fmt"
	"net/url"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/urlstruct"
)

type Book struct {
	tableName struct{} `pg:"alias:b"` //nolint:unused,structcheck

	ID        int64
	AuthorID  int64
	CreatedAt time.Time
}

type BookFilter struct {
	tableName struct{} `urlstruct:"b"` //nolint:unused,structcheck

	urlstruct.Pager
	AuthorID int64
}

func ExampleUnmarshal_filter() {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "",
		Database: "postgres",
	})
	defer db.Close()

	values := url.Values{
		"author_id": {"123"},
		"page":      {"2"},
		"limit":     {"100"},
	}
	filter := new(BookFilter)
	err := urlstruct.Unmarshal(values, filter)
	if err != nil {
		panic(err)
	}

	filter.Pager.MaxLimit = 100     // default max limit is 1000
	filter.Pager.MaxOffset = 100000 // default max offset is 1000000

	// Following query generates:
	//
	// SELECT "b"."id", "b"."author_id", "b"."created_at"
	// FROM "books" AS "b"
	// WHERE "b".author_id = 123
	// LIMIT 100 OFFSET 100

	var books []*Book
	_ = db.Model(&books).
		WhereStruct(filter).
		Limit(filter.Pager.GetLimit()).
		Offset(filter.Pager.GetOffset()).
		Select()

	fmt.Println("author", filter.AuthorID)
	fmt.Println("limit", filter.GetLimit())
	fmt.Println("offset", filter.GetOffset())
	// Output: author 123
	// limit 100
	// offset 100
}
