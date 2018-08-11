package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// release type representing Helm releases which are described in the desired state
type release struct {
	Name            string
	Description     string
	Namespace       string
	Enabled         bool
	Chart           string
	Version         string
	ValuesFile      string   `yaml:"valuesFile"`
	ValuesFiles     []string `yaml:"valuesFiles"`
	Purge           bool
	Test            bool
	Protected       bool
	Wait            bool
	Priority        int
	TillerNamespace string
	Set             map[string]string
}

// validateRelease validates if a release inside a desired state meets the specifications or not.
// check the full specification @ https://github.com/gofunky/helmsman/docs/desired_state_spec.md
func validateRelease(appLabel string, r *release, names map[string]map[string]bool, s state) (bool, string) {
	_, err := os.Stat(r.ValuesFile)
	if r.Name == "" {
		r.Name = appLabel
	} else if r.TillerNamespace != "" && r.TillerNamespace != "kube-system" {
		if v, ok := s.Namespaces[r.TillerNamespace]; !ok {
			return false, "tillerNamespace specified, but the namespace specified does not exist!"
		} else if !v.InstallTiller {
			return false, "tillerNamespace specified, but that namespace does not have installTiller set to true."
		}
	} else if names[r.Name][getDesiredTillerNamespace(r)] {
		return false, "release name must be unique within a given Tiller."
	} else if nsOverride == "" && r.Namespace == "" {
		return false, "release targeted namespace can't be empty."
	} else if nsOverride == "" && r.Namespace != "" && !checkNamespaceDefined(r.Namespace, s) {
		return false, "release " + r.Name + " is using namespace [ " + r.Namespace + " ] which is not defined in the Namespaces section of your desired state file." +
			" Release [ " + r.Name + " ] can't be installed in that Namespace until its defined."
	} else if r.Chart == "" || !strings.ContainsAny(r.Chart, "/") {
		return false, "chart can't be empty and must be of the format: repo/chart."
	} else if r.Version == "" {
		return false, "version can't be empty."
	} else if r.ValuesFile != "" && (!isOfType(r.ValuesFile, ".yaml") || err != nil) {
		return false, "valuesFile must be a valid file path for a yaml file, Or can be left empty."
	} else if r.ValuesFile != "" && len(r.ValuesFiles) > 0 {
		return false, "valuesFile and valuesFiles should not be used together."
	} else if len(r.ValuesFiles) > 0 {
		for _, filePath := range r.ValuesFiles {
			if _, pathErr := os.Stat(filePath); !isOfType(filePath, ".yaml") || pathErr != nil {
				return false, "the value for valueFile '" + filePath + "' must be a valid file path for a yaml file."
			}
		}
	} else if r.Priority != 0 && r.Priority > 0 {
		return false, "priority can only be 0 or negative value, positive values are not allowed."
	}

	if names[r.Name] == nil {
		names[r.Name] = make(map[string]bool)
	}
	if r.TillerNamespace != "" {
		names[r.Name][r.TillerNamespace] = true
	} else if s.Namespaces[r.Namespace].InstallTiller {
		names[r.Name][r.Namespace] = true
	} else {
		names[r.Name]["kube-system"] = true
	}

	return true, ""
}

// checkNamespaceDefined checks if a given namespace is defined in the namespaces section of the desired state file
func checkNamespaceDefined(ns string, s state) bool {
	_, ok := s.Namespaces[ns]
	if !ok {
		return false
	}
	return true
}

// overrideNamespace overrides a release defined namespace with a new given one
func overrideNamespace(r *release, newNs string) {
	log.Println("INFO: overriding namespace for app:  " + r.Name)
	r.Namespace = newNs
}

// print prints the details of the release
func (r release) print() {
	fmt.Println("")
	fmt.Println("\tname : ", r.Name)
	fmt.Println("\tdescription : ", r.Description)
	fmt.Println("\tnamespace : ", r.Namespace)
	fmt.Println("\tenabled : ", r.Enabled)
	fmt.Println("\tchart : ", r.Chart)
	fmt.Println("\tversion : ", r.Version)
	fmt.Println("\tvaluesFile : ", r.ValuesFile)
	fmt.Println("\tvaluesFiles : ", strings.Join(r.ValuesFiles, ","))
	fmt.Println("\tpurge : ", r.Purge)
	fmt.Println("\ttest : ", r.Test)
	fmt.Println("\tprotected : ", r.Protected)
	fmt.Println("\twait : ", r.Wait)
	fmt.Println("\tpriority : ", r.Priority)
	fmt.Println("\ttillerNamespace : ", r.TillerNamespace)
	fmt.Println("\tvalues to override from env:")
	printMap(r.Set)
	fmt.Println("------------------- ")
}
