package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/gonuts/commander"
	"github.com/gonuts/flag"
	"github.com/sbinet/go-github-client/client"
)

func git_make_cmd_create() *commander.Command {
	cmd := &commander.Command{
		Run:       git_run_cmd_create,
		UsageLine: "create [options] <repo>",
		Short:     "create a new repository on github",
		Long: `
create creates a new git repository on github.

ex:
 $ goctogit create -descr "mana-core is a fine repo" mana-core
 $ goctogit create -descr "a helloworld repo" -u mylogin hello
 $ goctogit create -descr "a hellowrold repo" -org myorganization hello
`,
		Flag: *flag.NewFlagSet("git-create", flag.ExitOnError),
	}
	cmd.Flag.String("u", "", "github user account")
	cmd.Flag.String("org", "", "github organization account")
	cmd.Flag.String("descr", "", "description of the new github repository")

	return cmd
}

func git_run_cmd_create(cmd *commander.Command, args []string) {
	n := "github-" + cmd.Name()
	if len(args) <= 0 {
		err := fmt.Errorf("%s: you need to give a repository name", n)
		handle_err(err)
	}

	repo_name := args[0]

	user := cmd.Flag.Lookup("u").Value.Get().(string)
	org := cmd.Flag.Lookup("org").Value.Get().(string)
	descr := cmd.Flag.Lookup("descr").Value.Get().(string)

	if user == "" {
		v, err := Cfg.String("go-octogit", "username")
		handle_err(err)
		user = v
	}

	password, err := Cfg.String("go-octogit", "password")
	handle_err(err)

	ghc, err := client.NewGithubClient(user, password, client.AUTH_USER_PASSWORD)
	handle_err(err)

	account := user
	url := path.Join("user", "repos")
	if org != "" {
		account = org
		url = path.Join("orgs", org, "repos")
	}

	fmt.Printf("%s: creating repository [%s] with account [%s]...\n",
		n, repo_name, account)
	if descr != "" {
		fmt.Printf("%s: descr: %q\n", n, descr)
	}

	data, err := json.Marshal(
		map[string]interface{}{
			"name":          repo_name,
			"description":   descr,
			"homepage":      "",
			"private":       false,
			"has_issues":    true,
			"has_wiki":      true,
			"has_downloads": true,
		})
	handle_err(err)

	req, err := ghc.NewAPIRequest("POST", url, bytes.NewBuffer(data))
	handle_err(err)

	resp, err := ghc.RunRequest(req, new(http.Client))
	handle_err(err)

	sc := resp.RawHttpResponse.StatusCode
	if !resp.IsSuccess() && sc != 201 {
		err = fmt.Errorf("%s: request did not succeed. got (status=%d) %v\n", n, resp.RawHttpResponse.StatusCode, resp.RawHttpResponse)
		handle_err(err)
	}
	fmt.Printf("%s: you can push from an existing repository like so:\n", n)
	fmt.Printf(`
git remote add origin https://github.com/%s/%s.git
git push -u origin master

git remote add origin git@github.com:%s/%s
git push -u origin master
`, account, repo_name, account, repo_name)
	fmt.Printf("%s: creating repository [%s] with account [%s]... [done]\n",
		n, repo_name, account)
}

// EOF
