package model

type UserSettingResponse struct {
	UserId            string `json:"userId"`
	SlackRedirectURI  string `json:"slackRedirectURI"`
	SlackConnected    bool   `json:"slackConnected"`
	LinearRedirectURI string `json:"linearRedirectURI"`
	LinearConnected   bool   `json:"linearConnected"`
}
