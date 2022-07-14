package gh

import (
	"fmt"
	"time"
)

func (gc GithubClient) ListReviewComments(owner, repo string, pr int) (ReviewsCommentsList, error){
	u := fmt.Sprintf(
		"%s/repos/%s/%s/pulls/%d/comments",
		"https://api.github.com", owner, repo, pr,
	)

	res := ReviewsCommentsList{}
	err := gc.call(u, &res)
	return res, err
}

type ReviewsCommentsList []ReviewsComments

type ReviewsComments struct {
	User              User      `json:"user"`
	Body              string    `json:"body"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	AuthorAssociation string    `json:"author_association"`
}
