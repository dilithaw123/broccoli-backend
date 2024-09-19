package group

type Group struct {
	ID            uint64   `json:"id"             db:"id"`
	Name          string   `json:"name"           db:"name"`
	AllowedEmails []string `json:"allowed_emails" db:"allowed_emails"`
	Timezone      string   `json:"timezone"       db:"timezone"`
}

func NewGroup(name string, allowedEmails []string, tz string) Group {
	return Group{
		Name:          name,
		AllowedEmails: allowedEmails,
		Timezone:      tz,
	}
}
