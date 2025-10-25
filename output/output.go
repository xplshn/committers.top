package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"most-active-github-users-counter/github"
	"most-active-github-users-counter/top"
)

type Format func(results github.GithubSearchResults, writer io.Writer, options top.Options) error

func PlainOutput(results github.GithubSearchResults, writer io.Writer, options top.Options) error {
	users := GithubUserList(results.Users)
	fmt.Fprintln(writer, "USERS\n--------")
	for i, user := range users {
		fmt.Fprintf(writer, "#%+v: %+v (%+v):%+v (%+v) %+v\n", i+1, user.Name, user.Login, user.ContributionCount, user.Company, strings.Join(user.Organizations, ","))
	}
	fmt.Fprintln(writer, "\nORGANIZATIONS\n--------")
	for i, org := range users.TopOrgs(10) {
		fmt.Fprintf(writer, "#%+v: %+v (%+v)\n", i+1, org.Name, org.MemberCount)
	}
	return nil
}

func CsvOutput(results github.GithubSearchResults, writer io.Writer, options top.Options) error {
	users := GithubUserList(results.Users)
	w := csv.NewWriter(writer)
	if err := w.Write([]string{"rank", "name", "login", "contributions", "company", "organizations"}); err != nil {
		return err
	}
	for i, user := range users {
		rank := strconv.Itoa(i + 1)
		name := user.Name
		login := user.Login
		contribs := strconv.Itoa(user.ContributionCount)
		orgs := strings.Join(user.Organizations, ",")
		company := user.Company
		if err := w.Write([]string{rank, name, login, contribs, company, orgs}); err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}

type ContributionsSelector func(github.User) int

type Filter func(github.User) bool

func (users GithubUserList) TopBy(selector func(github.User) int, userFilter Filter, amount int) GithubUserList {
	cloned := clone(users)
	if userFilter != nil {
		var filtered []github.User
		for _, u := range cloned {
			if userFilter(u) {
				filtered = append(filtered, u)
			}
		}
		cloned = filtered
	}
	sort.Slice(cloned, func(i, j int) bool {
		return selector(cloned[i]) > selector(cloned[j])
	})
	return trim(cloned, amount)
}

func YamlOutput(results github.GithubSearchResults, writer io.Writer, options top.Options) error {
	users := GithubUserList(results.Users)

	outputOrganizations := func(orgs Organizations) {
		for i, org := range orgs {
			fmt.Fprintf(
				writer,
				`
  - rank: %+v
    name: %+v
    membercount: %+v
`,
				i+1,
				strconv.QuoteToASCII(org.Name),
				org.MemberCount)
		}
	}

	outputUsers := func(user []github.User, cs ContributionsSelector) {
		for i, u := range user {
			contributionCount := cs(u)
			fmt.Fprintf(
				writer,
				`
  - rank: %+v
    name: %+v
    login: %+v
    avatarUrl: %+v
    contributions: %+v
    company: %+v
    organizations: %+v
`,
				i+1,
				strconv.QuoteToASCII(u.Name),
				strconv.QuoteToASCII(u.Login),
				u.AvatarURL,
				contributionCount,
				strconv.QuoteToASCII(u.Company),
				strconv.QuoteToASCII(strings.Join(u.Organizations, ",")))
		}
	}

	topCommits := users.TopBy(func(u github.User) int { return u.CommitsCount }, nil, options.Amount)
	fmt.Fprintln(writer, "users:")
	outputUsers(topCommits, func(u github.User) int { return u.CommitsCount })

	topPublic := users.TopBy(func(u github.User) int { return u.PublicContributionCount }, nil, options.Amount)
	fmt.Fprintln(writer, "users_public_contributions:")
	outputUsers(topPublic, func(u github.User) int { return u.PublicContributionCount })

	topTotal := users.TopBy(func(u github.User) int { return u.ContributionCount }, nil, options.Amount)
	fmt.Fprintln(writer, "\nprivate_users:")
	outputUsers(topTotal, func(u github.User) int { return u.ContributionCount })

	fmt.Fprintln(writer, "\norganizations:")
	outputOrganizations(topCommits.TopOrgs(10))
	fmt.Fprintln(writer, "\npublic_contributions_organizations:")
	outputOrganizations(topPublic.TopOrgs(10))
	fmt.Fprintln(writer, "\nprivate_organizations:")
	outputOrganizations(topTotal.TopOrgs(10))

	fmt.Fprintf(writer, "generated: %+v\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(writer, "min_followers_required: %+v\n", results.MinimumFollowerCount)
	fmt.Fprintf(writer, "total_user_count: %+v\n", results.TotalUserCount)

	if options.PresetTitle != "" && options.PresetChecksum != "" {
		fmt.Fprintf(writer, "title: %+v\n", options.PresetTitle)
		fmt.Fprintf(writer, "definition_checksum: %+v\n", options.PresetChecksum)
	}

	return nil
}

var companyLogin = regexp.MustCompile(`^\@([a-zA-Z0-9]+)$`)

func trim(users GithubUserList, numTop int) GithubUserList {
	if numTop == 0 {
		numTop = 256
	}
	if len(users) < numTop {
		numTop = len(users)
	}
	return users[:numTop]
}

func clone(users GithubUserList) GithubUserList {
	usersCloned := make(GithubUserList, len(users))
	copy(usersCloned, users)
	return usersCloned
}

type GithubUserList []github.User

func (slice GithubUserList) MinFollowers() int {
	if len(slice) == 0 {
		return 0
	}
	followers := math.MaxInt32
	for _, user := range slice {
		if user.FollowerCount < followers {
			followers = user.FollowerCount
		}
	}
	return followers
}

type Organization struct {
	Name        string
	MemberCount int
}

type Organizations []Organization

func (slice Organizations) Len() int {
	return len(slice)
}

func (slice Organizations) Less(i, j int) bool {
	return slice[i].MemberCount > slice[j].MemberCount
}

func (slice Organizations) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice GithubUserList) TopOrgs(count int) Organizations {
	orgsMap := make(map[string]int)
	for _, user := range slice {
		userOrgs := user.Organizations
		orgMatches := companyLogin.FindStringSubmatch(strings.Trim(user.Company, " "))
		if len(orgMatches) > 0 {
			orgLogin := companyLogin.FindStringSubmatch(strings.Trim(user.Company, " "))[1]
			if len(orgLogin) > 0 && !contains(userOrgs, orgLogin) {
				userOrgs = append(userOrgs, orgLogin)
			}
		}

		for _, o := range userOrgs {
			org := strings.ToLower(o)
			orgsMap[org] = orgsMap[org] + 1
		}
	}

	orgs := Organizations{}

	for k, v := range orgsMap {
		orgs = append(orgs, Organization{Name: k, MemberCount: v})
	}
	sort.Sort(orgs)
	if len(orgs) > count {
		return orgs[:count]
	}
	return orgs
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
