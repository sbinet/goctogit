package main

import (
	"fmt"
	"net/http"
	"path"

	"github.com/sbinet/go-commander"
	"github.com/sbinet/go-flag"
	"github.com/sbinet/go-github-client/client"
)

func git_make_cmd_dl_ls() *commander.Command {
	cmd := &commander.Command{
		Run:       git_run_cmd_dl_ls,
		UsageLine: "dl-ls [options] repo",
		Short:     "lists the available downloads of github repository",
		Long: `
dl-ls lists the available downloads of a github repository.

ex:
 $ goctogit dl-ls mana-core
 $ goctogit dl-ls -org my-organization mana-core
`,
		Flag: *flag.NewFlagSet("git-dl-ls", flag.ExitOnError),
	}
	cmd.Flag.String("u", "", "github user account")
	cmd.Flag.String("org", "", "github organization account")

	return cmd
}

func git_run_cmd_dl_ls(cmd *commander.Command, args []string) {
	n := "github-" + cmd.Name()
	if len(args) != 1 {
		err := fmt.Errorf("%s: needs a github repository name", n)
		handle_err(err)
	}

	repo_name := args[0]
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
	// GET /repos/:owner/:repo/downloads
	if org != "" {
		account = org
	}
	url := path.Join("repos", account, repo_name, "downloads")

	fmt.Printf("%s: listing downloads for repository [%s] with account [%s]...\n",
		n, repo_name, account)

	req, err := ghc.NewAPIRequest("GET", url, nil)
	handle_err(err)

	resp, err := ghc.RunRequest(req, new(http.Client))
	handle_err(err)

	if !resp.IsSuccess() {
		err = fmt.Errorf("%s: request did not succeed. got (status=%d) %v\n", n, resp.RawHttpResponse.StatusCode, resp.RawHttpResponse)
		handle_err(err)
	}

	json_arr, err := resp.JsonArray()
	if err != nil {
		handle_err(err)
	}
	for _, elmt := range json_arr {
		json := client.JsonMap(elmt.(map[string]interface{}))
		fmt.Printf("=== %s\n",
			json.GetString("name"),
			)
		fmt.Printf("%3s id=%v\n", "", int64(json.GetFloat("id")))
		fmt.Printf("%3s sz=%v bytes\n", "", int64(json.GetFloat("size")))
		descr := json.GetString("description")
		if descr != "" {
			fmt.Printf("%3s descr=%q\n", "", descr)
		}
		fmt.Printf("%3s %s\n", "", json.GetString("html_url"))
	}

	fmt.Printf("%s: listing downloads for repository [%s] with account [%s]... [done]\n",
		n, repo_name, account)
}

// EOF
