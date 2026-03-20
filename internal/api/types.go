package api

// Issue is the raw JSON shape returned by the Jira REST API.
// Only fields used by jiru are included.
type Issue struct {
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}

// IssueFields contains the fields nested under an issue.
type IssueFields struct {
	Summary     string      `json:"summary"`
	Description any         `json:"description"` // string (v2) or ADF object (v3)
	Status      NameField   `json:"status"`
	Priority    NameField   `json:"priority"`
	Assignee    UserField   `json:"assignee"`
	Reporter    UserField   `json:"reporter"`
	IssueType   IssueType   `json:"issuetype"`
	Parent      *ParentRef  `json:"parent,omitempty"`
	Labels      []string    `json:"labels"`
	Created     string      `json:"created"`
	Updated     string      `json:"updated"`
	Comment     CommentWrap `json:"comment"`
}

// NameField is a JSON object with a "name" field.
type NameField struct {
	Name string `json:"name"`
}

// UserField holds user information from the API.
type UserField struct {
	Name        string `json:"name"`        // Username (v2 / Server)
	DisplayName string `json:"displayName"` // Full name
	AccountID   string `json:"accountId"`   // Cloud account ID
}

// IssueType holds issue type information.
type IssueType struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Subtask bool   `json:"subtask"`
}

// ParentRef is a reference to a parent issue.
type ParentRef struct {
	Key string `json:"key"`
}

// CommentWrap holds the comment container from the API.
type CommentWrap struct {
	Comments []Comment `json:"comments"`
	Total    int       `json:"total"`
}

// Comment is a single issue comment.
type Comment struct {
	Author UserField `json:"author"`
	Body   any       `json:"body"` // string (v2) or ADF object (v3)
}

// SearchResult is the response from search endpoints.
type SearchResult struct {
	Issues        []*Issue `json:"issues"`
	Total         int      `json:"total"`
	MaxResults    int      `json:"maxResults"`
	StartAt       int      `json:"startAt"`
	IsLast        bool     `json:"isLast"`        // Agile v1 — unreliable
	NextPageToken string   `json:"nextPageToken"` // v3 JQL search
}

// BoardResult is the response from the boards endpoint.
type BoardResult struct {
	MaxResults int     `json:"maxResults"`
	Total      int     `json:"total"`
	IsLast     bool    `json:"isLast"`
	Boards     []Board `json:"values"`
}

// Board represents a Jira board.
type Board struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// SprintResult is the response from the sprints endpoint.
type SprintResult struct {
	MaxResults int      `json:"maxResults"`
	IsLast     bool     `json:"isLast"`
	Sprints    []Sprint `json:"values"`
}

// Sprint represents a Jira sprint.
type Sprint struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"state"` // "active", "closed", "future"
}

// Project represents a Jira project.
type Project struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Type string `json:"projectTypeKey"`
}

// ProjectVersion represents a project release/version.
type ProjectVersion struct {
	Name     string `json:"name"`
	Released bool   `json:"released"`
	Archived bool   `json:"archived"`
}

// User represents a user from the user search endpoint.
type User struct {
	AccountID   string `json:"accountId"`
	DisplayName string `json:"displayName"`
	Active      bool   `json:"active"`
}

// MeResponse is the response from the /myself endpoint.
type MeResponse struct {
	DisplayName string `json:"displayName"`
	Name        string `json:"name"`
}

// CreateResponse is the response from creating an issue.
type CreateResponse struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

// TransitionResponse is the response from the transitions endpoint.
type TransitionResponse struct {
	Transitions []Transition `json:"transitions"`
}

// Transition represents an available status transition.
type Transition struct {
	ID   string    `json:"id"`
	Name string    `json:"name"`
	To   NameField `json:"to"`
}

// CreateMetaResponse is the response from the create metadata endpoint.
type CreateMetaResponse struct {
	Projects []CreateMetaProject `json:"projects"`
}

// CreateMetaProject holds project-level create metadata.
type CreateMetaProject struct {
	Key        string           `json:"key"`
	IssueTypes []CreateMetaType `json:"issuetypes"`
}

// CreateMetaType holds issue type info from create metadata.
type CreateMetaType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CreateMetaFieldsResponse is the response from the create metadata fields endpoint.
type CreateMetaFieldsResponse struct {
	Values []CreateMetaField `json:"values"`
}

// CreateMetaField is a single field definition from create metadata.
type CreateMetaField struct {
	FieldID  string `json:"fieldId"`
	Name     string `json:"name"`
	Required bool   `json:"required"`
	Schema   struct {
		Type  string `json:"type"`
		Items string `json:"items,omitempty"`
	} `json:"schema"`
	AllowedValues []struct {
		Value string `json:"value"`
		Name  string `json:"name"`
	} `json:"allowedValues"`
}

// IssueLinkType represents a type of link between issues.
type IssueLinkType struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Inward  string `json:"inward"`
	Outward string `json:"outward"`
}

// IssueLinkTypesResponse wraps the link types endpoint response.
type IssueLinkTypesResponse struct {
	IssueLinkTypes []IssueLinkType `json:"issueLinkTypes"`
}

// StatusResponse is a single status from the /status endpoint.
type StatusResponse struct {
	Name           string `json:"name"`
	StatusCategory struct {
		Key string `json:"key"`
	} `json:"statusCategory"`
}

// LabelResponse is the response from the /label endpoint.
type LabelResponse struct {
	Values []string `json:"values"`
}

// BoardConfigResponse is the response from the board configuration endpoint.
type BoardConfigResponse struct {
	Filter struct {
		ID string `json:"id"`
	} `json:"filter"`
}

// FilterResponse is the response from the filter endpoint.
type FilterResponse struct {
	JQL string `json:"jql"`
}
