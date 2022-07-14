package gh

import (
	"fmt"
	"time"
)

func (gc GithubClient) ListPullRequests(owner, repo, author string) (SearchResult, error) {
	// TODO add pagination
	u := fmt.Sprintf(
		"%s?q=is:pr+repo:%s/%s+author:%s&per_page=%d&sort=created&order=asc",
		"https://api.github.com/search/issues", owner, repo, author, 100,
	)

	res := SearchResult{}
	err := gc.call(u, &res)
	return res, err
}

type SearchResult struct {
	TotalCount int    `json:"total_count"`
	Incomplete bool   `json:"incomplete_results"`
	Items      []Item `json:"items"`
}

type Item struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	Comments  int       `json:"comments"`
	CreatedAt time.Time `json:"created_at"`
	ClosedAt  time.Time `json:"closed_at"`
	Draft     bool      `json:"draft"`
	User      User      `json:"user"`
}

