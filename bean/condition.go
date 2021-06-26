package bean

import (
	"regexp"
	"strings"
)

func skipBranch(c *Condition, branch string) bool {
	return !c.Match(branch)
}
func skipCommitMessages(c *Condition, branch string) bool {
	return !c.Match(branch)
}
func skipCommitNotes(c *Condition, branch string) bool {
	return !c.Match(branch)
}

func (c *Condition) Match(v string) bool {
	if c == nil {
		return false
	}
	if c.Include != nil && c.Exclude != nil {
		return c.Includes(v) && !c.Excludes(v)
	}

	if c.Include != nil && c.Includes(v) {
		return true
	}

	if c.Exclude != nil && !c.Excludes(v) {
		return true
	}

	return false
}

func (c *Condition) Excludes(v string) bool {
	for _, in := range c.Exclude {
		if in == "" {
			continue
		}
		if in == v {
			return true
		}
		if isMatch(v, in) {
			return true
		}
		reg, err := regexp.Compile(in)
		if err != nil {
			return false
		}
		match := reg.Match([]byte(strings.Replace(v, "\n", "", -1)))
		if match {
			return true
		}
	}
	return false
}

func (c *Condition) Includes(v string) bool {
	for _, in := range c.Include {
		if in == "" {
			continue
		}
		if in == v {
			return true
		}
		if isMatch(v, in) {
			return true
		}
		reg, err := regexp.Compile(in)
		if err != nil {
			return false
		}
		match := reg.Match([]byte(strings.Replace(v, "\n", "", -1)))
		if match {
			return true
		}
	}
	return false
}

func isMatch(s string, p string) bool {
	m, n := len(s), len(p)
	dp := make([][]bool, m+1)
	for i := 0; i <= m; i++ {
		dp[i] = make([]bool, n+1)
	}
	dp[0][0] = true
	for i := 1; i <= n; i++ {
		if p[i-1] == '*' {
			dp[0][i] = true
		} else {
			break
		}
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if p[j-1] == '*' {
				dp[i][j] = dp[i][j-1] || dp[i-1][j]
			} else if s[i-1] == p[j-1] {
				dp[i][j] = dp[i-1][j-1]
			}
		}
	}
	return dp[m][n]
}
