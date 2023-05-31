package main

import "time"

type WorkflowJobConclusion string
type WorkflowJobStatus string

const (
	ConclusionSuccess        WorkflowJobConclusion = "success"
	ConclusionFailure        WorkflowJobConclusion = "failure"
	ConclusionNull           WorkflowJobConclusion = "null"
	ConclusionSkipped        WorkflowJobConclusion = "skipped"
	ConclusionCancelled      WorkflowJobConclusion = "cancelled"
	ConclusionActionRequired WorkflowJobConclusion = "action_required"
	ConclusionNeutral        WorkflowJobConclusion = "neutral"
	ConclusionTimedOut       WorkflowJobConclusion = "timed_out"

	StatusQueued     WorkflowJobStatus = "queued"
	StatusInProgress WorkflowJobStatus = "in_progress"
	StatusCompleted  WorkflowJobStatus = "completed"
	StatusWaiting    WorkflowJobStatus = "waiting"
)

func (s WorkflowJobStatus) String() string {
	return string(s)
}

func (s WorkflowJobConclusion) String() string {
	return string(s)
}

type WorkflowJobWebhook struct {
	Action      string `json:"action"`
	WorkflowJob struct {
		ID           int64     `json:"id"`
		RunID        int64     `json:"run_id"`
		WorkflowName string    `json:"workflow_name"`
		HeadBranch   string    `json:"head_branch"`
		RunURL       string    `json:"run_url"`
		RunAttempt   int       `json:"run_attempt"`
		NodeID       string    `json:"node_id"`
		HeadSha      string    `json:"head_sha"`
		URL          string    `json:"url"`
		HTMLURL      string    `json:"html_url"`
		Status       string    `json:"status"`
		Conclusion   string    `json:"conclusion"`
		CreatedAt    time.Time `json:"created_at"`
		StartedAt    time.Time `json:"started_at"`
		CompletedAt  time.Time `json:"completed_at"`
		Name         string    `json:"name"`
		Steps        []struct {
			Name        string    `json:"name"`
			Status      string    `json:"status"`
			Conclusion  string    `json:"conclusion"`
			Number      int       `json:"number"`
			StartedAt   time.Time `json:"started_at"`
			CompletedAt time.Time `json:"completed_at"`
		} `json:"steps"`
		CheckRunURL     string   `json:"check_run_url"`
		Labels          []string `json:"labels"`
		RunnerID        int      `json:"runner_id"`
		RunnerName      string   `json:"runner_name"`
		RunnerGroupID   int      `json:"runner_group_id"`
		RunnerGroupName string   `json:"runner_group_name"`
	} `json:"workflow_job"`
	Repository struct {
		ID       int    `json:"id"`
		NodeID   string `json:"node_id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Private  bool   `json:"private"`
		Owner    struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"owner"`
		HTMLURL                  string        `json:"html_url"`
		Description              string        `json:"description"`
		Fork                     bool          `json:"fork"`
		URL                      string        `json:"url"`
		ForksURL                 string        `json:"forks_url"`
		KeysURL                  string        `json:"keys_url"`
		CollaboratorsURL         string        `json:"collaborators_url"`
		TeamsURL                 string        `json:"teams_url"`
		HooksURL                 string        `json:"hooks_url"`
		IssueEventsURL           string        `json:"issue_events_url"`
		EventsURL                string        `json:"events_url"`
		AssigneesURL             string        `json:"assignees_url"`
		BranchesURL              string        `json:"branches_url"`
		TagsURL                  string        `json:"tags_url"`
		BlobsURL                 string        `json:"blobs_url"`
		GitTagsURL               string        `json:"git_tags_url"`
		GitRefsURL               string        `json:"git_refs_url"`
		TreesURL                 string        `json:"trees_url"`
		StatusesURL              string        `json:"statuses_url"`
		LanguagesURL             string        `json:"languages_url"`
		StargazersURL            string        `json:"stargazers_url"`
		ContributorsURL          string        `json:"contributors_url"`
		SubscribersURL           string        `json:"subscribers_url"`
		SubscriptionURL          string        `json:"subscription_url"`
		CommitsURL               string        `json:"commits_url"`
		GitCommitsURL            string        `json:"git_commits_url"`
		CommentsURL              string        `json:"comments_url"`
		IssueCommentURL          string        `json:"issue_comment_url"`
		ContentsURL              string        `json:"contents_url"`
		CompareURL               string        `json:"compare_url"`
		MergesURL                string        `json:"merges_url"`
		ArchiveURL               string        `json:"archive_url"`
		DownloadsURL             string        `json:"downloads_url"`
		IssuesURL                string        `json:"issues_url"`
		PullsURL                 string        `json:"pulls_url"`
		MilestonesURL            string        `json:"milestones_url"`
		NotificationsURL         string        `json:"notifications_url"`
		LabelsURL                string        `json:"labels_url"`
		ReleasesURL              string        `json:"releases_url"`
		DeploymentsURL           string        `json:"deployments_url"`
		CreatedAt                time.Time     `json:"created_at"`
		UpdatedAt                time.Time     `json:"updated_at"`
		PushedAt                 time.Time     `json:"pushed_at"`
		GitURL                   string        `json:"git_url"`
		SSHURL                   string        `json:"ssh_url"`
		CloneURL                 string        `json:"clone_url"`
		SvnURL                   string        `json:"svn_url"`
		Homepage                 string        `json:"homepage"`
		Size                     int           `json:"size"`
		StargazersCount          int           `json:"stargazers_count"`
		WatchersCount            int           `json:"watchers_count"`
		Language                 string        `json:"language"`
		HasIssues                bool          `json:"has_issues"`
		HasProjects              bool          `json:"has_projects"`
		HasDownloads             bool          `json:"has_downloads"`
		HasWiki                  bool          `json:"has_wiki"`
		HasPages                 bool          `json:"has_pages"`
		HasDiscussions           bool          `json:"has_discussions"`
		ForksCount               int           `json:"forks_count"`
		MirrorURL                interface{}   `json:"mirror_url"`
		Archived                 bool          `json:"archived"`
		Disabled                 bool          `json:"disabled"`
		OpenIssuesCount          int           `json:"open_issues_count"`
		License                  interface{}   `json:"license"`
		AllowForking             bool          `json:"allow_forking"`
		IsTemplate               bool          `json:"is_template"`
		WebCommitSignoffRequired bool          `json:"web_commit_signoff_required"`
		Topics                   []interface{} `json:"topics"`
		Visibility               string        `json:"visibility"`
		Forks                    int           `json:"forks"`
		OpenIssues               int           `json:"open_issues"`
		Watchers                 int           `json:"watchers"`
		DefaultBranch            string        `json:"default_branch"`
	} `json:"repository"`
	Organization struct {
		Login            string `json:"login"`
		ID               int    `json:"id"`
		NodeID           string `json:"node_id"`
		URL              string `json:"url"`
		ReposURL         string `json:"repos_url"`
		EventsURL        string `json:"events_url"`
		HooksURL         string `json:"hooks_url"`
		IssuesURL        string `json:"issues_url"`
		MembersURL       string `json:"members_url"`
		PublicMembersURL string `json:"public_members_url"`
		AvatarURL        string `json:"avatar_url"`
		Description      string `json:"description"`
	} `json:"organization"`
	Enterprise struct {
		ID          int       `json:"id"`
		Slug        string    `json:"slug"`
		Name        string    `json:"name"`
		NodeID      string    `json:"node_id"`
		AvatarURL   string    `json:"avatar_url"`
		Description string    `json:"description"`
		WebsiteURL  string    `json:"website_url"`
		HTMLURL     string    `json:"html_url"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	} `json:"enterprise"`
	Sender struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"sender"`
	Installation struct {
		ID     int    `json:"id"`
		NodeID string `json:"node_id"`
	} `json:"installation"`
}
