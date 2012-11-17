package main

import (
	"fmt"
	"net/http"
	"path"

	"github.com/sbinet/go-commander"
	"github.com/sbinet/go-flag"
	"github.com/sbinet/go-github-client/client"
)

func git_make_cmd_dl_rm() *commander.Command {
	cmd := &commander.Command{
		Run:       git_run_cmd_dl_rm,
		UsageLine: "dl-rm [options] -repo=repo file-id",
		Short:     "deletes a download on github by id",
		Long: `
dl-rm deletes a download from a github repository by id.

ex:
 $ goctogit dl-rm -repo=mana-core 1
 $ goctogit dl-rm -repo=mana-core -org my-organization 1
`,
		Flag: *flag.NewFlagSet("git-dl-rm", flag.ExitOnError),
	}
	cmd.Flag.String("u", "", "github user account")
	cmd.Flag.String("repo", "", "name of the github repository")
	cmd.Flag.String("org", "", "github organization account")

	return cmd
}

func git_run_cmd_dl_rm(cmd *commander.Command, args []string) {
	n := "github-" + cmd.Name()
	if len(args) != 1 {
		err := fmt.Errorf("%s: needs a file-id to delete", n)
		handle_err(err)
	}
	file_id := args[0]

	repo_name := cmd.Flag.Lookup("repo").Value.Get().(string)
	if repo_name == "" {
		err := fmt.Errorf("%s: needs a github repository name to delete from", n)
		handle_err(err)
	}

	user := cmd.Flag.Lookup("u").Value.Get().(string)
	org := cmd.Flag.Lookup("org").Value.Get().(string)

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
	// DELETE /repos/:owner/:repo/downloads/:id
	if org != "" {
		account = org
	}
	url := path.Join("repos", account, repo_name, "downloads", file_id)

	fmt.Printf("%s: deleting download id=%s from [%s/%s]...\n",
		n, file_id, account, repo_name)

	req, err := ghc.NewAPIRequest("DELETE", url, nil)
	handle_err(err)

	resp, err := ghc.RunRequest(req, new(http.Client))
	handle_err(err)

	sc := resp.RawHttpResponse.StatusCode
	switch sc {
	case 204:
		// all good
	case 404:
		err = fmt.Errorf("%s: no such file-id\n", n)
	default:
		err = fmt.Errorf("%s: request did not succeed. got (status=%d) %v\n", n, resp.RawHttpResponse.StatusCode, resp.RawHttpResponse)
	}
	handle_err(err)

	fmt.Printf("%s: deleting download id=%s from [%s/%s]... [done]\n",
		n, file_id, account, repo_name)
}

// EOF
