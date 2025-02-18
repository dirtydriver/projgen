package templater

import (
	"fmt"
	"strings"
	"sync"
	"text/template"
	"text/template/parse"
)

var uniqCheck map[string]bool
var uniqElements []string

func addElement(e string) {
	if uniqCheck[e] {
		return
	}
	uniqElements = append(uniqElements, e)
	uniqCheck[e] = true
}

func collectPlaceholders(node parse.Node, placeholders map[string]struct{}) {
	switch n := node.(type) {
	case *parse.ListNode:
		if n == nil {
			return
		}
		for _, child := range n.Nodes {
			collectPlaceholders(child, placeholders)
		}
	case *parse.ActionNode:
		collectPlaceholders(n.Pipe, placeholders)
	case *parse.PipeNode:
		for _, cmd := range n.Cmds {
			collectPlaceholders(cmd, placeholders)
		}
	case *parse.CommandNode:
		for _, arg := range n.Args {
			collectPlaceholders(arg, placeholders)
		}
	case *parse.FieldNode:
		// FieldNode represents expressions like {{ .GroupID }}
		// The Ident slice contains the field names, which we join with dots.
		placeholder := "." + strings.Join(n.Ident, ".")
		placeholders[placeholder] = struct{}{}
	case *parse.ChainNode:
		// ChainNode represents chained accesses like {{ .User.Name }}
		var base string
		if fieldNode, ok := n.Node.(*parse.FieldNode); ok {
			base = "." + strings.Join(fieldNode.Ident, ".")
		} else {
			// Fallback: use the node's string representation.
			base = fmt.Sprintf("%v", n.Node)
		}
		if len(n.Field) > 0 {
			base += "." + strings.Join(n.Field, ".")
		}
		placeholders[base] = struct{}{}
	}
}

func CollectParameters(tempFiles []string) ([]string, error) {

	var (
		wg         sync.WaitGroup
		mu         sync.Mutex
		parameters = make(map[string]struct{})
		errChan    = make(chan error, len(tempFiles))
	)

	for _, file := range tempFiles {
		wg.Add(1)
		go func(file string) {

			defer wg.Done()
			tmpl, err := template.ParseFiles(file)
			if err != nil {
				errChan <- err
				return
			}
			localParameters := make(map[string]struct{})
			collectPlaceholders(tmpl.Root, localParameters)
			mu.Lock()
			for param := range localParameters {
				parameters[param] = struct{}{}
			}
			mu.Unlock()
		}(file)

	}
	wg.Wait()
	close(errChan)
	var placeholderList []string
	for param := range parameters {

		placeholderList = append(placeholderList, param)
	}

	return placeholderList, nil

}
