package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Hugal31/mePicture/config"
	"github.com/Hugal31/mePicture/database"
	"github.com/Hugal31/mePicture/picture"
)

func pictureUsage() {
	fmt.Fprintln(os.Stderr, "Manage pictures\n"+
		"\n"+
		"Usage:\n"+
		"\n"+
		"\tmePicture picture command [arguments]\n"+
		"\n"+
		"The commands are:\n"+
		"\n"+
		"\tadd target tagName...        Tag pictures\n"+
		"\tlist [-p] [-f] [tagName...]  List pictures, filter with the tags given in parameter\n"+
		"\t                             Add -p to prevent the display of tag namse\n"+
		"\t                             Add -f to show full path\n"+
		"\tremove target tagName...     Remove tags from picture\n"+
		"\tdelete target...             Remove all tags from target\n"+
		"\n"+
		"target:  Image file or directory")
	os.Exit(1)
}

// Handler for list command
func CommandPicture(args []string) {
	if len(args) == 0 {
		pictureUsage()
	}

	switch args[0] {
	case "add":
		pictureAddCommand(args[1:])
		break
	case "list":
		pictureListCommand(args[1:])
		break
	case "remove":
		pictureRemoveCommand(args[1:])
		break
	case "delete":
		pictureDeleteCommand(args[1:])
		break
	default:
		pictureUsage()
	}
}

func addFileTags(path string, file os.FileInfo, tagNames []string, db *database.DB) {
	if file.IsDir() {
		subFiles, err := ioutil.ReadDir(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, subfile := range subFiles {
			addFileTags(path+string(os.PathSeparator)+subfile.Name(), subfile, tagNames, db)
		}
	} else if filepath.Ext(path) == ".png" || filepath.Ext(path) == ".jpg" { // TODO Refactor
		rel := getPicturePath(path)
		pic := db.PictureAdd(rel)
		db.PictureAddTags(&pic, tagNames)
	}
}

func PictureAddTags(path string, tagNames []string) {
	checkTagNames(tagNames)

	file, err := os.Stat(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	db := database.Open()
	defer db.Close()

	db.Begin()
	db.AddTags(tagNames)
	addFileTags(path, file, tagNames, db)
	db.Commit()
}

func pictureAddCommand(args []string) {
	if len(args) < 2 {
		pictureUsage()
	}
	PictureAddTags(args[0], args[1:])
}

func PictureList(displayTags, fullPath bool, tags []string) {
	db := database.Open()
	defer db.Close()

	var pictures picture.PictureSlice
	if len(tags) != 0 {
		pictures = db.ListPictureWithTags(tags)
	} else {
		pictures = db.ListPicture()
	}
	for _, pic := range pictures {
		if fullPath {
			fmt.Printf("%s%c%s", config.GetConfig().PicturesRoot, os.PathSeparator, pic.Name)
		} else {
			fmt.Print(pic.Name)
		}
		if displayTags {
			fmt.Print(", \t\t")
			for i, tag := range pic.Tags {
				if i == 0 {
					fmt.Print(tag.Name)
				} else {
					fmt.Printf(", %s", tag.Name)
				}
			}
		}
		fmt.Println()
	}
}

func pictureListCommand(args []string) {
	// Check if the '-p' or '-f' flag is present
	// TODO Use a real argument parser
	displayTags := true
	fullPath := false
	for i := 0; i < len(args); {
		if args[i] == "-p" {
			displayTags = false
			args = append(args[:i], args[i+1:]...)
		} else if args[i] == "-f" {
			fullPath = true
			args = append(args[:i], args[i+1:]...)
		} else if len(args[i]) > 0 && args[i][0] == '-' {

		} else {
			i++
		}
	}
	PictureList(displayTags, fullPath, args)
}

func pictureRemoveCommand(args []string) {
	if len(args) < 2 {
		pictureUsage()
	}

	db := database.Open()
	defer db.Close()

	path := getPicturePath(args[0])
	pic := db.PictureFromPath(path)

	for _, tagName := range args[1:] {
		t := db.TagFromName(tagName)
		db.PictureRemoveTag(pic, t)
	}
}

func pictureDeleteCommand(args []string) {
	if len(args) < 1 {
		pictureUsage()
	}

	db := database.Open()
	defer db.Close()

	path := getPicturePath(args[0])
	pic := db.PictureFromPath(path)
	db.PictureDelete(pic)
}

func getPicturePath(path string) string {
	path, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	rel, err := filepath.Rel(config.GetConfig().PicturesRoot, path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s is not in the picture root\n", path)
		os.Exit(1)
	}
	return rel
}
