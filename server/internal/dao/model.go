package dao

type tableName string

const (
	httpRequest      = tableName("http_request")
	httpResponse     = tableName("http_response")
	slackAccessToken = tableName("slack_access_token")
	user             = tableName("user")
)

type HttpRequest struct {
	Id               string  `db:"id"`
	URI              *string `db:"uri"`
	Header           *string `db:"header"`
	Body             *string `db:"body"`
	OriginatedFromUs bool    `db:"originated_from_us"`
	CreatedAt        int64   `db:"created_at"`
}

type HttpResponse struct {
	Id               string  `db:"id"`
	Header           *string `db:"header"`
	Body             *string `db:"body"`
	Status           *int
	OriginatedFromUs bool    `db:"originated_from_us"`
	RequestID        *string `db:"request_id"`
	CreatedAt        int64   `db:"created_at"`
}

type SlackAccessToken struct {
	Id           string  `db:"id"`
	UserId       string  `db:"user_id"`
	TeamId       string  `db:"team_id"`
	AccessToken  string  `db:"access_token"`
	RefreshToken *string `db:"refresh_token"`
	ExpiredIn    int     `db:"expired_in"`
	CreatedAt    int64   `db:"created_at"`
}

type LinearAccessToken struct {
	Id           string `db:"id"`
	UserId       string `db:"user_id"`
	AccessToken  string `db:"access_token"`
	RefreshToken string `db:"refresh_token"`
	ExpiredAt    int    `db:"expired_at"`
	CreatedAt    int64  `db:"created_at"`
}

type User struct {
	Id        string `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
	CreatedAt int64  `db:"created_at"`
	UpdatedAt int64  `db:"updated_at"`
}

type UserSession struct {
	Id        string `db:"id"`
	UserId    string `db:"user_id"`
	Token     string `db:"token"`
	ExpiredAt int64  `db:"expired_at"`
	CreatedAt int64  `db:"created_at"`
	UpdatedAt int64  `db:"updated_at"`
	Deleted   bool   `db:"deleted"`
}

type OauthStateToken struct {
	Id        string `db:"id"`
	Service   string `db:"service"`
	UserId    string `db:"user_id"`
	Token     string `db:"token"`
	CreatedAt int64  `db:"created_at"`
}

type SlackMessage struct {
	Id                   string `db:"id"`
	SlackClientMessageId string `db:"slack_client_message_id"`
	SlackThreadTS        string `db:"slack_thread_ts"`
	SlackTS              string `db:"slack_ts"`
	SlackChannelId       string `db:"slack_channel_id"`
	SlackTeamId          string `db:"slack_team_id"`
	LinearIssueId        string `db:"linear_issue_id"`
	LinearCommentId      string `db:"linear_comment_id"`
	CreatedAt            int64  `db:"created_at"`
	Deleted              bool   `db:"deleted"`
}
