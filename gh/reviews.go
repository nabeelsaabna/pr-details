package gh

import (
	"fmt"
	"time"
)

func (gc GithubClient) ListReviews(owner, repo string, pr int) (ReviewsList, error) {
	u := fmt.Sprintf(
		"%s/repos/%s/%s/pulls/%d/reviews",
		"https://api.github.com", owner, repo, pr,
	)

	res := ReviewsList{}
	err := gc.call(u, &res)
	return res, err
}

type ReviewsList []Review

type Review struct {
	User        User      `json:"user"`
	Body        string    `json:"body"`
	State       string    `json:"state"`
	SubmittedAt time.Time `json:"submitted_at"`
}
