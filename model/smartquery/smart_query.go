package smartquery

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/sql"
)

var EXTRA_DEBUG = false

type SmartQuery struct {
	Name, Comment, Query, OrderBy string
}

func (r *SmartQuery) ValidateQuery() error {
	squirrel := Squirrelizer{}
	_, err := squirrel.BuildSelect("playlist1", "user666", r.OrderBy)
	if err != nil {
		return err
	}
	fakeSelect := squirrel.Sql
	fakeStatement := fakeSelect + r.Query
	println("################################################")
	fmt.Printf("SQL=[%v]", fakeStatement)
	println()
	println("################################################")
	asBytes := []byte(fakeStatement)
	parser := sitter.NewParser()
	parser.SetLanguage(sql.GetLanguage())
	tree, _ := parser.ParseCtx(context.Background(), nil, asBytes)

	err = r.walkTree(tree.RootNode(), 0, asBytes)
	if err != nil {
		return err
	}
	return nil
}

// walk tree looking for syntax errors and/or blacklisted tokens
func (r *SmartQuery) walkTree(node *sitter.Node, depth int, statement []byte) error {
	if EXTRA_DEBUG {
		fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		fmt.Printf("[DEPTH%d] ID=%d TYPE=%v", depth, node.ID(), node.Type())
		fmt.Println()
		fmt.Println(node.Content(statement))
		fmt.Println("+++++")
		fmt.Println(node.String())
		fmt.Println("+++++")
		fmt.Println(node.Symbol())
		fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	}
	if node.IsError() {
		return fmt.Errorf("sql parse failed %v", node.Content(statement))
	}
	if depth == 3 && node.Type() == "order_by" {
		r.OrderBy = r.extractOrderByTarget(node.Content(statement))
	}

	count := int(node.ChildCount())

	for i := 0; i < count; i++ {
		n := node.Child(i)
		err := r.walkTree(n, 1+depth, statement)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SmartQuery) extractOrderByTarget(orderByStatement string) string {
	// println("+++++++++++")
	// fmt.Printf("[1] ][%v]", orderByStatement)
	// println()
	// println("+++++++++++")

	result := strings.TrimSpace(orderByStatement[8:])
	// println("+++++++++++")
	// fmt.Printf("[2] ][%v]", result)
	// println()
	// println("+++++++++++")

	return result
}

func (r *SmartQuery) MarshalJSON() ([]byte, error) {
	tmp := struct {
		Name    string `json:"name"`
		Comment string `json:"comment"`
		Query   string `json:"query"`
		OrderBy string `json:"orderby"`
	}{
		Name:    r.Name,
		Comment: r.Comment,
		Query:   r.Query,
		OrderBy: r.OrderBy,
	}

	representation, err := json.Marshal(tmp)
	if err != nil {
		return nil, err
	}
	return representation, nil
}

func (r *SmartQuery) UnmarshalJSON(data []byte) (SmartQuery, error) {
	var aux struct {
		Name    string `json:"name"`
		Comment string `json:"comment"`
		Query   string `json:"query"`
		OrderBy string `json:"orderby"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return SmartQuery{}, err
	}
	return SmartQuery{
		Name:    aux.Name,
		Comment: aux.Comment,
		Query:   aux.Query,
		OrderBy: aux.OrderBy}, nil
}
