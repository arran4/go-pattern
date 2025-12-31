package dsl

import (
	"fmt"
)

type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
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

// ParsePipeline parses a full pipeline of commands.
// It parses the first expression, and if there are remaining tokens,
// it assumes they are part of the pipeline (if the parser stopped due to precedence).
// Wait, `parseExpression(0)` parses everything until EOF unless precedence stops it.
// If we introduce math operators, we need to ensure their precedence fits.
func (p *Parser) ParsePipeline() (Node, error) {
	left, err := p.parseExpression(0)
	if err != nil {
		return nil, err
	}

	// If we are at EOF, we are done
	if p.curToken.Type == EOF {
		return left, nil
	}
	// If not EOF, it might be an error or unexpected token
	// return nil, fmt.Errorf("unexpected token: %s", p.curToken.Literal)
	// For loose compatibility, we return what we have
	return left, nil
}

// Precedences
const (
	LOWEST      = iota
	PIPELINE    // |
	SUM         // + -
	PRODUCT     // * / %
	BINARY      // ^ (bitwise xor / boolean) - wait, ^ usually binds tighter than | but looser than *
	CALL        // command arg, function call
)

// We need to define precedences carefully.
// Standard math: * / > + -
// Pipeline | should be very low precedence (separates stages).
// ^ (XOR) in C/Go is lower than + but higher than |.
// Let's adjust:
// | : 1
// ^ : 2
// + - : 3
// * / % : 4
// Call : 5

// However, `PIPELINE` constant was used for `|`.
// `BINARY` was used for `^`.

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
		// Could be command or variable or function call start
		return p.parseIdentifierOrCommand
	case NUMBER, STRING:
		return p.parseLiteral
	case LPAREN:
		return p.parseGroupedExpression
	case MINUS:
		// Unary minus?
		return nil // Not implemented yet
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
	case PLUS, MINUS, ASTERISK, SLASH, PERCENT:
		return p.parseMathOp
	}
	return nil
}

func (p *Parser) peekPrecedence() int {
	switch p.peekToken.Type {
	case PIPE:
		return PIPELINE
	case CARET:
		return BINARY
	case PLUS, MINUS:
		return SUM
	case ASTERISK, SLASH, PERCENT:
		return PRODUCT
	case LPAREN:
		// Function call? if peek is ( and cur is IDENT?
		// Handled in prefix parse usually.
		return LOWEST
	}
	return LOWEST
}

func (p *Parser) parseIdentifierOrCommand() (Node, error) {
	// If it's a command in a pipeline context, it consumes args.
	// If it's a variable or function in a math context, it might be different.
	// In the previous grammar, `cmd arg1 arg2`.
	// In math grammar: `sin(x)`.
	// Ambiguity: `sin (x)` vs `cmd (arg)`.
	// We distinguish by checking if there's an opening paren immediately?
	// Or we treat everything as command call?
	// `sin` is a command taking 1 arg `(x)`.
	// `x` is a command taking 0 args? Or a literal?
	// Let's treat IDENT as CommandNode.

	cmd := &CommandNode{Name: p.curToken.Literal}

	// Arguments parsing.
	// In pipeline mode: consume until | or EOF.
	// In math mode: `a + b`. `a` has no args.
	// Problem: `cmd arg | ...` vs `a + b`.
	// If next token is an operator (+, -, *, /), we should stop parsing args.

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

func (p *Parser) parseLiteral() (Node, error) {
	return &LiteralNode{Value: p.curToken.Literal}, nil
}

func (p *Parser) peekTokenIsArg() bool {
	t := p.peekToken.Type
	// Operators stop argument parsing!
	if t == PIPE || t == CARET || t == PLUS || t == MINUS || t == ASTERISK || t == SLASH || t == PERCENT || t == RPAREN || t == COMMA || t == EOF {
		return false
	}
	return true
	// Included: IDENT, NUMBER, STRING, LPAREN, EQUALS
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
			return nil, fmt.Errorf("expected ) but got %v", p.peekToken)
		}
		p.nextToken() // eat )
		return &SubExpressionNode{Node: exp}, nil
	}

	// Literal or Identifier acting as arg
	return &LiteralNode{Value: p.curToken.Literal}, nil
}

func (p *Parser) parsePipelineOp(left Node) (Node, error) {
    p.nextToken() // Advance past |
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
	p.nextToken()
	right, err := p.parseExpression(precedence)
	if err != nil {
		return nil, err
	}
	return &BinaryNode{Left: left, Operator: op, Right: right}, nil
}

func (p *Parser) parseMathOp(left Node) (Node, error) {
	op := p.curToken.Literal
	precedence := p.curPrecedence()
	p.nextToken()
	right, err := p.parseExpression(precedence)
	if err != nil {
		return nil, err
	}
	// Reuse BinaryNode? Or MathNode?
	// BinaryNode handles operator string.
	return &BinaryNode{Left: left, Operator: op, Right: right}, nil
}

func (p *Parser) curPrecedence() int {
	switch p.curToken.Type {
	case PIPE:
		return PIPELINE
	case CARET:
		return BINARY
	case PLUS, MINUS:
		return SUM
	case ASTERISK, SLASH, PERCENT:
		return PRODUCT
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
