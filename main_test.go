package main

func init() {
	repoStatus, err := addHelmRepos(map[string]string{"stable": "https://kubernetes-charts.storage.googleapis.com"})
	if !repoStatus {
		panic(err)
	}
}
