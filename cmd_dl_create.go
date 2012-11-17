package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	//net_url "net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/sbinet/go-commander"
	"github.com/sbinet/go-flag"
	"github.com/sbinet/go-github-client/client"
)

type s3_data struct {
	Url            string `json:"url"`
	HtmlUrl        string `json:"html_url"`
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Descr          string `json:"description"`
	Size           int    `json:"size"`
	DlCount        int    `json:"download_count"`
	ContentType    string `json:"content_type"`
	Policy         string `json:"policy"`
	Signature      string `json:"signature"`
	Bucket         string `json:"bucket"`
	AccessKeyId    string `json:"accesskeyid"`
	Path           string `json:"path"`
	Acl            string `json:"acl"`
	ExpirationDate string `json:"expirationdate"`
	Prefix         string `json:"prefix"`
	MimeType       string `json:"mime_type"`
	Redirect       bool   `json:"redirect"`
	S3Url          string `json:"s3_url"`
}

func get_s3_data(r *client.GithubResult) (s3 s3_data, err error) {

	if r.RawHttpResponse.ContentLength == 0 {
		// EMPTY OBJECT
		err = json.Unmarshal(([]byte)("[]"), &s3)
	} else if r.RawHttpResponse.ContentLength != 0 {
		defer r.RawHttpResponse.Body.Close()
		data := make([]byte, 0, r.RawHttpResponse.ContentLength)
		data, err = ioutil.ReadAll(r.RawHttpResponse.Body)

		if err != nil {
			return s3, err
		}

		err = json.Unmarshal(data, &s3)
	}

	return s3, err
}

func git_make_cmd_dl_create() *commander.Command {
	cmd := &commander.Command{
		Run:       git_run_cmd_dl_create,
		UsageLine: "dl-create [options] -f=file -repo=repo",
		Short:     "creates a new download on github",
		Long: `
dl-create creates a new download asset on a github repository.

ex:
 $ goctogit dl-create -descr "new tarball" -f=foo.tar.gz -repo=mana-core
 $ goctogit dl-create -descr "new tarball" -f=foo.tar.gz -repo=mana-core -org my-organyzation
`,
		Flag: *flag.NewFlagSet("git-dl-create", flag.ExitOnError),
	}
	cmd.Flag.String("descr", "", "description of the new github repository")
	cmd.Flag.String("f", "", "path to file to upload on github")
	cmd.Flag.String("u", "", "github user account")
	cmd.Flag.String("repo", "", "name of the github repository")
	cmd.Flag.String("org", "", "github organization account")

	return cmd
}

func git_run_cmd_dl_create(cmd *commander.Command, args []string) {
	n := "github-" + cmd.Name()
	if len(args) != 0 {
		err := fmt.Errorf("%s: does NOT take any positional parameter", n)
		handle_err(err)
	}

	repo_name := cmd.Flag.Lookup("repo").Value.Get().(string)
	fname := cmd.Flag.Lookup("f").Value.Get().(string)

	if repo_name == "" {
		err := fmt.Errorf("%s: needs a github repository name where to put the download", n)
		handle_err(err)
	}
	if fname == "" {
		err := fmt.Errorf("%s: needs a path to a file to upload", n)
		handle_err(err)
	}
	fname, err := filepath.Abs(fname)
	if err != nil {
		handle_err(err)
	}
	if !path_exists(fname) {
		err := fmt.Errorf("%s: needs a path to an EXISTING file to upload", n)
		handle_err(err)
	}

	fi, err := os.Lstat(fname)
	if err != nil {
		handle_err(err)
	}

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
	// POST /repos/:owner/:repo/downloads
	if org != "" {
		account = org
	}
	url := path.Join("repos", account, repo_name, "downloads")

	fmt.Printf("%s: creating download for repository [%s] with account [%s]...\n",
		n, repo_name, account)
	if descr != "" {
		fmt.Printf("%s: descr: %q\n", n, descr)
	}

	data, err := json.Marshal(
		map[string]interface{}{
			"name":        fi.Name(),
			"size":        fi.Size(),
			"description": descr,
			//"content_type": "FIXME: TODO",
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

	s3, err := get_s3_data(resp)
	if err != nil {
		handle_err(err)
	}

	curl_cmd := []string{}
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)
	for _, v := range [][]string{
		{"key", s3.Path},
		{"acl", "public-read"},
		{"success_action_status", "201"},
		{"Filename", s3.Name},
		{"AWSAccessKeyId", s3.AccessKeyId},
		{"Policy", s3.Policy},
		{"Signature", s3.Signature},
		//{"Content-Type",          s3.ContentType},
	} {
		curl_cmd = append(curl_cmd, "-F", `"`+v[0]+`=`+v[1]+`"`)
		fmt.Printf("--write-field: %q %q\n", v[0], v[1])
		err = body_writer.WriteField(v[0], v[1])
		if err != nil {
			handle_err(err)
		}
	}


	{
		//FIXME --- we shouldn't do that...
			curl_cmd = append(curl_cmd, "-F", `"file=@`+fname+`"`, "-v", s3.S3Url)
		fmt.Printf("curl-cmd: %v\n", curl_cmd)
		cmd := exec.Command("curl", curl_cmd...)
		if cmd != nil {
			err = cmd.Run()
			if err != nil {
				handle_err(err)
			}
		}
	}

	if false {
		file_writer, err := body_writer.CreateFormFile("file", fname)
		if err != nil {
			handle_err(err)
		}
		//defer file_writer.Close()

		fh, err := os.Open(fname)
		if err != nil {
			handle_err(err)
		}
		defer fh.Close()

		_, err = io.Copy(file_writer, fh)
		if err != nil {
			handle_err(err)
		}

		content_type := body_writer.FormDataContentType()
		err = body_writer.Close()
		if err != nil {
			handle_err(err)
		}

		fmt.Printf("===> %s\n", content_type)
		/*
		 s3_resp, err := http.Post(s3.S3Url, content_type, body_buf)
		 if err != nil {
		 fmt.Printf("\n%s: s3-request failed:\n%v\n%v\n", n, err, s3_resp)
		 handle_err(err)
		 }

		 if s3_resp.StatusCode != 201 {
		 //err = fmt.Errorf("%s: s3-request did not succeed. got (status=%d) %v\n", n, s3_resp.StatusCode, s3_resp)
		 //handle_err(err)
		 } else {
		 fmt.Printf("%s: response:\n%v\n", n, s3_resp)
		 }
		 */
	}
	fmt.Printf("%s: creating download for repository [%s] with account [%s]... [done]\n",
		n, repo_name, account)
}

// EOF
