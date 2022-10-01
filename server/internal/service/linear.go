package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"avocado.com/internal/dao"
	"avocado.com/internal/lib/linear"
	"avocado.com/internal/lib/mErrors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"golang.org/x/oauth2"
)

func (s *Service) ProcessLinearCallback(body []byte, header http.Header) (resp interface{}, err error) {
	sv, err := slack.NewSecretsVerifier(header, s.config.SlackConfig.SigningSecret)
	if err != nil {
		return nil, err
	}

	if _, err := sv.Write(body); err != nil {
		return nil, err
	}

	if err = sv.Ensure(); err != nil {
		return nil, err
	}

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		return
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err = json.Unmarshal([]byte(body), &r)
		if err != nil {
			return
		}
		return gin.H{"challenge": r.Challenge}, err
	} else if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			fmt.Printf("receoved messageEvent %+v\n", ev)
		}
	}

	return
}

func (s *Service) ProcessLinearRedirect(code string, stateToken string) (resp interface{}, err error) {
	conf := &oauth2.Config{
		ClientID:     s.config.Linearconfig.ClientID,
		ClientSecret: s.config.Linearconfig.ClientSecret,
		Scopes:       []string{"read", "issues:create", "comments:create"},
		RedirectURL:  s.config.Linearconfig.RedirectURI,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://linear.app/oauth/authorize",
			TokenURL: "https://api.linear.app/oauth/token",
		},
	}

	ctx := context.Background()
	var token *oauth2.Token
	if token, err = conf.Exchange(ctx, code); err != nil {
		return
	}

	fmt.Println("token")
	fmt.Printf("%+v", token)
	fmt.Println(token.Expiry)
	fmt.Println(token.Expiry.Unix())

	var user *dao.User
	if user, err = s.getUserFromStateToken("linear", stateToken); err != nil {
		return
	} else if user == nil {
		return nil, mErrors.UserNotFoundError
	}

	tokenId := uuid.NewString()
	linearToken := dao.LinearAccessToken{
		Id:           tokenId,
		UserId:       user.Id,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiredAt:    int(token.Expiry.Unix()),
		CreatedAt:    time.Now().Unix(),
	}

	if err = s.dao.SaveLinearAccessToken(&linearToken, s.dao.DB); err != nil {
		fmt.Println("couldn't save linear token")
		fmt.Println(err)
		return
	}

	return
}

func (s *Service) ProcessLinearWebhook(body []byte) (resp interface{}, err error) {
	var payload linear.WebhookPayload
	if payload, err = linear.ParseWebhookPayload(body); err != nil {
		return
	}

	s.c <- func() {
		if payload.Action == "create" && payload.Type == "Comment" {
			if _, err := s.processCreateComment(payload); err != nil {
				fmt.Println(err)
			}
		} else if payload.Action == "update" && payload.Type == "Issue" {
			if _, err := s.processUpdateIssue(payload); err != nil {
				fmt.Println(err)
			}
		} else if payload.Action == "create" && payload.Type == "Reaction" {
			// s.processCommentReaction(payload)
		}
	}

	return
}

func (s *Service) processCreateComment(payload linear.WebhookPayload) (resp interface{}, err error) {
	switch ev := payload.Data.(type) {
	case linear.CommentData:
		{
			// ignore if message was created by us
			var previousSlackMessage *dao.SlackMessage
			if previousSlackMessage, err = s.dao.GetSlackMessageByLinearCommentId(ev.Id, s.dao.DB); err != nil {
				return
			} else if previousSlackMessage != nil {
				return
			}

			// post reply
			var slackMessage *dao.SlackMessage
			if slackMessage, err = s.dao.GetRootSlackMessageByLinearIssueId(ev.Issue.Id, s.dao.DB); err != nil {
				return
			} else if slackMessage == nil {
				return
			}

			// ignore if thread is already deleted
			if slackMessage.Deleted {
				return
			}

			var accessToken *dao.SlackAccessToken
			if accessToken, err = s.dao.GetSlackAccessTokenByTeamId(slackMessage.SlackTeamId, s.dao.DB); err != nil {
				return
			} else if accessToken == nil {
				return
			}

			// generate message body
			msgOptions := []slack.MsgOption{}
			messageBody := ev.Body
			markdownRegex := regexp.MustCompile(`!?\[[^][]+]\((https?://[^()]+)\)`)
			attachments := markdownRegex.FindAllStringSubmatch(messageBody, -1)
			for i := range attachments {
				msgOptions = append(msgOptions, slack.MsgOptionAttachments(slack.Attachment{
					Fallback: attachments[i][0],
					ImageURL: attachments[i][1],
				}))
			}
			msgOptions = append(msgOptions, slack.MsgOptionText(markdownRegex.ReplaceAllString(messageBody, "")+"\n-"+ev.User.Name, true))
			msgOptions = append(msgOptions, slack.MsgOptionTS(slackMessage.SlackTS))

			slackClient := slack.New(accessToken.AccessToken)
			var respChannel, respTimestamp string
			if respChannel, respTimestamp, err = slackClient.PostMessage(slackMessage.SlackChannelId, msgOptions...); err != nil {
				return "", err
			}

			newSlackMessage := dao.SlackMessage{
				Id:              uuid.NewString(),
				SlackTS:         respTimestamp,
				SlackChannelId:  respChannel,
				SlackTeamId:     slackMessage.SlackTeamId,
				LinearIssueId:   ev.Issue.Id,
				LinearCommentId: ev.Id,
				CreatedAt:       time.Now().Unix(),
			}
			if err = s.dao.SaveSlackMessage(&newSlackMessage, s.dao.DB); err != nil {
				return
			}
		}
	}
	return
}

func (s *Service) processUpdateIssue(payload linear.WebhookPayload) (resp interface{}, err error) {
	switch ev := payload.Data.(type) {
	case linear.IssueData:
		{
			var slackMessage *dao.SlackMessage
			if slackMessage, err = s.dao.GetRootSlackMessageByLinearIssueId(ev.Id, s.dao.DB); err != nil {
				return
			} else if slackMessage == nil {
				return
			}

			var accessToken *dao.SlackAccessToken
			if accessToken, err = s.dao.GetSlackAccessTokenByTeamId(slackMessage.SlackTeamId, s.dao.DB); err != nil {
				return
			} else if accessToken == nil {
				return
			}

			slackClient := slack.New(accessToken.AccessToken)
			messageRef := slack.NewRefToMessage(slackMessage.SlackChannelId, slackMessage.SlackTS)
			var existingReactions []slack.ItemReaction
			if existingReactions, err = slackClient.GetReactions(messageRef, slack.GetReactionsParameters{Full: true}); err != nil {
				return
			}

			if ev.State.Name == "In Progress" {
				for _, reaction := range existingReactions {
					if reaction.Name == "eyes" {
						return
					}
				}

				messageRef := slack.NewRefToMessage(slackMessage.SlackChannelId, slackMessage.SlackTS)
				slackClient.RemoveReaction("white_check_mark", messageRef)
				slackClient.AddReaction("eyes", messageRef)
			} else if ev.State.Name == "Done" {
				for _, reaction := range existingReactions {
					if reaction.Name == "white_check_mark" {
						return
					}
				}
				messageRef := slack.NewRefToMessage(slackMessage.SlackChannelId, slackMessage.SlackTS)
				slackClient.RemoveReaction("eyes", messageRef)
				slackClient.AddReaction("white_check_mark", messageRef)
			}
		}
	}
	return
}

func (s *Service) processCommentReaction(payload linear.WebhookPayload) (resp interface{}, err error) {
	fmt.Println("comment reaction")
	switch ev := payload.Data.(type) {
	case linear.ReactionData:
		{
			var slackMessage *dao.SlackMessage
			if slackMessage, err = s.dao.GetSlackMessageByLinearCommentId(ev.Comment.Id, s.dao.DB); err != nil {
				return
			} else if slackMessage == nil {
				return
			}

			fmt.Println("slack message is", slackMessage)

			var accessToken *dao.SlackAccessToken
			if accessToken, err = s.dao.GetSlackAccessTokenByTeamId(slackMessage.SlackTeamId, s.dao.DB); err != nil {
				return
			} else if accessToken == nil {
				fmt.Println("access token is nil")
				return
			}

			slackClient := slack.New(accessToken.AccessToken)
			messageRef := slack.NewRefToMessage(slackMessage.SlackChannelId, slackMessage.SlackTS)
			err := slackClient.AddReaction(ev.Emoji, messageRef)
			fmt.Println("error: ", err)
		}
	}
	return
}
