package dsl

import (
	"fmt"
)

type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
	errors    []string
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParsePipeline() (Node, error) {
	left, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}

	// If we are at EOF, we are done
	if p.curToken.Type == EOF {
		return left, nil
	}
    // If not EOF, check if we have unconsumed tokens that indicate error
    if p.curToken.Type != EOF {
         // This can happen if precedence stopped parsing but we are not at EOF.
         // E.g. "cmd arg )"
    }

	return left, nil
}

// Precedences
const (
	LOWEST      = iota
	PIPELINE    // |
	BINARY      // ^
	CALL        // command arg
)

func (p *Parser) parseExpression(precedence int) (Node, error) {
	prefix := p.prefixParseFn()
	if prefix == nil {
		return nil, fmt.Errorf("no prefix parse function for %s", p.curToken.Literal)
	}

	leftExp, err := prefix()
	if err != nil {
		return nil, err
	}

	for p.curToken.Type != EOF && precedence < p.peekPrecedence() {
		infix := p.infixParseFn(p.peekToken.Type)
		if infix == nil {
			return leftExp, nil
		}

		p.nextToken()
		leftExp, err = infix(leftExp)
		if err != nil {
			return nil, err
		}
	}

	return leftExp, nil
}

type prefixParseFn func() (Node, error)
type infixParseFn func(Node) (Node, error)

func (p *Parser) prefixParseFn() prefixParseFn {
	switch p.curToken.Type {
	case IDENT:
		return p.parseCommand
	case LPAREN:
		return p.parseGroupedExpression
	default:
		return nil
	}
}

func (p *Parser) infixParseFn(t TokenType) infixParseFn {
	switch t {
	case PIPE:
		return p.parsePipelineOp
	case CARET:
		return p.parseBinaryOp
	}
	return nil
}

func (p *Parser) peekPrecedence() int {
	switch p.peekToken.Type {
	case PIPE:
		return PIPELINE
	case CARET:
		return BINARY
	}
	return LOWEST
}

func (p *Parser) parseCommand() (Node, error) {
	cmd := &CommandNode{Name: p.curToken.Literal}

	for p.peekTokenIsArg() {
		p.nextToken()
		arg, err := p.parseArg()
		if err != nil {
			return nil, err
		}
		cmd.Args = append(cmd.Args, arg)
	}
	return cmd, nil
}

func (p *Parser) peekTokenIsArg() bool {
	t := p.peekToken.Type
	return t == IDENT || t == NUMBER || t == STRING || t == LPAREN || t == EQUALS
}

func (p *Parser) parseArg() (ArgNode, error) {
	if p.curToken.Type == IDENT && p.peekToken.Type == EQUALS {
		key := p.curToken.Literal
		p.nextToken() // eat key (cur is now EQUALS)
		p.nextToken() // eat = (cur is now value start)

		val, err := p.parseArgValue()
		if err != nil {
			return nil, err
		}
		return &KeyValueNode{Key: key, Value: val}, nil
	}

	return p.parseArgValue()
}

func (p *Parser) parseArgValue() (ArgNode, error) {
	if p.curToken.Type == LPAREN {
		p.nextToken() // eat (
		exp, err := p.parseExpression(LOWEST) // recurse
		if err != nil {
			return nil, err
		}
		if p.peekToken.Type != RPAREN {
			// consume token to check if it matches RPAREN?
            // current token is the last token of expression. peekToken should be RPAREN.
			return nil, fmt.Errorf("expected ) but got %v", p.peekToken)
		}
		p.nextToken() // eat )
		return &SubExpressionNode{Node: exp}, nil
	}

	// Literal
	return &LiteralNode{Value: p.curToken.Literal}, nil
}

func (p *Parser) parsePipelineOp(left Node) (Node, error) {
	// A | B | C
    p.nextToken() // Advance past the pipe operator!
	right, err := p.parseExpression(PIPELINE) // use PIPELINE precedence
	if err != nil {
		return nil, err
	}

	if pipe, ok := left.(*PipelineNode); ok {
		pipe.Nodes = append(pipe.Nodes, right)
		return pipe, nil
	}

	return &PipelineNode{Nodes: []Node{left, right}}, nil
}

func (p *Parser) parseBinaryOp(left Node) (Node, error) {
	op := p.curToken.Literal
	precedence := p.curPrecedence()
	p.nextToken() // Advance past operator
	right, err := p.parseExpression(precedence)
	if err != nil {
		return nil, err
	}
	return &BinaryNode{Left: left, Operator: op, Right: right}, nil
}

func (p *Parser) curPrecedence() int {
	switch p.curToken.Type {
	case PIPE:
		return PIPELINE
	case CARET:
		return BINARY
	}
	return LOWEST
}

func (p *Parser) parseGroupedExpression() (Node, error) {
	p.nextToken()
	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	if p.peekToken.Type != RPAREN {
		return nil, fmt.Errorf("expected )")
	}
	p.nextToken()
	return &GroupNode{Inner: exp}, nil
}
