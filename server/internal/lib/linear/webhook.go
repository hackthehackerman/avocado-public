package linear

import "encoding/json"

func ParseWebhookPayload(body []byte) (payload WebhookPayload, err error) {
	if err = json.Unmarshal(body, &payload); err != nil {
		return
	}

	if payload.Type == "Comment" {
		var comment CommentData
		if err = json.Unmarshal(*payload.DataRaw, &comment); err != nil {
			return
		}
		payload.Data = comment
	} else if payload.Type == "Issue" {
		var issue IssueData
		if err = json.Unmarshal(*payload.DataRaw, &issue); err != nil {
			return
		}
		payload.Data = issue
	} else if payload.Type == "Reaction" {
		var issue ReactionData
		if err = json.Unmarshal(*payload.DataRaw, &issue); err != nil {
			return
		}
		payload.Data = issue
	}

	return
}
