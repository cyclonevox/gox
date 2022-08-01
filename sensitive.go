package gox

import (
	`fmt`
	`regexp`
	`strings`
)

type Sensitive string

func (s Sensitive) MarshalJSON() ([]byte, error) {
	mobile := `^[+]86[-]1([3,4,5,6,7,8,9][0-9])\d{8}$`
	reg := regexp.MustCompile(mobile)

	if reg.MatchString(string(s)) {
		hiding := fmt.Sprintf(`"%s%s%s"`, s[:7], "****", s[11:])

		return []byte(hiding), nil
	}

	if len(s) > 1 {
		hiding := fmt.Sprintf(`"%s%s"`, s[:1], strings.Repeat("*", len(s)-1))

		return []byte(hiding), nil
	}

	return []byte(fmt.Sprintf(`"%s"`, s)), nil
}
