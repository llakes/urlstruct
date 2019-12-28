# urlstruct decodes url.Values into structs

[![Build Status](https://travis-ci.org/go-pg/urlstruct.svg?branch=master)](https://travis-ci.org/go-pg/urlstruct)
[![GoDoc](https://godoc.org/github.com/go-pg/urlstruct?status.svg)](https://godoc.org/github.com/go-pg/urlstruct)

## Example

Following example decodes URL query `?page=2&limit=100&author_id=123...` into a struct and uses [go-pg](https://github.com/go-pg/pg) feature `WhereStruct` to autogenerate WHERE clause:

```go
type Book struct {
	tableName struct{} `pg:"alias:b"`

	ID        int64
	AuthorID  int64
	FieldID    int64
	Field      string
	CreatedAt time.Time
}

type BookFilter struct {
	tableName struct{} `urlstruct:"b"`

	urlstruct.Pager
	AuthorID  int64
	AuthorIdNEQ  int64
	Field      string
	FieldNEQ   string
	FieldIEQ   string
	FieldMatch string
	FieldIdLT    int64
	FieldIdLTE   int64
	FieldIdGT    int64
	FieldIdGTE   int64
	CreatedAt string
	CreatedAtLT string
	CreatedAtGT string
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
		"author_id__exclude": {"123"},
		"field": {"123"},
		"field__ieq": {"123"},
		"field__exclude": {"123"},
		"field_id": {"123"},
		"field_id__ieq": {"123"},
		"field_id__exclude": {"123"},
		"field_id__gt": {"123"},
		"field_id__lt": {"123"},
		"created_at": {"2006-01-02"},
		"created_at__gt": {"2006-01-02"},
		"created_at__lt": {"2006-01-02"},
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
	//...
	// LIMIT 100 OFFSET 100

	var books []*Book
	_ = db.Model(&books).
		WhereStruct(filter).
		Limit(filter.Pager.GetLimit()).
		Offset(filter.Pager.GetOffset()).
		Select()

	fmt.Println("author", filter.AuthorID)
	fmt.Println("limit", filter.GetLimit())
	fmt.Println("offset", filter.GetLimit())
	// Output: author 123
	// limit 100
	// offset 100
}
```
