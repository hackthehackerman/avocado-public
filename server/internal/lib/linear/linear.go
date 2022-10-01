package linear

import (
	"context"
	"fmt"
	"log"

	"github.com/machinebox/graphql"
	"golang.org/x/oauth2"
)

type LinearClient struct {
	gc            *graphql.Client
	accessToken   string
	teamIdMapping map[string]string
	stateMapping  map[string]map[string]string
	labelMapping  map[string]map[string]string
}

func NewLinearClient(accessToken string) *LinearClient {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := graphql.NewClient("https://api.linear.app/graphql", graphql.WithHTTPClient(httpClient))
	ret := LinearClient{
		gc:            client,
		accessToken:   accessToken,
		teamIdMapping: make(map[string]string),
		stateMapping:  make(map[string]map[string]string),
		labelMapping:  make(map[string]map[string]string),
	}
	client.Log = func(s string) { log.Println(s) }
	return &ret
}

func (l *LinearClient) CreateNewIssue(id, title, description, threadDeepLink, username, avatarurl string, labelNames []string) (err error) {
	teamId, err := l.teamId()
	if err != nil {
		return
	}

	labelIds := []string{}
	for _, labelName := range labelNames {
		labelId, err := l.labelId(labelName)
		if err != nil {
			return err
		}
		labelIds = append(labelIds, "\""+labelId+"\"")
	}

	req := graphql.NewRequest(fmt.Sprintf(`
	mutation IssueCreate {
				issueCreate(
				  input: {
					id: "%s"
					title: "%s"
					description: "%s"
					teamId: "%s"
					createAsUser: "%s"
					displayIconUrl: "%s"
					labelIds: %s
				  }
				) 
			  {
				success
				issue {
				id
				title
				}
			  }
			}
	`, id, title, description, teamId, username, avatarurl, labelIds))
	type Success struct {
		Issue struct {
			Id    string
			title string
		}
	}
	var respData Success
	if err = l.gc.Run(context.Background(), req, &respData); err != nil {
		return
	}

	// make attachment
	type AttachmentSuccess struct {
		Attachment struct {
			Id string
		}
	}
	var attachmentRespData AttachmentSuccess
	req = graphql.NewRequest(fmt.Sprintf(`
	mutation{
		attachmentCreate(input:{
		  issueId: "%s"
		  title: "Open in slack"
		  url: "%s"
		  iconUrl: "https://a.slack-edge.com/80588/marketing/img/icons/icon_slack_hash_colored.png"
		}){
		  success
		  attachment{
			id
		  }
		}
	  }
	`, id, threadDeepLink))
	if err = l.gc.Run(context.Background(), req, &attachmentRespData); err != nil {
		return
	}

	return
}

func (l *LinearClient) AddComment(issueId, commentId, body, username, avatarurl string) (err error) {
	req := graphql.NewRequest(fmt.Sprintf(`
	mutation CommentCreate {
				commentCreate(
				  input: {
					issueId: "%s"
					id: "%s"
					body: "%s"
					createAsUser: "%s"
					displayIconUrl: "%s"
				  }
				) 
			  {
				success
				comment {
				id
				}
			  }
			}
	`, issueId, commentId, body, username, avatarurl))

	type Success struct {
		Comment struct {
			Id string
		}
	}
	var respData Success
	if err = l.gc.Run(context.Background(), req, &respData); err != nil {
		return
	}
	return
}

func (l *LinearClient) AddCommentRaw(issueId, commentId, body string) (err error) {
	req := graphql.NewRequest(fmt.Sprintf(`
	mutation CommentCreate {
				commentCreate(
				  input: {
					issueId: "%s"
					id: "%s"
					body: "%s"
				  }
				) 
			  {
				success
				comment {
				id
				}
			  }
			}
	`, issueId, commentId, body))

	type Success struct {
		Comment struct {
			Id string
		}
	}
	var respData Success
	if err = l.gc.Run(context.Background(), req, &respData); err != nil {
		return
	}
	return
}

func (l *LinearClient) AddReaction(commentId, emoji string) (err error) {
	req := graphql.NewRequest(fmt.Sprintf(`
	mutation ReactionCreate {
		reactionCreate(
				  input: {
					commentId: "%s"
					emoji: "%s"
				  }
				) 
			  {
				success
				reaction {
				id
				}
			  }
			}
	`, commentId, emoji))

	if err = l.gc.Run(context.Background(), req, nil); err != nil {
		return
	}
	return
}

func (l *LinearClient) UpdateIssueState(issueId, state string) (err error) {
	var states map[string]string
	if states, err = l.states(); err != nil {
		return
	}

	req := graphql.NewRequest(fmt.Sprintf(`
	mutation IssueUpdate {
		issueUpdate(
				  input: {
					stateId: "%s"
				  }
				  id: "%s"
				)
			  {
				success
				issue {
				id
				}
			  }
			}
	`, states[state], issueId))

	if err = l.gc.Run(context.Background(), req, nil); err != nil {
		return
	}
	return
}

func (l *LinearClient) DeleteComment(commentId string) (err error) {
	req := graphql.NewRequest(fmt.Sprintf(`
	mutation CommentDelete {
		commentDelete(
			id: "%s"
		) 
	  {
		success
	  }
	}
	`, commentId))

	if err = l.gc.Run(context.Background(), req, nil); err != nil {
		return
	}
	return
}

func (l *LinearClient) teamId() (teamId string, err error) {
	if teamId, ok := l.teamIdMapping[l.accessToken]; ok {
		return teamId, nil
	}

	req := graphql.NewRequest(
		`query Teams {
			teams {
			  nodes {
				id
				name
			  }
			}
		  }`)
	type Resp struct {
		Teams struct {
			Nodes []struct {
				Id   string
				Name string
			}
		}
	}
	var respData Resp
	if err = l.gc.Run(context.Background(), req, &respData); err != nil {
		return
	}
	l.teamIdMapping[l.accessToken] = respData.Teams.Nodes[0].Id
	return respData.Teams.Nodes[0].Id, nil
}

func (l *LinearClient) states() (states map[string]string, err error) {
	if stateMapping, ok := l.stateMapping[l.accessToken]; ok {
		return stateMapping, nil
	}

	req := graphql.NewRequest(
		`query WorkflowStates($filter: WorkflowStateFilter) {
			workflowStates(filter: $filter){
			  nodes {
				id
				name
			  }
			}
		  }`)
	type Resp struct {
		WorkflowState struct {
			Nodes []struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			} `json:"nodes"`
		} `json:"workflowStates"`
	}
	var respData Resp
	if err = l.gc.Run(context.Background(), req, &respData); err != nil {
		return
	}

	l.stateMapping[l.accessToken] = make(map[string]string)
	for _, node := range respData.WorkflowState.Nodes {
		l.stateMapping[l.accessToken][node.Name] = node.Id
	}

	return l.stateMapping[l.accessToken], nil
}

func (l *LinearClient) labelId(labelName string) (id string, err error) {
	if _, ok := l.labelMapping[l.accessToken]; !ok {
		l.reloadLabels()
	}

	if id, ok := l.labelMapping[l.accessToken][labelName]; ok {
		return id, nil
	}

	req := graphql.NewRequest(
		fmt.Sprintf(`mutation IssueLabelCreate{
			issueLabelCreate(
			  replaceTeamLabels: false
			  input: {
				name: "%s"
			  }
			){
			  success
			  issueLabel{
				id
				name
			  }
			}
		  }`, labelName))

	if err = l.gc.Run(context.Background(), req, nil); err != nil {
		return
	}

	l.reloadLabels()

	return l.labelMapping[l.accessToken][labelName], nil
}

func (l *LinearClient) reloadLabels() (err error) {
	req := graphql.NewRequest(
		`query IssueLabels($filter: IssueLabelFilter) {
			issueLabels(filter: $filter){
			  nodes {
				id
				name
			  }
			}
		  }`)
	type Resp struct {
		IssueLabels struct {
			Nodes []struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			} `json:"nodes"`
		} `json:"issueLabels"`
	}
	var respData Resp
	if err = l.gc.Run(context.Background(), req, &respData); err != nil {
		return
	}

	newmapping := make(map[string]string)
	for _, node := range respData.IssueLabels.Nodes {
		newmapping[node.Name] = node.Id
	}
	l.labelMapping[l.accessToken] = newmapping
	return
}
