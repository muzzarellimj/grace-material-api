package service

import (
	"fmt"
	"os"

	"github.com/muzzarellimj/grace-material-api/database"
)

func FetchRelationship[M interface{}](connection database.PgxConnection, table string, constraint string) (M, error) {
	var zero M

	relationshipSlice, err := FetchFragmentSlice[M](connection, table, constraint)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch initial relationship slice: %s\n", err)

		return zero, err
	}

	if len(relationshipSlice) > 0 {
		return relationshipSlice[0], nil
	}

	return zero, nil
}

func FetchRelationshipSlice[M interface{}](connection database.PgxConnection, table string, constraint string) ([]M, error) {
	var zero []M

	relationshipSlice, err := FetchFragmentSlice[M](connection, table, constraint)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch relationship slice: %v\n", err)

		return zero, err
	}

	return relationshipSlice, nil
}
