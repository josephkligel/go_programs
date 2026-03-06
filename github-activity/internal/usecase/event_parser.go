package usecase

import (
	"fmt"
)

// Helper functions for safe type assertions
func getString(m map[string]any, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getMap(m map[string]any, key string) map[string]any {
	if val, ok := m[key]; ok {
		if innerMap, ok := val.(map[string]any); ok {
			return innerMap
		}
	}
	return nil
}

func getArray(m map[string]any, key string) []any {
	if val, ok := m[key]; ok {
		if arr, ok := val.([]any); ok {
			return arr
		}
	}
	return nil
}

func ParseSingleAction(event map[string]any) string {
	if action, ok := event["type"].(string); ok {
		var actionString string
		payload := getMap(event, "payload")
		
		switch action {
		case "CommitCommentEvent":
			if payload != nil {
				comment := getMap(payload, "comment")
				actionString = fmt.Sprintf("- Commented on commit: %v", getString(comment, "body"))
			} else {
				actionString = "- Commented on commit"
			}
			return actionString
		case "CreateEvent":
			if payload != nil {
				actionString = fmt.Sprintf("- Created %v called %v", getString(payload, "ref_type"), getString(payload, "ref"))
			} else {
				actionString = "- Created something"
			}
			return actionString
		case "DeleteEvent":
			if payload != nil {
				actionString = fmt.Sprintf("- Deleted %v called %v", getString(payload, "ref_type"), getString(payload, "ref"))
			} else {
				actionString = "- Deleted something"
			}
			return actionString
		case "DiscussionEvent":
			if payload != nil {
				discussion := getMap(payload, "discussion")
				actionString = fmt.Sprintf("- Started discussion: %v", getString(discussion, "title"))
			} else {
				actionString = "- Started discussion"
			}
			return actionString
		case "ForkEvent":
			if payload != nil {
				forkee := getMap(payload, "forkee")
				actionString = fmt.Sprintf("- Forked repository to %v", getString(forkee, "full_name"))
			} else {
				actionString = "- Forked repository"
			}
			return actionString
		case "GollumEvent":
			if payload != nil {
				pages := getArray(payload, "pages")
				if len(pages) > 0 {
					if page, ok := pages[0].(map[string]any); ok {
						actionString = fmt.Sprintf("- Edited wiki page: %v", getString(page, "page_name"))
					} else {
						actionString = "- Edited wiki page"
					}
				} else {
					actionString = "- Edited wiki page"
				}
			} else {
				actionString = "- Edited wiki page"
			}
			return actionString
		case "IssueCommentEvent":
			if payload != nil {
				comment := getMap(payload, "comment")
				actionString = fmt.Sprintf("- Commented on issue: %v", getString(comment, "body"))
			} else {
				actionString = "- Commented on issue"
			}
			return actionString
		case "IssuesEvent":
			if payload != nil {
				issue := getMap(payload, "issue")
				actionString = fmt.Sprintf("- %v issue: %v", getString(payload, "action"), getString(issue, "title"))
			} else {
				actionString = "- Modified issue"
			}
			return actionString
		case "MemberEvent":
			if payload != nil {
				member := getMap(payload, "member")
				actionString = fmt.Sprintf("- %v member: %v", getString(payload, "action"), getString(member, "login"))
			} else {
				actionString = "- Modified member"
			}
			return actionString
		case "PublicEvent":
			if payload != nil {
				repository := getMap(payload, "repository")
				actionString = fmt.Sprintf("- Made repository public: %v", getString(repository, "full_name"))
			} else {
				actionString = "- Made repository public"
			}
			return actionString
		case "PullRequestEvent":
			if payload != nil {
				pr := getMap(payload, "pull_request")
				actionString = fmt.Sprintf("- %v pull request: %v", getString(payload, "action"), getString(pr, "title"))
			} else {
				actionString = "- Modified pull request"
			}
			return actionString
		case "PullRequestReviewEvent":
			if payload != nil {
				pr := getMap(payload, "pull_request")
				actionString = fmt.Sprintf("- Reviewed pull request: %v", getString(pr, "title"))
			} else {
				actionString = "- Reviewed pull request"
			}
			return actionString
		case "PullRequestReviewCommentEvent":
			if payload != nil {
				comment := getMap(payload, "comment")
				actionString = fmt.Sprintf("- Commented on pull request: %v", getString(comment, "body"))
			} else {
				actionString = "- Commented on pull request"
			}
			return actionString
		case "PushEvent":
			if payload != nil {
				repository := getMap(payload, "repository")
				actionString = fmt.Sprintf("- Pushed to %v at %v", getString(payload, "ref"), getString(repository, "full_name"))
			} else {
				actionString = "- Pushed code"
			}
			return actionString
		case "ReleaseEvent":
			if payload != nil {
				release := getMap(payload, "release")
				actionString = fmt.Sprintf("- %v release: %v", getString(payload, "action"), getString(release, "name"))
			} else {
				actionString = "- Modified release"
			}
			return actionString
		case "WatchEvent":
			if payload != nil {
				repo := getMap(event, "repo")
				actionString = fmt.Sprintf("- %v watching repository: %v", getString(payload, "action"), getString(repo, "name"))
			} else {
				actionString = "- Watching repository"
			}
			return actionString
		default:
			actionString = "Unknown action: " + action
			return actionString
		}
	} else {
		return "Could not parse action type"
	}
}

func ParseActions(data []map[string]any) []string {
	var actions []string
	var singleAction string
	for _, event := range data {
		singleAction = ParseSingleAction(event)
		actions = append(actions, singleAction)
	}
	return actions
}
