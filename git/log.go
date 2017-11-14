package git

import (
	"errors"
	"sort"
	"strings"
)

// Log is the main struct that we use to track Git log for a repository.
type Log struct {
	Hash      string
	Author    Author
	Committer Author
}

type Author struct {
	Name  string
	Email string
}

type Authors []Author

type Logs []Log

func (r *Repo) Logs() (logs Logs, err error) {
	cmd, stdout, stderr := r.Git("--no-pager", "log", "--pretty=%h;%aN,%ae;%cN,%ce")
	if err := cmd.Run(); err != nil {
		return logs, errors.New(stderr.String())
	}

	for _, line := range strings.Split(stdout.String(), "\n") {
		parts := strings.SplitN(line, ";", 3)
		if len(parts) != 3 {
			continue
		}
		authorParts := strings.SplitN(strings.TrimSpace(parts[1]), ",", 2)
		if len(authorParts) != 2 {
			continue
		}
		committerParts := strings.SplitN(strings.TrimSpace(parts[2]), ",", 2)
		if len(committerParts) != 2 {
			continue
		}

		logs = append(logs, Log{
			Hash: strings.TrimSpace(parts[0]),
			Author: Author{
				Name:  authorParts[0],
				Email: authorParts[1],
			},
			Committer: Author{
				Name:  committerParts[0],
				Email: committerParts[1],
			},
		})
	}

	return logs, err
}

func (l Logs) GetAuthors() (authors Authors) {
	for _, log := range l {
		authors = appendUniqueAuthor(authors, log.Author)
		authors = appendUniqueAuthor(authors, log.Committer)
	}
	sort.Sort(authors)
	return authors
}

func appendUniqueAuthor(authors Authors, author Author) Authors {
	for _, anAuthor := range authors {
		if anAuthor.Name == author.Name && anAuthor.Email == author.Email {
			return authors
		}
	}
	return append(authors, author)
}

func (slice Authors) Len() int {
	return len(slice)
}

func (slice Authors) Less(i, j int) bool {
	switch strings.Compare(strings.ToUpper(slice[i].Name), strings.ToUpper(slice[j].Name)) {
	case -1:
		return true
	case 0, 1:
		return false
	default:
		return false
	}
}

func (slice Authors) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
