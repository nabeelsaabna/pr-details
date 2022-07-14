package runner

type prData struct {
	number    int
	title     string
	prType    string
	state     string
	created   string
	closed    string
	duration  string
	comments  int
	reviewers map[string]int
}

