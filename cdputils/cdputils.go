package cdputils

import (
	"fmt"
	"strings"

	"github.com/chromedp/cdproto/cdp"
)

// FirstChildByClass returns the first node child found in the provided node that contains the provided class or nil if no matching node was found. Recursive = false will only search direct children of the provided node
func FirstChildByClass(n *cdp.Node, className string, recursive bool) *cdp.Node {

	for _, c := range n.Children {

		if HasClass(c, className) {
			return c
		}

		if recursive {
			if r := FirstChildByClass(c, className, recursive); r != nil {
				return r
			}
		}
	}

	return nil
}

// NodeText recursively returns the inner text of an HTML node
func NodeInnerText(n *cdp.Node) string {

	buf := &strings.Builder{}

	buf.WriteString(NodeText(n))

	if len(n.Children) > 0 {
		for _, c := range n.Children {
			str := NodeInnerText(c)

			if str != "" {
				buf.WriteString(fmt.Sprintf(" %s", str))
			}
		}
	}

	return strings.TrimSpace(buf.String())
}

// NodeText returns the inner text of an HTML node
func NodeText(n *cdp.Node) string {

	return strings.TrimSpace(n.NodeValue)
}

// HasClass returns true if a node has a particular class, false otherwise
func HasClass(n *cdp.Node, className string) bool {

	classes := n.AttributeValue("class")
	return strings.Contains(classes, className)
}
