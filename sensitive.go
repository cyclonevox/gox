package gox

import (
	`fmt`
	`strings`
)

type Sensitive string

func (s Sensitive) MarshalJSON() ([]byte, error) {
	if len(s) == 15 {
		hiding := fmt.Sprintf(`"%s%s%s"`, s[:7], "****", s[11:])

		return []byte(hiding), nil
	}

	if len(s) == 18 {
		hiding := fmt.Sprintf(`"%s%s%s"`, s[:3], "***********", s[14:])

		return []byte(hiding), nil
	}

	if len(s) > 1 {
		hiding := fmt.Sprintf(`"%s%s"`, s[:1], strings.Repeat("*", len(s)-1))

		return []byte(hiding), nil
	}

	return []byte(fmt.Sprintf(`"%s"`, s)), nil
}
