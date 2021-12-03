package api

import (
	"fmt"
	"strings"
)

// Errors is a collection of REST errors
type Errors struct {
	Errors []*Error `json:"errors"`
}

func (es *Errors) Error() string {
	if len(es.Errors) == 1 {
		return fmt.Sprintf("1 error occurred:\n\t* %s\n", es.Errors[0])
	}

	points := make([]string, len(es.Errors))
	for i, err := range es.Errors {
		points[i] = fmt.Sprintf("* %s", err)
	}

	return fmt.Sprintf(
		"%d errors occurred:\n\t%s\n\n",
		len(es.Errors), strings.Join(points, "\n\t"))
}

// Error represent an error returned by the REST API
type Error struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func (e *Error) Error() string {
	detail := strings.Replace(e.Detail, "\"", "'", -1)
	detail = strings.Replace(detail, "\n", "", -1)
	detail = strings.Replace(detail, "\t", "", -1)
	return fmt.Sprintf("ID: %q, Status: %d, Title: %q, Detail: %q", e.ID, e.Status, e.Title, detail)
}
