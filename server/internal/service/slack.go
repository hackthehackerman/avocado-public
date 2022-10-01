package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"avocado.com/internal/dao"
	"avocado.com/internal/lib/linear"
	"avocado.com/internal/lib/mErrors"
	"avocado.com/internal/lib/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func (s *Service) ProcessSlackEvents(body []byte, header http.Header) (resp interface{}, err error) {
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
		s.c <- func() {
			var slackAccessToken *dao.SlackAccessToken
			if slackAccessToken, err = s.dao.GetSlackAccessTokenByTeamId(eventsAPIEvent.TeamID, s.dao.DB); err != nil {
			} else if slackAccessToken == nil {
			}

			var user *dao.User
			if user, err = s.dao.GetUserById(slackAccessToken.UserId, s.dao.DB); err != nil {
			} else if user == nil {
			}

			var linearAccessToken *dao.LinearAccessToken
			if linearAccessToken, err = s.dao.GetLinearAccessToken(user.Id, s.dao.DB); err != nil {
			}

			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.MessageEvent:
				if ev.SubType == "message_deleted" {
					// close issue when someone delete the message

					var message *dao.SlackMessage
					if message, err = s.dao.GetSlackMessageByThreadTs(ev.PreviousMessage.TimeStamp, s.dao.DB); err != nil {
						return
					} else if message == nil {
						return
					}


					message.Deleted = true
					if err = s.dao.UpdateSlackMessage(message, s.dao.DB); err != nil {
						return
					}

					if message.LinearCommentId == "" {
						lc := linear.NewLinearClient(linearAccessToken.AccessToken)
						lc.UpdateIssueState(message.LinearIssueId, "Canceled")
						lc.AddCommentRaw(message.LinearIssueId, uuid.NewString(), "The original slack thread was deleted")
					} else {
						lc := linear.NewLinearClient(linearAccessToken.AccessToken)
						lc.DeleteComment(message.LinearCommentId)
					}
					return
				}

				if ev.SubType == "" || ev.SubType == "file_share" {
					if ev.ThreadTimeStamp == "" || ev.ThreadTimeStamp == ev.TimeStamp {
						// root message
						slackMessage := dao.SlackMessage{
							Id:                   uuid.NewString(),
							SlackClientMessageId: ev.ClientMsgID,
							SlackThreadTS:        ev.TimeStamp,
							SlackTS:              ev.TimeStamp,
							SlackChannelId:       ev.Channel,
							SlackTeamId:          eventsAPIEvent.TeamID,
							LinearIssueId:        uuid.NewString(),
							LinearCommentId:      "",
							CreatedAt:            time.Now().Unix(),
						}

						if err = s.dao.SaveSlackMessage(&slackMessage, s.dao.DB); err != nil {
							return
						}

						slackClient := slack.New(slackAccessToken.AccessToken)
						params := slack.PermalinkParameters{
							Channel: ev.Channel,
							Ts:      ev.TimeStamp,
						}
						threadURL, err := slackClient.GetPermalink(&params)
						if err != nil {
							return
						}

						var user *slack.User
						if user, err = slackClient.GetUserInfo(ev.User); err != nil {
							return
						}

						issueBody := ev.Text
						if len(ev.Files) > 0 {
							issueBody = issueBody + "\n attachment: \n\n"
							for _, f := range ev.Files {
								issueBody = issueBody + " " + f.Name + " " + f.URLPrivate + "\n"
							}
						}

						if len(ev.Blocks.BlockSet) > 0 {
							issueBody = issueBody + "\n" + parseBlocks(ev.Blocks)
						}
						issueBody = strings.Replace(issueBody, "\n", "\\n", -1)

						lc := linear.NewLinearClient(linearAccessToken.AccessToken)
						lc.CreateNewIssue(slackMessage.LinearIssueId, util.SafeSubtring(ev.Text, 50), issueBody, threadURL, user.RealName, user.Profile.Image48, []string{"generated-from-slack"})

						return
					} else {
						// threaded reply
						// check if this is sent from us already
						var sentSlackMessage *dao.SlackMessage
						if sentSlackMessage, err = s.dao.GetSlackMessageByMessageTs(ev.TimeStamp, s.dao.DB); err != nil {
							return
						} else if sentSlackMessage != nil {
							return
						}

						//
						var parentSlackMessage *dao.SlackMessage
						if parentSlackMessage, err = s.dao.GetSlackMessageByMessageTs(ev.ThreadTimeStamp, s.dao.DB); err != nil {
							return
						} else if parentSlackMessage == nil {
							return
						}

						slackMessage := dao.SlackMessage{
							Id:                   uuid.NewString(),
							SlackClientMessageId: ev.ClientMsgID,
							SlackThreadTS:        ev.ThreadTimeStamp,
							SlackTS:              ev.TimeStamp,
							SlackChannelId:       ev.Channel,
							SlackTeamId:          eventsAPIEvent.TeamID,
							LinearIssueId:        parentSlackMessage.LinearIssueId,
							LinearCommentId:      uuid.NewString(),
							CreatedAt:            time.Now().Unix(),
						}

						if err = s.dao.SaveSlackMessage(&slackMessage, s.dao.DB); err != nil {
							return
						}

						slackClient := slack.New(slackAccessToken.AccessToken)
						var user *slack.User
						if user, err = slackClient.GetUserInfo(ev.User); err != nil {
							return
						}

						commentBody := ev.Text
						if len(ev.Files) > 0 {
							commentBody = commentBody + "\\n attachment: \\n\\n"
							for _, f := range ev.Files {
								commentBody = commentBody + " " + f.Name + " " + f.URLPrivate + "\\n"
							}
						}
						lc := linear.NewLinearClient(linearAccessToken.AccessToken)
						lc.AddComment(parentSlackMessage.LinearIssueId, slackMessage.LinearCommentId, commentBody, user.RealName, user.Profile.Image48)
					}
				}
			case *slackevents.ReactionAddedEvent:
				if ev.Item.Type == "message" {
					var slackMessage *dao.SlackMessage
					if slackMessage, err = s.dao.GetSlackMessageByMessageTs(ev.Item.Timestamp, s.dao.DB); err != nil {
						return
					} else if slackMessage == nil {
						return
					}
					if slackMessage.LinearCommentId == "" {
						// root
						if ev.Reaction == "eyes" {
							lc := linear.NewLinearClient(linearAccessToken.AccessToken)
							lc.UpdateIssueState(slackMessage.LinearIssueId, "In Progress")
						} else if ev.Reaction == "white_check_mark" {
							lc := linear.NewLinearClient(linearAccessToken.AccessToken)
							lc.UpdateIssueState(slackMessage.LinearIssueId, "Done")
						} else if ev.Reaction == "x" {
							lc := linear.NewLinearClient(linearAccessToken.AccessToken)
							lc.UpdateIssueState(slackMessage.LinearIssueId, "Canceled")
						}
					} else {
						// comment
						lc := linear.NewLinearClient(linearAccessToken.AccessToken)
						lc.AddReaction(slackMessage.LinearCommentId, ev.Reaction)
					}
				}
			}
		}

	}

	return
}

func (s *Service) ProcessSlackRedirect(code string, stateToken string) (resp interface{}, err error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	oauthResp, err := slack.GetOAuthV2Response(client, s.config.SlackConfig.ClientID, s.config.SlackConfig.ClientSecret, code, s.config.SlackConfig.RedirectURI)
	if err != nil {
		return nil, err
	}


	if !oauthResp.Ok {
		return nil, mErrors.Error{Code: http.StatusUnauthorized, Msg: "Failed to retrieve oauth token"}
	}

	var user *dao.User
	if user, err = s.getUserFromStateToken("slack", stateToken); err != nil {
		return
	} else if user == nil {
		return nil, mErrors.UserNotFoundError
	}

	id := uuid.NewString()
	ts := time.Now().Unix()
	accessToken := &dao.SlackAccessToken{
		Id:           id,
		UserId:       user.Id,
		TeamId:       oauthResp.Team.ID,
		AccessToken:  oauthResp.AccessToken,
		RefreshToken: &oauthResp.RefreshToken,
		ExpiredIn:    oauthResp.ExpiresIn,
		CreatedAt:    ts,
	}
	if err = s.dao.SaveSlackAccessToken(accessToken, s.dao.DB); err != nil {
		return
	}

	return
}

func parseBlocks(blocks slack.Blocks) string {
	var buffer bytes.Buffer
	for _, b := range blocks.BlockSet {
		buffer.WriteString(parseBlock(b))
	}
	buffer.WriteString("\n")
	return buffer.String()
}

func parseBlock(block slack.Block) string {
	var buffer bytes.Buffer

	switch block.BlockType() {
	case slack.MBTSection:
		b := block.(*slack.SectionBlock)
		buffer.WriteString("\n")
		buffer.WriteString(b.Text.Text)
		for _, f := range b.Fields {
			buffer.WriteString("\n\n")
			buffer.WriteString(f.Text)
		}
	case slack.MBTDivider:
	case slack.MBTImage:
		b := block.(*slack.ImageBlock)
		buffer.WriteString("\n")
		buffer.WriteString(fmt.Sprintf("![%s](%s)", b.AltText, b.ImageURL))
	case slack.MBTAction:
	case slack.MBTContext:
	case slack.MBTFile:
	case slack.MBTInput:
	case slack.MBTHeader:
	case slack.MBTRichText:
	}

	return buffer.String()
}
