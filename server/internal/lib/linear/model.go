package linear

import "encoding/json"

type WebhookPayload struct {
	Action      string           `json:"action"`
	Type        string           `json:"type"`
	CreatedAt   string           `json:"createdAt"`
	DataRaw     *json.RawMessage `json:"data"`
	Data        interface{}
	Url         string       `json:"url"`
	UpdatedFrom *interface{} `json:"updatedFrom"`
}

type CommentData struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Body      string `json:"body"`
	IssueId   string `json:"issueId"`
	UserId    string `json:"userId"`
	Issue     struct {
		Id    string `json:"id"`
		Title string `json:"title"`
	}
	User struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
}

type IssueData struct {
	Id                  string   `json:"id"`
	CreatedAt           string   `json:"createdAt"`
	UpdatedAt           string   `json:"updatedAt"`
	Number              int      `json:"number"`
	Title               string   `json:"title"`
	Description         string   `json:"description"`
	Priority            int      `json:"priority"`
	BoardOrder          int      `json:"boardOrder"`
	SortOrder           float32  `json:"sortOrder"`
	StartedAt           string   `json:"startedAt"`
	TeamId              string   `json:"teamIId"`
	PreviousIdentifiers []string `json:"previousIdentifiers"`
	StateId             string   `json:"stateId"`
	PriorityLabel       string   `json:"priorityLabel"`
	SubscriberIds       []string `json:"subscriberIds"`
	LabelIds            []string `json:"labelIds"`
	State               struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Color string `json:"color"`
		Type  string `json:"type"`
	} `json:"state"`
	Team struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Key  string `json:"key"`
	}
}

type ReactionData struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Emoji     string `json:"emoji"`
	UserId    string `json:"userId"`
	Comment   struct {
		Id   string `json:"id"`
		Body string `json:"body"`
	} `json:"comment"`
	User struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
}
