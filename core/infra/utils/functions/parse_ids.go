package functions

import (
	"strconv"
	"strings"
)

func ParseIds(raw string) ([]int, error) {
	if strings.TrimSpace(raw) == "" {
		return []int{}, nil
	}

	parts := strings.Split(raw, ",")
	var ids []int
	for _, p := range parts {
		id, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
