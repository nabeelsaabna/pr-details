package runner

import (
	"fmt"
	"github.com/nabeelys/pr-details/gh"
	"io"
	"os"
	"sort"
	"strings"
)

var client *gh.GithubClient = nil

func SetClient(cl *gh.GithubClient) {
	client = cl
}

func Run(owner, repository, author string) (err error) {
	if client == nil {
		err = fmt.Errorf("missing client")
		return
	}

	fmt.Printf("Starting processing: %s/%s for %s\n", owner, repository, author)

	fileName := fmt.Sprintf("prs_review_%s_%s_%s.csv", owner, repository, author)
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println("error creating file ", fileName)
		return err
	}
	res, err := client.ListPullRequests(owner, repository, author)

	if err != nil {
		fmt.Println("error calling ListPullRequests")
		return err
	}

	fmt.Printf("Total Count: %d\n", res.TotalCount)
	fmt.Printf("Incompelte: %t\n", res.Incomplete)
	fmt.Printf("Items: %d\n", len(res.Items))

	var lines []prData
	for _, v := range res.Items {
		line := handlePullRequest(v, owner, repository)
		lines = append(lines, line)
	}

	sorted := sortReviewers(lines)
	printToFile(lines, sorted, f)
	return
}

func handlePullRequest(it gh.Item, owner, repository string) prData {
	lcTitle := strings.ToLower(it.Title)
	var prType = " - "
	if strings.Contains(lcTitle, "fix") || strings.Contains(lcTitle, "bug") {
		prType = "bug"
	} else if strings.Contains(lcTitle, "refactor") {
		prType = "refactor"
	} else if strings.Contains(lcTitle, "feat") {
		prType = "feature"

	}

	closedAt := " - "
	duration := " - "
	if it.CreatedAt.Before(it.ClosedAt) {
		duration = fmt.Sprintf("%d", (it.ClosedAt.Unix()-it.CreatedAt.Unix())/(60*60*24))
		closedAt = it.ClosedAt.Format("2006-01-02")
	}

	reviewers, err := aggregateReviewsAndComments(it.Number, owner, repository)
	if err != nil {
		fmt.Printf("error aggregateReviewsAndComments: %v\n", err)
		return prData{}
	}

	return prData{
		number:    it.Number,
		title:     it.Title,
		prType:    prType,
		state:     it.State,
		created:   it.CreatedAt.Format("2006-01-02"),
		closed:    closedAt,
		duration:  duration,
		comments:  it.Comments,
		reviewers: reviewers,
	}
}

func aggregateReviewsAndComments(prNumber int, owner, repository string) (map[string]int, error) {
	result := map[string]int{}
	reviews, err := client.ListReviews(owner, repository, prNumber)
	if err != nil {
		fmt.Printf("error ListReviews: %v\n", err)
		return result, err
	}

	comments, err := client.ListReviewComments(owner, repository, prNumber)
	if err != nil {
		fmt.Printf("error ListReviewComments: %v\n", err)
		return result, err
	}

	for _, r := range reviews {
		count, found := result[r.User.Login]
		if !found {
			count = 0
		}
		result[r.User.Login] = count + 1
	}

	for _, c := range comments {
		count, found := result[c.User.Login]
		if !found {
			count = 0
		}
		result[c.User.Login] = count + 1
	}

	return result, err
}

func sortReviewers(lines []prData) []string {
	namesMap := map[string]bool{}
	for _, l := range lines {
		for name, _ := range l.reviewers {
			namesMap[name] = true
		}
	}

	var namesList []string
	for name := range namesMap {
		namesList = append(namesList, name)
	}

	sort.Strings(namesList)
	return namesList
}

func printToFile(lines []prData, reviewersList []string, f io.Writer) {
	if len(lines) == 0 {
		return
	}

	headersStr := ""
	for _, h := range reviewersList {
		headersStr = fmt.Sprintf("%s , %s", headersStr, h)
	}

	_, _ = fmt.Fprintf(f, "number, type, title, state, created, closed, days, comments %s\n", headersStr)

	for _, l := range lines {
		str := fmt.Sprintf("%d, %s, %s, %s, %s, %s, %s, %d",
			l.number, l.prType, l.title, l.state, l.created, l.closed, l.duration, l.comments)
		for _, n := range reviewersList {
			v := l.reviewers[n]
			if v > 0 {
				str = fmt.Sprintf("%s, %d", str, v)
			} else {
				str = fmt.Sprintf("%s, ", str)
			}
		}
		_, _ = fmt.Fprintf(f, "%s\n", str)
	}
}

