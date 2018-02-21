package main

import (
	"fmt"
	"io/ioutil"
	"os"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	yaml "gopkg.in/yaml.v2"
)

func fatal(format string, a ...interface{}) {
	s := "Error: " + format + "\n"
	if a != nil {
		fmt.Printf(s, a)
	} else {
		fmt.Print(s)
	}
	os.Exit(1)
}

// Project is a repo that will map to any number of services
type Project struct {
	Repo string `yaml:"repo"`
}

func main() {
	fmt.Println("Let's assemble!")

	projectsFile, err := ioutil.ReadFile("projects.yaml")
	if err != nil {
		fatal("Failed to load repos.yaml file")
	}

	var projects map[string]Project
	err = yaml.Unmarshal(projectsFile, &projects)
	if err != nil {
		fatal("Failed to read repos.yaml file, %v", err)
	}

	pem, _ := ioutil.ReadFile("/Users/danielfarrelly/.ssh/id_rsa")
	signer, err := ssh.ParsePrivateKey(pem)
	if err != nil {
		fatal("ssh key %v", err)
	}
	auth := &ssh2.PublicKeys{User: "git", Signer: signer}

	// auth, err := ssh.NewSSHAgentAuth("git")
	// if err != nil {
	// 	fatal("Failed to set up git auth")
	// }

	for project, config := range projects {
		fmt.Printf("Project: %s in %s\n", project, config.Repo)

		_, err := git.PlainClone(project, false, &git.CloneOptions{
			URL:      config.Repo,
			Progress: os.Stdout,
			Auth:     auth,
		})
		if err != nil {
			fatal("Failed to clone repo for %s, %v", project, err)
		}
	}

	// fmt.Println(projects)
	fmt.Printf("%#v\n", projects)
}
