package github

import "github.com/da4nik/ssci/types"

// User represents user in github webhook
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Repository represents repository section in github hook
type Repository struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Owner         User   `json:"owner"`
	Private       bool   `json:"private"`
	Description   string `json:"description"`
	Fork          bool   `json:"fork"`
	GitURL        string `json:"git_url"`
	SSHURL        string `json:"ssh_url"`
	CloneURL      string `json:"clone_url"`
	DefaultBranch string `json:"default_branch"`
	MasterBranch  string `json:"master_branch"`
}

// Commit represents commit in github webhook
type Commit struct {
	ID        string   `json:"id"`
	Message   string   `json:"message"`
	Timestamp string   `json:"timestamp"`
	Author    User     `json:"author"`
	Committer User     `json:"committer"`
	Added     []string `json:"added"`
	Removed   []string `json:"removed"`
	Modified  []string `json:"modified"`
}

// PushEvent push event payload
type PushEvent struct {
	Repository Repository `json:"repository"`
	Pusher     User       `json:"pusher"`
	HeadCommit Commit     `json:"head_commit"`
	Commits    []Commit   `json:"commits"`
}

// Notification coverts github event to Notification
func (pe PushEvent) Notification() types.Notification {
	return types.Notification{
		Name:     pe.Repository.FullName,
		CloneURL: pe.Repository.CloneURL,
	}
}
