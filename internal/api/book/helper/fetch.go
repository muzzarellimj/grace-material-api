package helper

import (
	"fmt"
	"os"

	"github.com/muzzarellimj/grace-material-api/internal/database"
	"github.com/muzzarellimj/grace-material-api/internal/database/connection"
	"github.com/muzzarellimj/grace-material-api/internal/database/service"
	model "github.com/muzzarellimj/grace-material-api/internal/model/book"
)

func FetchBook(constraint string) (model.Book, error) {
	bookFragment, err := service.FetchFragment[model.BookFragment](connection.Book, database.TableBookFragments, constraint)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch book with constraint '%s': %v\n", constraint, err)

		return model.Book{}, err
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

	return mapBook(bookFragment, authorFragmentSlice, publisherFragmentSlice, topicFragmentSlice), nil
}

func fetchAuthorFragmentSlice(bookFragment model.BookFragment) ([]model.BookAuthorFragment, error) {
	bookAuthorRelationshipSlice, err := service.FetchRelationshipSlice[model.BookAuthorRelationship](connection.Book, database.TableBookAuthorRelationships, fmt.Sprintf("book=%d", bookFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between book '%d' and authors: %v\n", bookFragment.ID, err)

		return []model.BookAuthorFragment{}, err
	}

	var authorFragmentSlice []model.BookAuthorFragment

	for _, relationship := range bookAuthorRelationshipSlice {
		authorFragment, err := service.FetchFragment[model.BookAuthorFragment](connection.Book, database.TableBookAuthorFragments, fmt.Sprintf("id=%d", relationship.Author))

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
	bookPublisherRelationshipSlice, err := service.FetchRelationshipSlice[model.BookPublisherRelationship](connection.Book, database.TableBookPublisherRelationships, fmt.Sprintf("book=%d", bookFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between book '%d' and publishers: %v\n", bookFragment.ID, err)

		return []model.BookPublisherFragment{}, err
	}

	var publisherFragmentSlice []model.BookPublisherFragment

	for _, relationship := range bookPublisherRelationshipSlice {
		publisherFragment, err := service.FetchFragment[model.BookPublisherFragment](connection.Book, database.TableBookPublisherFragments, fmt.Sprintf("id=%d", relationship.Publisher))

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
	bookTopicRelationshipSlice, err := service.FetchRelationshipSlice[model.BookTopicRelationship](connection.Book, database.TableBookTopicRelationships, fmt.Sprintf("book=%d", bookFragment.ID))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationships between book '%d' and topics: %v\n", bookFragment.ID, err)

		return []model.BookTopicFragment{}, err
	}

	var topicFragmentSlice []model.BookTopicFragment

	for _, relationship := range bookTopicRelationshipSlice {
		topicFragment, err := service.FetchFragment[model.BookTopicFragment](connection.Book, database.TableBookTopicFragments, fmt.Sprintf("id=%d", relationship.Topic))

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
