package group

type Group struct {
	ID            uint64   `json:"id"             db:"id"`
	Name          string   `json:"name"           db:"name"`
	AllowedEmails []string `json:"allowed_emails" db:"allowed_emails"`
}

func NewGroup(name string, allowedEmails []string) Group {
	return Group{
		Name:          name,
		AllowedEmails: allowedEmails,
	}
}
