package linkedin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Client struct {
	accessToken string
	personID    string
}

func New() *Client {
	return &Client{
		accessToken: os.Getenv("LINKEDIN_ACCESS_TOKEN"),
		personID:    os.Getenv("LINKEDIN_PERSON_ID"),
	}
}

type shareRequest struct {
	Author         string         `json:"author"`
	LifecycleState string         `json:"lifecycleState"`
	SpecificContent specificContent `json:"specificContent"`
	Visibility     visibility     `json:"visibility"`
}

type specificContent struct {
	ShareCommentary shareCommentary `json:"com.linkedin.ugc.ShareContent"`
}

type shareCommentary struct {
	ShareCommentary commentary     `json:"shareCommentary"`
	ShareMediaCategory string      `json:"shareMediaCategory"`
}

type commentary struct {
	Text string `json:"text"`
}

type visibility struct {
	MemberNetworkVisibility string `json:"com.linkedin.ugc.MemberNetworkVisibility"`
}

// PostToLinkedIn posts content to LinkedIn via UGC API
func (c *Client) PostToLinkedIn(content string) error {
	if c.accessToken == "" || c.personID == "" {
		return fmt.Errorf("LINKEDIN_ACCESS_TOKEN or LINKEDIN_PERSON_ID not set")
	}

	payload := shareRequest{
		Author:         fmt.Sprintf("urn:li:person:%s", c.personID),
		LifecycleState: "PUBLISHED",
		SpecificContent: specificContent{
			ShareCommentary: shareCommentary{
				ShareCommentary:    commentary{Text: content},
				ShareMediaCategory: "NONE",
			},
		},
		Visibility: visibility{
			MemberNetworkVisibility: "PUBLIC",
		},
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "https://api.linkedin.com/v2/ugcPosts", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("LinkedIn API error %d: %s", resp.StatusCode, string(b))
	}

	return nil
}

// GetAccessToken returns a new LinkedIn OAuth URL for the user to authorize
func GetAuthURL(clientID, redirectURI string) string {
	return fmt.Sprintf(
		"https://www.linkedin.com/oauth/v2/authorization?response_type=code&client_id=%s&redirect_uri=%s&scope=w_member_social",
		clientID, redirectURI,
	)
}
