package gitlabtype

type ObjectAttributes struct {
	ID                    int         `json:"id"`
	TargetBranch          string      `json:"target_branch"`
	SourceBranch          string      `json:"source_branch"`
	SourceProjectID       int         `json:"source_project_id"`
	AuthorID              int         `json:"author_id"`
	AssigneeID            int         `json:"assignee_id"`
	Title                 string      `json:"title"`
	CreatedAt             string      `json:"created_at"`
	UpdatedAt             string      `json:"updated_at"`
	StCommits             interface{} `json:"st_commits"`
	StDiffs               interface{} `json:"st_diffs"`
	MergeError            string      `json:"merge_error"`
	State                 string      `json:"state"`
	MergeStatus           string      `json:"merge_status"`
	TargetProjectID       int         `json:"target_project_id"`
	IID                   int         `json:"iid"`
	Description           string      `json:"description"`
	Source                Project     `json:"source"`
	Target                Project     `json:"target"`
	LastCommit            Commit      `json:"last_commit"`
	WorkInProgress        bool        `json:"work_in_progress"`
	URL                   string      `json:"url"`
	Action                string      `json:"action"`
	Assignee              User        `json:"assignee"`
	Author                User        `json:"author"`
	SourceBranchProtected bool        `json:"source_branch_protected"`
	TargetBranchProtected bool        `json:"target_branch_protected"`
	SourceProject         Project
	TargetProject         Project
	Labels                []struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Color       string `json:"color"`
		ProjectID   int    `json:"project_id"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
		Template    bool   `json:"template"`
		Description string `json:"description"`
		Type        string `json:"type"`
		GroupID     int    `json:"group_id"`
	} `json:"labels"`
	Oldrev string `json:"oldrev,omitempty"`
}

type WebhookMergeRequest struct {
	ObjectKind       string           `json:"object_kind"`
	User             User             `json:"user"`
	Project          Project          `json:"project"`
	ObjectAttributes ObjectAttributes `json:"object_attributes"`
}
