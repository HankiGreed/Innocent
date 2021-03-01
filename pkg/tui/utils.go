package tui

import (
	"io/ioutil"
	"log"

	"github.com/HankiGreed/termui/v3/widgets"
)

type nodeValue string

func (n nodeValue) String() string {
	return string(n)
}

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
			nodes = append(nodes, &widgets.TreeNode{Value: nodeValue(file.Name()), Nodes: nil})
		}
	}
	return nodes
}
