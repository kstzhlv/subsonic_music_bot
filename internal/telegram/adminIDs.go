package telegram

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseAdminChatIDs(value string) ([]int64, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}

	parts := strings.Split(value, ",")
	ids := make([]int64, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, fmt.Errorf(
				"parse admin chat IDs %q: %w",
				part,
				err,
			)
		}

		ids = append(ids, id)
	}

	return ids, nil
}
