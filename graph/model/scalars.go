package model

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalDate(t time.Time) graphql.Marshaler {
	if t.IsZero() {
		return graphql.Null
	}

	return graphql.WriterFunc(func(w io.Writer) {
		date := strings.Split(t.Format(time.RFC3339), "T")[0]
		io.WriteString(w, strconv.Quote(date))
	})
}

func UnmarshalDate(v interface{}) (time.Time, error) {
	if tmpStr, ok := v.(string); ok {
		return time.Parse(time.RFC3339, tmpStr+"T00:00:00Z") //TODO should take care of TimeZone better
	}
	return time.Time{}, errors.New("Date should be YYYY-MM-DD format")
}
