package helper

import (
	"fmt"
	"os"

	"github.com/muzzarellimj/grace-material-api/internal/database"
	"github.com/muzzarellimj/grace-material-api/internal/database/service"
	model "github.com/muzzarellimj/grace-material-api/internal/model/book"
)

func FetchBook(constraint string) (model.Book, error) {
	zero := model.Book{}

	bookFragment, err := service.FetchFragment[model.BookFragment](database.Connection, database.TableBookFragments, constraint)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch book with constraint '%s': %v\n", constraint, err)

		return zero, err
	}

	authorFragmentSlice, err := fetchAuthorFragmentSlice(bookFragment)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch authors related to book '%d': %v\n", bookFragment.ID, err)
	}

	publisherFragmentSlice, err := fetchPublisherFragmentSlice(bookFragment)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch publishers related to book '%d': %v\n", bookFragment.ID, err)
	}

	topicFragmentSlice, err := fetchTopicFragmentSlice(bookFragment)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch topics related to book '%d': %v\n", bookFragment.ID, err)
	}

	book := mapBook(bookFragment, authorFragmentSlice, publisherFragmentSlice, topicFragmentSlice)

	return book, nil
}

func FetchBookSlice(constraintSlice []string) ([]model.Book, []error) {
	var bookSlice []model.Book
	var errSlice []error

	for _, constraint := range constraintSlice {
		book, err := FetchBook(constraint)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch and map book with constraint '%s': %v\n", constraint, err)

			errSlice = append(errSlice, err)
		}

		if book.ID != 0 {
			bookSlice = append(bookSlice, book)
		}
	}

	return bookSlice, errSlice
}

func FetchBookExistenceSlice() ([]int, []error) {
	idSlice, err := service.FetchExistenceSlice(database.Connection, database.TableBookFragments)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch existence slice: %v\n", err)

		return []int{}, []error{err}
	}

	if len(idSlice) == 0 {
		fmt.Fprint(os.Stderr, "Existence slice appears to be empty.")

		return []int{}, nil
	}

	return idSlice, nil
}

func fetchAuthorFragmentSlice(bookFragment model.BookFragment) ([]model.BookAuthorFragment, error) {
	bookAuthorRelationshipSlice, err := service.FetchRelationshipSlice[model.BookAuthorRelationship](database.Connection, database.TableBookAuthorRelationships, fmt.Sprintf("book=%d", bookFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between book '%d' and authors: %v\n", bookFragment.ID, err)

		return []model.BookAuthorFragment{}, err
	}

	var authorFragmentSlice []model.BookAuthorFragment

	for _, relationship := range bookAuthorRelationshipSlice {
		authorFragment, err := service.FetchFragment[model.BookAuthorFragment](database.Connection, database.TableBookAuthorFragments, fmt.Sprintf("id=%d", relationship.Author))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch author '%d': %v\n", relationship.Author, err)
		}

		if authorFragment.ID != 0 {
			authorFragmentSlice = append(authorFragmentSlice, authorFragment)
		}
	}

	return authorFragmentSlice, nil
}

func fetchPublisherFragmentSlice(bookFragment model.BookFragment) ([]model.BookPublisherFragment, error) {
	bookPublisherRelationshipSlice, err := service.FetchRelationshipSlice[model.BookPublisherRelationship](database.Connection, database.TableBookPublisherRelationships, fmt.Sprintf("book=%d", bookFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between book '%d' and publishers: %v\n", bookFragment.ID, err)

		return []model.BookPublisherFragment{}, err
	}

	var publisherFragmentSlice []model.BookPublisherFragment

	for _, relationship := range bookPublisherRelationshipSlice {
		publisherFragment, err := service.FetchFragment[model.BookPublisherFragment](database.Connection, database.TableBookPublisherFragments, fmt.Sprintf("id=%d", relationship.Publisher))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch publisher '%d': %v\n", relationship.Publisher, err)
		}

		if publisherFragment.ID != 0 {
			publisherFragmentSlice = append(publisherFragmentSlice, publisherFragment)
		}
	}

	return publisherFragmentSlice, nil
}

func fetchTopicFragmentSlice(bookFragment model.BookFragment) ([]model.BookTopicFragment, error) {
	bookTopicRelationshipSlice, err := service.FetchRelationshipSlice[model.BookTopicRelationship](database.Connection, database.TableBookTopicRelationships, fmt.Sprintf("book=%d", bookFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between book '%d' and topics: %v\n", bookFragment.ID, err)

		return []model.BookTopicFragment{}, err
	}

	var topicFragmentSlice []model.BookTopicFragment

	for _, relationship := range bookTopicRelationshipSlice {
		topicFragment, err := service.FetchFragment[model.BookTopicFragment](database.Connection, database.TableBookTopicFragments, fmt.Sprintf("id=%d", relationship.Topic))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch topic '%d': %v\n", relationship.Topic, err)
		}

		if topicFragment.ID != 0 {
			topicFragmentSlice = append(topicFragmentSlice, topicFragment)
		}
	}

	return topicFragmentSlice, nil
}

func mapBook(bookFragment model.BookFragment, authorFragmentSlice []model.BookAuthorFragment, publisherFragmentSlice []model.BookPublisherFragment, topicFragmentSlice []model.BookTopicFragment) model.Book {
	if authorFragmentSlice == nil {
		authorFragmentSlice = make([]model.BookAuthorFragment, 0)
	}

	if publisherFragmentSlice == nil {
		publisherFragmentSlice = make([]model.BookPublisherFragment, 0)
	}

	if topicFragmentSlice == nil {
		topicFragmentSlice = make([]model.BookTopicFragment, 0)
	}

	return model.Book{
		ID:               bookFragment.ID,
		Title:            bookFragment.Title,
		Subtitle:         bookFragment.Subtitle,
		Description:      bookFragment.Description,
		Authors:          authorFragmentSlice,
		Publishers:       publisherFragmentSlice,
		Topics:           topicFragmentSlice,
		PublishDate:      bookFragment.PublishDate,
		Pages:            bookFragment.Pages,
		ISBN10:           bookFragment.ISBN10,
		ISBN13:           bookFragment.ISBN13,
		Image:            bookFragment.Image,
		EditionReference: bookFragment.EditionReference,
		WorkReference:    bookFragment.WorkReference,
	}
}
