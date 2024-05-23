package gitlabtype

type GitLabEvent string

const (
	PushEvent              GitLabEvent = "Push Hook"
	TagPushEvent           GitLabEvent = "Tag Push Hook"
	IssueEvent             GitLabEvent = "Issue Hook"
	ConfidentialIssueEvent GitLabEvent = "Confidential Issue Hook"
	MergeRequestEvent      GitLabEvent = "Merge Request Hook"
	WikiPageEvent          GitLabEvent = "Wiki Page Hook"
	PipelineEvent          GitLabEvent = "Pipeline Hook"
	BuildEvent             GitLabEvent = "Build Hook"
	NoteEvent              GitLabEvent = "Note Hook"
	ConfidentialNoteEvent  GitLabEvent = "Confidential Note Hook"
	JobEvent               GitLabEvent = "Job Hook"
	DeploymentEvent        GitLabEvent = "Deployment Hook"
	ReleaseEvent           GitLabEvent = "Release Hook"
)
