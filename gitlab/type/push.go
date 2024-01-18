package gitlabtype

type WebhookPush struct {
	ObjectKind   string  `json:"object_kind"`
	EventName    string  `json:"event_name"`
	Before       string  `json:"before"`
	After        string  `json:"after"`
	Ref          string  `json:"ref"`
	CheckoutSHA  string  `json:"checkout_sha"`
	UserID       int     `json:"user_id"`
	UserName     string  `json:"user_name"`
	UserUsername string  `json:"user_username"`
	UserEmail    string  `json:"user_email"`
	UserAvatar   string  `json:"user_avatar"`
	ProjectID    int     `json:"project_id"`
	Project      Project `json:"project"`
	Repository   struct {
		Name        string `json:"name"`
		URL         string `json:"url"`
		Description string `json:"description"`
		Homepage    string `json:"homepage"`
	} `json:"repository"`
	Commits      []Commit `json:"commits"`
	TotalCommits int      `json:"total_commits_count"`
}
