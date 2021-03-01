package tui

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/HankiGreed/termui/v3/widgets"
	"github.com/gabriel-vasile/mimetype"
)

type nodeValue string

func (n nodeValue) String() string {
	return string(n)
}

// Gets the list of files in the directory
func getTreeNodesFromDirectoryListing(directory string) []*widgets.TreeNode {
	directoryListing, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatalln(directory, err)
	}

	var nodes []*widgets.TreeNode
	for _, file := range directoryListing {
		if file.IsDir() {
			nodes = append(nodes, &widgets.TreeNode{
				Value: nodeValue(file.Name()),
				Nodes: []*widgets.TreeNode{
					{
						Value: nodeValue("Just a placeholder to show that the tree exists"),
						Nodes: nil,
					},
				},
			})
		} else {
			if mime, _ := mimetype.DetectFile(directory + "/" + file.Name()); strings.Contains(mime.String(), "audio") {
				nodes = append(nodes, &widgets.TreeNode{Value: nodeValue(file.Name()), Nodes: nil})
			}
		}
	}
	return nodes
}
