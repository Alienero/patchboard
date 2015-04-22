// Copyright © 2015 Alienero. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package git

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

const realm = "gitter"

func (server *GitServer) createRepo(c web.C, w http.ResponseWriter, r *http.Request) {
	reponame := c.URLParams["reponame"]
	repopath := path.Join(server.workspace, reponame)
	if _, err := os.Stat(repopath); err == nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Repo `%s` already exists!\n", reponame)
	} else {
		gitInitCmd := exec.Command("git", "init", "--bare", repopath)
		_, err := gitInitCmd.CombinedOutput()
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Initialize git repo `%s` failed!\n", reponame)
		} else {
			fmt.Fprintf(w, "Empty git repo `%s` initialized!\n", reponame)
		}
	}
}

func (server *GitServer) deleteRepo(c web.C, w http.ResponseWriter, r *http.Request) {
	reponame := c.URLParams["reponame"]
	repopath := path.Join(server.workspace, reponame)
	if _, err := os.Stat(repopath); os.IsNotExist(err) {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Repo `%s` does not exist!\n", reponame)
	} else {
		err := os.RemoveAll(repopath)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Delete repo `%s` failed!\n", reponame)
		} else {
			fmt.Fprintf(w, "Repo `%s` deleted!\n", reponame)
		}
	}
}

func (server *GitServer) inforefs(c web.C, w http.ResponseWriter, r *http.Request) {
	reponame := c.URLParams["reponame"]
	repopath := path.Join(server.workspace, reponame)
	service := r.FormValue("service")
	if len(service) > 0 {
		w.Header().Add("Content-type", fmt.Sprintf("application/x-%s-advertisement", service))
		gitLocalCmd := exec.Command(
			"git",
			string(service[4:]),
			"--stateless-rpc",
			"--advertise-refs",
			repopath)
		out, err := gitLocalCmd.CombinedOutput()
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintln(w, "Internal Server Error")
			w.Write(out)
		} else {
			serverAdvert := fmt.Sprintf("# service=%s", service)
			length := len(serverAdvert) + 4
			fmt.Fprintf(w, "%04x%s0000", length, serverAdvert)
			w.Write(out)
		}
	} else {
		fmt.Fprintln(w, "Invalid request")
		w.WriteHeader(400)
	}
}

func (server *GitServer) rpc(c web.C, w http.ResponseWriter, r *http.Request) {
	reponame := c.URLParams["reponame"]
	repopath := path.Join(server.workspace, reponame)
	command := c.URLParams["command"]
	if len(command) > 0 {

		w.Header().Add("Content-type", fmt.Sprintf("application/x-git-%s-result", command))
		w.WriteHeader(200)

		gitCmd := exec.Command("git", command, "--stateless-rpc", repopath)

		cmdIn, _ := gitCmd.StdinPipe()
		cmdOut, _ := gitCmd.StdoutPipe()
		body := r.Body

		gitCmd.Start()
		io.Copy(cmdIn, body)
		io.Copy(w, cmdOut)

		if command == "receive-pack" {
			updateCmd := exec.Command("git", "--git-dir", repopath, "update-server-info")
			updateCmd.Start()
		}
	} else {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Invalid Request")
	}
}

func (server *GitServer) generic(c web.C, w http.ResponseWriter, r *http.Request) {
	reponame := c.URLParams["reponame"]
	repopath := path.Join(server.workspace, reponame)
	filepath := path.Join(server.workspace, r.URL.String())
	if strings.HasPrefix(filepath, repopath) {
		http.ServeFile(w, r, filepath)
	} else {
		w.WriteHeader(404)
	}
}

type Loginer func(id, psw string) bool

type GitServer struct {
	workspace string
	login     Loginer
}

func NewGitServer(workspace string, login Loginer) *GitServer {
	return &GitServer{
		workspace: workspace,
		login:     login,
	}
}

func (server *GitServer) Start(addr string) {
	// r := web.New()
	// create and delete repo
	goji.Put("/:reponame", server.createRepo)
	goji.Delete("/:reponame", server.deleteRepo)

	// get repo info/refs
	goji.Get("/:reponame/info/refs", server.inforefs)
	// goji.Head("/:reponame/info/refs", server.inforefs)

	// RPC request on repo
	goji.Post(regexp.MustCompile("^/(?P<reponame>[^/]+)/git-(?P<command>[^/]+)$"), server.rpc)

	// access file contents
	goji.Get("/:reponame/*", server.generic)
	// goji.Head("/:reponame/*", server.generic)

	// start serving
	goji.Use(server.basicAuthMiddleware)
	http.ListenAndServe(addr, goji.DefaultMux)
}

func (server *GitServer) basicAuthMiddleware(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		id, psw, ok := r.BasicAuth()
		if server.login(id, psw) && ok {
			h.ServeHTTP(w, r)
		} else {
			w.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=\"%s\"", realm))
			w.WriteHeader(401)
			fmt.Fprintln(w, "Unauthorized")
		}
	}
	return http.HandlerFunc(fn)
}
