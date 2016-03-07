package dependencygraph

import "errors"

func (g *Graph) addChildToBack(curNode *Node, childName string) {
	childNode, _ := g.GetOrCreateNode(childName)
	curNode.edge.PushBack(childNode)
}

func (g *Graph) GetOrCreateNode(nodeName string) (ret *Node, not_created bool) {
	if g.nodeList == nil {
		g.nodeList = make(map[string]*Node)
	}
	if ret, not_created = g.nodeList[nodeName]; !not_created {
		ret = &Node{nodeName, NewList()}
		g.nodeList[nodeName] = ret
	}
	return
}

//AddChildrens Add the childs in order at the back of the node list (to keep the order!)
func (g *Graph) AddChildrens(nodeName string, childrenNames ...string) {
	curNode, _ := g.GetOrCreateNode(nodeName)
	for _, childName := range childrenNames {
		g.addChildToBack(curNode, childName)
	}
}

func (g *Graph) walk(curNode, parentNode *Node, f func(string, string, *Graph) error, resolved, seen *List) error {
	parentPath := ""
	if parentNode != nil {
		parentPath = parentNode.path
	}
	if err := f(curNode.path, parentPath, g); err != nil {
		return err
	}
	seen.PushFront(curNode)
	for e := curNode.edge.Front(); e != nil; e = e.Next() {
		val := e.Value
		if resolved.Find(val) == nil {
			if seen.Find(val) != nil {
				return errors.New("CIRCULAR dependencies found with:\nParent: " + val.path + seen.String())
			}
			if err := g.walk(val, curNode, f, resolved, seen); err != nil {
				return err
			}
		}
	}
	resolved.PushBack(curNode)
	return nil
}

func (g *Graph) Walk(entryPoint string, f func(string, string, *Graph) error) ([]string, error) {
	entryNode, _ := g.GetOrCreateNode(entryPoint)
	resolved := NewList()
	seen := NewList()
	if err := g.walk(entryNode, nil, f, resolved, seen); err != nil {
		return nil, err
	}
	ret := make([]string, resolved.Len())
	i := 0
	for e := resolved.Front(); e != nil; e = e.Next() {
		ret[i] = e.Value.path
		i += 1
	}
	return ret, nil
}
