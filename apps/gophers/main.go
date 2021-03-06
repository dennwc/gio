// SPDX-License-Identifier: Unlicense OR MIT

package main

// A Gio program that displays Go contributors from GitHub. See https://gioui.org for more information.

import (
	"context"
	"flag"
	"fmt"
	"image"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"

	_ "image/jpeg"
	_ "image/png"

	_ "net/http/pprof"

	"gioui.org/ui"
	"gioui.org/ui/app"
	"gioui.org/ui/gesture"
	"gioui.org/ui/key"
	"gioui.org/ui/layout"

	"github.com/google/go-github/v24/github"
)

type App struct {
	w *app.Window

	ui *UI

	updateUsers   chan []*user
	commitsResult chan []*github.Commit
	ctx           context.Context
	ctxCancel     context.CancelFunc
}

var (
	profile = flag.Bool("profile", false, "serve profiling data at http://localhost:6060")
	stats   = flag.Bool("stats", false, "show rendering statistics")
	token   = flag.String("token", "", "Github authentication token")
)

func main() {
	flag.Parse()
	initProfiling()
	if *token == "" {
		fmt.Println("The quota for anonymous GitHub API access is very low. Specify a token with -token to avoid quota errors.")
		fmt.Println("See https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line.")
	}
	go func() {
		w := app.NewWindow(
			app.WithWidth(ui.Dp(400)),
			app.WithHeight(ui.Dp(800)),
			app.WithTitle("Gophers"),
		)
		if err := newApp(w).run(); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

func initProfiling() {
	if !*profile {
		return
	}
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}

func (a *App) run() error {
	a.ui.profiling = *stats
	ops := new(ui.Ops)
	var cfg app.Config
	for {
		select {
		case users := <-a.updateUsers:
			a.ui.users = users
			a.ui.userClicks = make([]gesture.Click, len(users))
			a.w.Invalidate()
		case commits := <-a.commitsResult:
			a.ui.selectedUser.commits = commits
			a.w.Invalidate()
		case e := <-a.w.Events():
			switch e := e.(type) {
			case key.Event:
				switch e.Name {
				case key.NameEscape:
					os.Exit(0)
				case 'P':
					if e.Modifiers.Contain(key.ModCommand) {
						a.ui.profiling = !a.ui.profiling
						a.w.Invalidate()
					}
				}
			case app.DestroyEvent:
				return e.Err
			case app.StageEvent:
				if e.Stage >= app.StageRunning {
					if a.ctxCancel == nil {
						a.ctx, a.ctxCancel = context.WithCancel(context.Background())
					}
					if a.ui.users == nil {
						go a.fetchContributors()
					}
				} else {
					if a.ctxCancel != nil {
						a.ctxCancel()
						a.ctxCancel = nil
					}
				}
			case *app.CommandEvent:
				switch e.Type {
				case app.CommandBack:
					if a.ui.selectedUser != nil {
						a.ui.selectedUser = nil
						e.Cancel = true
						a.w.Invalidate()
					}
				}
			case app.UpdateEvent:
				ops.Reset()
				cfg = e.Config
				cs := layout.RigidConstraints(e.Size)
				a.ui.Layout(&cfg, a.w.Queue(), ops, cs)
				a.w.Update(ops)
			}
		}
	}
}

func newApp(w *app.Window) *App {
	a := &App{
		w:             w,
		updateUsers:   make(chan []*user),
		commitsResult: make(chan []*github.Commit, 1),
	}
	fetch := func(u string) {
		a.fetchCommits(a.ctx, u)
	}
	a.ui = newUI(fetch)
	return a
}

func githubClient(ctx context.Context) *github.Client {
	var tc *http.Client
	if *token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: *token},
		)
		tc = oauth2.NewClient(ctx, ts)
	}
	return github.NewClient(tc)
}

func (a *App) fetchContributors() {
	client := githubClient(a.ctx)
	cons, _, err := client.Repositories.ListContributors(a.ctx, "golang", "go", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "github: failed to fetch contributors: %v\n", err)
		return
	}
	var users []*user
	userErrs := make(chan error, len(cons))
	avatarErrs := make(chan error, len(cons))
	for _, con := range cons {
		con := con
		avatar := con.GetAvatarURL()
		if avatar == "" {
			continue
		}
		u := &user{
			login: con.GetLogin(),
		}
		users = append(users, u)
		go func() {
			guser, _, err := client.Users.Get(a.ctx, u.login)
			if err != nil {
				avatarErrs <- err
				return
			}
			u.name = guser.GetName()
			u.company = guser.GetCompany()
			avatarErrs <- nil
		}()
		go func() {
			a, err := fetchImage(avatar)
			u.avatar = a
			userErrs <- err
		}()
	}
	for i := 0; i < len(cons); i++ {
		if err := <-userErrs; err != nil {
			fmt.Fprintf(os.Stderr, "github: failed to fetch user: %v\n", err)
		}
		if err := <-avatarErrs; err != nil {
			fmt.Fprintf(os.Stderr, "github: failed to fetch avatar: %v\n", err)
		}
	}
	// Drop users with no avatar or name.
	for i := len(users) - 1; i >= 0; i-- {
		if u := users[i]; u.name == "" || u.avatar == nil {
			users = append(users[:i], users[i+1:]...)
		}
	}
	a.updateUsers <- users
}

func fetchImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetchImage: http.Get(%q): %v", url, err)
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fetchImage: image decode failed: %v", err)
	}
	return img, nil
}

func (a *App) fetchCommits(ctx context.Context, user string) {
	go func() {
		gh := githubClient(ctx)
		repoCommits, _, err := gh.Repositories.ListCommits(ctx, "golang", "go", &github.CommitsListOptions{
			Author: user,
		})
		if err != nil {
			log.Printf("failed to fetch commits: %v", err)
			return
		}
		var commits []*github.Commit
		for _, commit := range repoCommits {
			if c := commit.GetCommit(); c != nil {
				commits = append(commits, c)
			}
		}
		a.commitsResult <- commits
	}()
}
