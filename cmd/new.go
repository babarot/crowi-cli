package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/b4b4r07/crowi/api"
	"github.com/b4b4r07/crowi/cli"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [FILE/DIR]",
	Short: "Create a new page",
	Long:  `Create a new page. if no args are given, write a page with editor`,
	RunE:  new,
}

var (
	imageURL = regexp.MustCompile(`!\[.*?\]\((.+?)\)`)
	urlSafe  = strings.NewReplacer(
		`^`, `-`, // for Crowi's regexp
		`$`, `-`,
		`*`, ``,
		`%`, ``, // query
		`?`, ``,
		`.`, `_`,
	)
	datePath = regexp.MustCompile(`^/.*/(.*?/\d{4}/\d{2}/\d{2})$`)
)

func new(cmd *cobra.Command, args []string) error {
	var (
		pages []page
		err   error
	)
	switch {
	case len(args) == 0:
		pages, err = makeFromEditor()
	case len(args) > 0:
		pages, err = makeFromArgs(args)
	}
	if err != nil {
		return err
	}

	if len(pages) == 0 {
		return errors.New("no pages")
	}

	client, err := cli.NewClient()
	if err != nil {
		return err
	}

	apipage := api.NewPage(client)

	// create pages
	for _, page := range pages {
		res, err := apipage.Create(page.path, page.body)
		if err != nil {
			log.Printf("[ERROR] %v", err.Error())
			continue
		}
		if !res.OK {
			log.Printf("[ERROR] %v", res.Error)
			continue
		}
		cli.Underline("Created", res.Page.ID)
		// Attachments
		if imageURL.MatchString(page.body) {
			var (
				find = imageURL.FindAllStringSubmatch(page.body, -1)
				file = find[0][1]
				id   = res.Page.ID
				body = page.body
			)
			if _, err := os.Stat(file); err == nil {
				apipage.Attach(id, file)
				images, err := apipage.Images(id)
				if err != nil {
					continue
				}
				// get attachments URLs and replace body with these
				for _, image := range images.Attachments {
					if image.OriginalName == filepath.Base(file) {
						body = imageURL.ReplaceAllString(body, fmt.Sprintf("![](%s)", image.URL))
					}
				}
				// update if changed
				if body != page.body {
					apipage.Update(id, body)
				}
			}
		}
	}

	return nil
}

// Constituent elements of the page that is not yet made
// (the page to be made from now)
type page struct {
	path, body string
}

func makeFromEditor() (pages []page, err error) {
	user := cli.Conf.Crowi.User
	if user == "" {
		return pages, errors.New("config user not defined")
	}
	date := time.Now().Format("2006/01/02")
	defaultPath := path.Join("/user", user, cli.Conf.Crowi.PageName, date)
	cli.ScanDefaultString = defaultPath + "/"

	pagepath, err := cli.Scan(color.YellowString("Path> "))
	if err != nil {
		return
	}
	if !filepath.HasPrefix(pagepath, "/") {
		return pages, errors.New("path: it must start with a slash")
	}
	// Do not make it a portal page
	pagepath = strings.TrimSuffix(pagepath, "/")

	f, err := cli.TempFile(filepath.Base(pagepath) + cli.Extention)
	defer os.Remove(f.Name())

	var content []byte
	matched := datePath.FindStringSubmatch(pagepath)
	if len(matched) == 0 {
		content = []byte(fmt.Sprintf("# %s", path.Base(pagepath)))
	} else {
		// matched
		content = []byte(fmt.Sprintf("# %s", matched[1]))
	}

	// write content and ignore error if occured
	f.Write(content)
	f.Sync()

	editor := cli.Conf.Core.Editor
	if editor == "" {
		return pages, errors.New("config editor not defined")
	}
	err = cli.Run(editor, f.Name())
	if err != nil {
		return
	}

	body := cli.FileContent(f.Name())
	if body == "" || body == string(content) {
		return pages, errors.New("did nothing due to no contents")
	}

	return []page{{
		path: urlSafe.Replace(pagepath),
		body: body,
	}}, nil
}

func makeFromArgs(args []string) (pages []page, err error) {
	var (
		mdfiles []string
	)

	isdir := func(path string) bool {
		if stat, err := os.Stat(path); err == nil && stat.IsDir() {
			return true
		}
		return false
	}

	for _, arg := range args {
		// if the arg is dir, walk within the dir and add them to slice
		// otherwise (regular file), just add it to slice
		if isdir(arg) {
			err = filepath.Walk(arg, func(arg string, info os.FileInfo, err error) error {
				// skip like .git
				if strings.HasPrefix(arg, ".") {
					return nil
				}
				if info.IsDir() {
					return nil
				}
				switch filepath.Ext(arg) {
				case ".md", ".mkd", ".markdown":
					mdfiles = append(mdfiles, arg)
				}
				return nil
			})
			if err != nil {
				return
			}
		} else {
			mdfiles = append(mdfiles, arg)
		}
	}

	if len(mdfiles) == 0 {
		return pages, errors.New("no markdown files")
	}

	for _, file := range mdfiles {
		file, _ = filepath.Abs(file)
		pagepath := strings.TrimRight(file[len(os.Getenv("HOME")):], filepath.Ext(file))
		cli.ScanDefaultString = filepath.Join("/user", cli.Conf.Crowi.User, pagepath)
		pagepath, err = cli.Scan(color.YellowString("Path> "))
		if err != nil {
			return pages, err
		}
		pages = append(pages, page{
			path: urlSafe.Replace(pagepath),
			body: cli.FileContent(file),
		})
	}

	return pages, err
}

func init() {
	RootCmd.AddCommand(newCmd)
}
