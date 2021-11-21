package models

type Data struct {
	HookName   string      `json:"hook_name"`
	Password   string      `json:"password"`
	HookId     int         `json:"hook_id"`
	HookUrl    string      `json:"hook_url"`
	Timestamp  string      `json:"timestamp"`
	Sign       string      `json:"sign"`
	Ref        string      `json:"ref"`
	Before     string      `json:"before"`
	After      string      `json:"after"`
	Created    bool        `json:"created"`
	Deleted    bool        `json:"deleted"`
	Compare    string      `json:"compare"`
	Commits    []Commit    `json:"commits"`
	HeadCommit interface{} `json:"head_commit"`
	Repository Repository  `json:"repository"`
}

type Repository struct {
	Name          string `json:"name"`
	Path          string `json:"path"`
	DefaultBranch string `json:"default_branch"`
	FullName      string `json:"full_name"`
	GitUrl        string `json:"git_url"`
}

type Author struct {
	Time      string `json:"time"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	User_name string `json:"user_name"`
	Url       string `json:"url"`
}

type Commit struct {
	TreeId    string   `json:"tree_id"`
	Message   string   `json:"message"`
	Timestamp string   `json:"timestamp"`
	Url       string   `json:"url"`
	Author    Author   `json:"author"`
	Committer Author   `json:"committer"`
	Modified  []string `json:"modified"`
}
