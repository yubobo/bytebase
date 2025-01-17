package mysql

// Framework code is generated by the generator.

import (
	"fmt"

	"github.com/pingcap/tidb/parser/ast"

	"github.com/bytebase/bytebase/plugin/advisor"
	"github.com/bytebase/bytebase/plugin/advisor/db"
)

var (
	_ advisor.Advisor = (*TableDisallowCreateTableAsAdvisor)(nil)
	_ ast.Visitor     = (*tableDisallowCreateTableAsChecker)(nil)
)

func init() {
	advisor.Register(db.MySQL, advisor.MySQLTableDisallowCreateTableAs, &TableDisallowCreateTableAsAdvisor{})
	advisor.Register(db.TiDB, advisor.MySQLTableDisallowCreateTableAs, &TableDisallowCreateTableAsAdvisor{})
}

// TableDisallowCreateTableAsAdvisor is the advisor checking for disallow CREATE TABLE ... AS ... statement.
type TableDisallowCreateTableAsAdvisor struct {
}

// Check checks for disallow CREATE TABLE ... AS ... statement.
func (*TableDisallowCreateTableAsAdvisor) Check(ctx advisor.Context, statement string) ([]advisor.Advice, error) {
	stmtList, errAdvice := parseStatement(statement, ctx.Charset, ctx.Collation)
	if errAdvice != nil {
		return errAdvice, nil
	}

	level, err := advisor.NewStatusBySQLReviewRuleLevel(ctx.Rule.Level)
	if err != nil {
		return nil, err
	}
	checker := &tableDisallowCreateTableAsChecker{
		level: level,
		title: string(ctx.Rule.Type),
	}

	for _, stmt := range stmtList {
		checker.text = stmt.Text()
		checker.line = stmt.OriginTextPosition()
		(stmt).Accept(checker)
	}

	if len(checker.adviceList) == 0 {
		checker.adviceList = append(checker.adviceList, advisor.Advice{
			Status:  advisor.Success,
			Code:    advisor.Ok,
			Title:   "OK",
			Content: "",
		})
	}
	return checker.adviceList, nil
}

type tableDisallowCreateTableAsChecker struct {
	adviceList []advisor.Advice
	level      advisor.Status
	title      string
	text       string
	line       int
}

// Enter implements the ast.Visitor interface.
func (v *tableDisallowCreateTableAsChecker) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.CreateTableStmt:
		if node.Select != nil {
			v.adviceList = append(v.adviceList, advisor.Advice{
				Status:  v.level,
				Code:    advisor.StatementCreateTableAs,
				Title:   v.title,
				Content: fmt.Sprintf("cannot create table `%s` by using CREATE TABLE ... [AS] SELECT ...", node.Table.Name.String()),
				Line:    node.OriginTextPosition(),
			})
		}
	default:
	}

	return in, false
}

// Leave implements the ast.Visitor interface.
func (*tableDisallowCreateTableAsChecker) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
