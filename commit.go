package main

import "math/rand"

type Commit struct {
	RepoName     string        `json:"repo_name"`
	AuthorName   string        `json:"author_name"`
	AuthorEmail  string        `json:"author_email"`
	ChangedFiles []ChangedFile `json:"changed_files"`
}

type ChangedFile struct {
	FileName   string `json:"file_name"`
	Language   string `json:"language"`
	Insertions int    `json:"insertions"`
	Deletions  int    `json:"deletions"`
}

func CreateRepo(c chan<- Commit) {
	for i := 0; i < 100; i++ {
		c <- createCommit()
	}
	close(c)
}

func createCommit() Commit {
	commit := Commit{
		RepoName:     "reponame",
		AuthorName:   "someone",
		AuthorEmail:  "something@email.com",
		ChangedFiles: make([]ChangedFile, 100),
	}

	for i := 0; i < 100; i++ {
		commit.ChangedFiles[i] = ChangedFile{
			FileName:   "filename",
			Language:   "language",
			Insertions: rand.Intn(50),
			Deletions:  rand.Intn(50),
		}
	}

	return commit
}
