package main

/*
 #include <unistd.h>
 #include <stdlib.h>
 #include <string.h>
*/
import "C"

import (
	"fmt"
	"net/http"
	"unsafe"

	"github.com/sbinet/go-commander"
	"github.com/sbinet/go-flag"
	"github.com/sbinet/go-github-client/client"
)

func git_make_cmd_login() *commander.Command {
	cmd := &commander.Command{
		Run:       git_run_cmd_login,
		UsageLine: "login [options]",
		Short:     "login stores the github authentication data",
		Long: `
login stores the github authentication data

ex:
 $ goctogit login
 $ goctogit login -u mylogin
`,
		Flag: *flag.NewFlagSet("git-create", flag.ExitOnError),
	}
	cmd.Flag.String("u", "", "github user account")

	return cmd
}

func git_run_cmd_login(cmd *commander.Command, args []string) {
	n := "github-" + cmd.Name()

	user := cmd.Flag.Lookup("u").Value.Get().(string)
	fmt.Printf("github username: ")
	if user == "" {
		_, err := fmt.Scanln(&user)
		handle_err(err)
	} else {
		fmt.Printf("%s\n", user)
	}

	password, err := getpasswd("github password: ")
	handle_err(err)

	section := "go-octogit"
	for k, v := range map[string]string{
		"username": user,
		"password": password,
	} {
		if Cfg.HasOption(section, k) {
			Cfg.RemoveOption(section, k)
		}
		if !Cfg.AddOption(section, k, v) {
			err := fmt.Errorf("%s: could not add option [%s] to section [%s]", n, k, section)
			panic(err.Error())
		}
	}

	// check credentials
	ghc, err := client.NewGithubClient(user, password, client.AUTH_USER_PASSWORD)
	handle_err(err)

	req, err := ghc.NewAPIRequest("GET", "authorizations", nil)
	handle_err(err)

	resp, err := ghc.RunRequest(req, new(http.Client))
	handle_err(err)

	if !resp.IsSuccess() {
		err = fmt.Errorf("%s: authentication failed\n%v\n", n, resp.RawHttpResponse)
		handle_err(err)
	}

	err = Cfg.WriteFile(CfgFname, 0600, "")
	handle_err(err)
}

func getpasswd(prompt string) (string, error) {
	c_prompt := C.CString(prompt)
	defer C.free(unsafe.Pointer(c_prompt))

	c_pwd := C.getpass(c_prompt)
	if c_pwd == nil {
		return "", fmt.Errorf("failed to get password")
	}
	passwd := C.GoString(c_pwd)
	C.free(unsafe.Pointer(c_pwd))

	return passwd, nil
}

// EOF
