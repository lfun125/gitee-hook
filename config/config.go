package config

type Model struct {
	Listen              string `yaml:"listen"`
	JenkinsUrl          string `yaml:"jenkins_url"`
	JenkinsUser         string `yaml:"jenkins_user"`
	JenkinsUserToken    string `yaml:"jenkins_user_token"`
	JenkinsProjectToken string `yaml:"jenkins_project_token"`
	GiteeSecret         string `yaml:"gitee_secret"`
	Repository          map[string]struct {
		Branches map[string]string `json:"branches"` // branch -> jenkins job
	} `yaml:"repository"`
}

type Push struct {
	Ref        string `json:"ref"`
	Repository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
	} `json:"repository"`
}
