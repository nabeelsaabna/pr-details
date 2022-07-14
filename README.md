# pr-details

Simple minimal program to call GitHub and retrieve information about PR history for a single author in a repository

The environment variable `PR_GITHUB_TOKEN` can be provided to overcome the Rate-Limit by GitHub

### Run:
```
go run main.go <owner> <repository> <author>
```

The result will be a `prs_review_<owner>_<repository>_<author>.csv` file

### Example:

```
$ go run main.go nabeelys Excel-To-PDF nabeelys
environment variable PR_GITHUB_TOKEN is not set
Starting processing: nabeelys/Excel-To-PDF for nabeelys
Total Count: 2
Incompelte: false
Items: 2
```

```
$ cat prs_review_nabeelys_Excel-To-PDF_nabeelys.csv
number, type, title, state, created, closed, days, comments
12,  - , Update dependency org.apache.poi:poi-ooxml to v5 remove from local, closed, 2022-06-19, 2022-06-19, 0, 0
14,  - , Update README.md, closed, 2022-07-14, 2022-07-14, 0, 0
```