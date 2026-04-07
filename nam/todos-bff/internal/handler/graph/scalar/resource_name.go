package scalar

import (
	"fmt"
	"io"
	"strconv"
)

type ResourceName string

// MarshalGQL serializes ResourceName for GraphQL responses
func (r ResourceName) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(string(r)))
}

// UnmarshalGQL deserializes ResourceName from GraphQL input
func (r *ResourceName) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("resource name must be a string")
	}
	*r = ResourceName(str)
	return nil
}
