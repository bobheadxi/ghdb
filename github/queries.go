package github

type repositoriesQuery struct {
	Viewer struct {
		Repositories struct {
			Nodes    []repoTableInfo
			PageInfo pageInfo
		} `graphql:"repositories(affiliations: $affs, ownerAffiliations: $oAffs, first: 100, after: $cursor)"`
	}
}
