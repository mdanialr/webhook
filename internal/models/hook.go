package models

// Request hold the most outer scope of incoming JSON from GitHub Webhook
type Request struct {
	Commits []Commit `json:"commits"`
	Branch  string   `json:"ref"`
}

// Commit hold message that identified whether it contains the necessary char or not
type Commit struct {
	Message   string    `json:"message"`
	Committer Committer `json:"committer"`
}

// Committer hold who did the commit
type Committer struct {
	Username string `json:"username"`
}
