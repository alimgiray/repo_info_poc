package main

import "github.com/alimgiray/repo_info_poc/proto"

func ProtoFromCommit(c Commit) *commit.ProtoCommit {
	return &commit.ProtoCommit{
		RepoName:    c.RepoName,
		AuthorName:  c.AuthorName,
		AuthorEmail: c.AuthorEmail,
		ChangedFiles: func() []*commit.ChangedFile {
			r := make([]*commit.ChangedFile, len(c.ChangedFiles))

			for k, v := range c.ChangedFiles {
				r[k] = &commit.ChangedFile{
					FileName:   v.FileName,
					Language:   v.Language,
					Insertions: int64(v.Insertions),
					Deletions:  int64(v.Deletions),
				}
			}

			return r
		}(),
	}
}

func CommitFromProto(p *commit.ProtoCommit) Commit {
	return Commit{
		RepoName:    p.RepoName,
		AuthorName:  p.AuthorName,
		AuthorEmail: p.AuthorEmail,
		ChangedFiles: func() []ChangedFile {
			r := make([]ChangedFile, len(p.ChangedFiles))

			for k, v := range p.ChangedFiles {
				r[k] = ChangedFile{
					FileName:   v.FileName,
					Language:   v.Language,
					Insertions: int(v.Insertions),
					Deletions:  int(v.Deletions),
				}
			}

			return r
		}(),
	}
}
