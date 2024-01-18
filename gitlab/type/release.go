package gitlabtype

type Sources struct {
	Format string `json:"format"`
	URL    string `json:"url"`
}
type Assets struct {
	Count   int           `json:"count"`
	Links   []interface{} `json:"links"`
	Sources []Sources     `json:"sources"`
}

type WebhookRelease struct {
	ID          int     `json:"id"`
	CreatedAt   string  `json:"created_at"`
	Description string  `json:"description"`
	Name        string  `json:"name"`
	ReleasedAt  string  `json:"released_at"`
	Tag         string  `json:"tag"`
	ObjectKind  string  `json:"object_kind"`
	Project     Project `json:"project"`
	URL         string  `json:"url"`
	Action      string  `json:"action"`
	Assets      Assets  `json:"assets"`
	Commit      Commit  `json:"commit"`
}
