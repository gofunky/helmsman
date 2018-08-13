package main

import (
	"log"
	"os"
)

var s state
var debug bool
var file string
var apply bool
var help bool
var v bool
var verbose bool
var nsOverride string
var checkCleanup bool
var skipValidation bool
var applyLabels bool
var version = "v1.4.0-rc"

func main() {

	// set the kubecontext to be used Or create it if it does not exist
	if !setKubeContext(s.Settings["kubeContext"]) {
		if r, msg := createContext(); !r {
			logError(msg)
		}
		checkCleanup = true
	}

	// add/validate namespaces
	addNamespaces(s.Namespaces)

	if r, msg := initHelm(); !r {
		logError(msg)
	}

	// check if helm Tiller is ready
	for k, v := range s.Namespaces {
		if v.InstallTiller {
			waitForTiller(k)
		}
	}

	if verbose {
		logVersions()
	}

	// add repos -- fails if they are not valid
	if r, msg := addHelmRepos(s.HelmRepos); !r {
		logError(msg)
	}

	if !skipValidation {
		// validate charts-versions exist in defined repos
		if r, msg := validateReleaseCharts(s.Apps); !r {
			logError(msg)
		}
	} else {
		log.Println("INFO: charts validation is skipped.")
	}

	log.Println("INFO: checking what I need to do for your charts ... ")

	p := makePlan(&s)
	cleanUntrackedReleases()

	p.sortPlan()
	p.printPlan()
	p.sendPlanToSlack()

	if apply {
		p.execPlan()
	}

	if checkCleanup {
		cleanup()
	}

	log.Println("INFO: completed successfully!")
}

// cleanup deletes the k8s certificates and keys files
// It also deletes any Tiller TLS certs and keys
func cleanup() {
	if _, err := os.Stat("ca.crt"); err == nil {
		deleteFile("ca.crt")
	}

	if _, err := os.Stat("ca.key"); err == nil {
		deleteFile("ca.key")
	}

	if _, err := os.Stat("client.crt"); err == nil {
		deleteFile("client.crt")
	}

	for k := range s.Namespaces {
		if _, err := os.Stat(k + "-tiller.cert"); err == nil {
			deleteFile(k + "-tiller.cert")
		}
		if _, err := os.Stat(k + "-tiller.key"); err == nil {
			deleteFile(k + "-tiller.key")
		}
		if _, err := os.Stat(k + "-ca.cert"); err == nil {
			deleteFile(k + "-ca.cert")
		}
		if _, err := os.Stat(k + "-client.cert"); err == nil {
			deleteFile(k + "-client.cert")
		}
		if _, err := os.Stat(k + "-client.key"); err == nil {
			deleteFile(k + "-client.key")
		}
	}
}
