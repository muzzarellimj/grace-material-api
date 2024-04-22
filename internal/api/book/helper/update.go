package helper

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/muzzarellimj/grace-material-api/internal/database"
	"github.com/muzzarellimj/grace-material-api/internal/database/service"
	model "github.com/muzzarellimj/grace-material-api/internal/model/book"
)

func UpdateBookFragment(book model.BookFragment) (int, error) {
	id, err := service.UpdateFragment(database.Connection, database.TableBookFragments, database.PropertiesBookFragments, fmt.Sprintf("id=%d", book.ID), pgx.NamedArgs{
		"title":             book.Title,
		"subtitle":          book.Subtitle,
		"description":       book.Description,
		"publish_date":      book.PublishDate,
		"pages":             book.Pages,
		"isbn10":            book.ISBN10,
		"isbn13":            book.ISBN13,
		"image":             book.Image,
		"edition_reference": book.EditionReference,
		"work_reference":    book.WorkReference,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to update book '%d' fragment: %v\n", book.ID, err)

		return 0, err
	}

	return id, nil
}
