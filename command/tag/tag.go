package tag

import (
	"fmt"
	"os"

	"github.com/Hugal31/mePicture/database"
)

func usage() {
	fmt.Fprintln(os.Stderr, "Manage tags\n"+
		"\n"+
		"Usage:\n"+
		"\n"+
		"\tmePicture tag command [arguments]\n"+
		"\n"+
		"The commands are:\n"+
		"\n"+
		"\tadd tagName...      Add a tag\n"+
		"\tlist                List tags")
	os.Exit(1)
}

func ListTags() {
	db := database.Open()
	defer db.Close()

	tags := db.ListTags()
	for _, tag := range tags {
		println(tag.Name)
	}
}

func listTagsCommand([]string) {
	ListTags()
}

func AddTag(tagName string) {
	db := database.Open()
	defer db.Close()
	db.AddTag(tagName)
}

func AddTags(tagNames []string) {
	db := database.Open()
	defer db.Close()
	db.AddTags(tagNames)
}

func addTagCommand(args []string) {
	if len(args) == 0 {
		usage() // TODO Usage spécifique de la commande add ?
	}
	AddTags(args)
}

func CommandTag(args []string) {
	if len(args) == 0 {
		usage()
	}
	switch args[0] {
	case "list":
		listTagsCommand(args[1:])
		break
	case "add":
		addTagCommand(args[1:])
		break
	default:
		usage()
	}
}