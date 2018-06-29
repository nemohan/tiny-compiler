package main

import (
	"errors"
	"fmt"
)

//TODO: use a struct may be better
//var currentLexeme int
var currentToken *tokenSymbol
var parseErr error

func Parse() {
	currentToken = GetToken()
	fmt.Printf("current token:%v\n", currentToken)
	parseStmtSequence()
}

func parseStmtSequence() {
	for parseErr == nil && currentToken.tokenType != tokenEOF {
		switch currentToken.tokenType {
		case tokenId:
			match(tokenId)
			if !match(tokenAssign) {
				parseErr = errors.New(" tokenAssign")
				return
			}
			parseAssignStmt()
			//don't care match or not
			matchStr(";")
		default:
			//TODO: error
			if matchStr("if") {
				fmt.Printf("match if\n")
				parseIfStmt()
				matchStr(";")
			}

			if matchStr("read") {
				parseReadStmt()
				matchStr(";")
				fmt.Printf("parse read done\n")
			}
			if matchStr("repeat") {

			}
			if matchStr("write") {

			}
		}
	}
}

func parseStmt() {

}

func parseIfStmt() {
	parseExp()
	if !matchStr("then") {
		//TODO:
	}
	fmt.Printf("match then\n")
	parseExp()

	//optional
	if matchStr("else") {
		parseStmtSequence()
	}
	if !matchStr("end") {

	}
}

func parseRepeatStmt() {

}

func parseAssignStmt() {
	fmt.Printf("parse assignStmt\n")
	parseExp()
}

func parseReadStmt() {
	switch currentToken.tokenType {
	case tokenId:
		fmt.Printf("parseReadStmt:%v\n", currentToken)
		match(tokenId)
	default:
		//TODO: error
	}
}

func parseWriteStmt() {

}

func parseExp() {
	parseSimpleExp()
	tokenType := currentToken.tokenType
	for tokenType == tokenLess || tokenType == tokenEqual {
		switch tokenLess {
		case tokenLess:
			match(tokenLess)
			parseSimpleExp()

		case tokenEqual:
			match(tokenEqual)
			parseSimpleExp()
		default:
			panic("parseExp")
		}
	}
}

func parseSimpleExp() {
	parseTerm()
	tokenType := currentToken.tokenType
	for tokenType == tokenAdd || tokenType == tokenMinus {
		switch tokenAdd {
		case tokenAdd:
			match(tokenAdd)
			parseSimpleExp()

		case tokenMinus:
			match(tokenMinus)
			parseSimpleExp()
		default:
			panic("parseSimpleExp")
		}
	}

}

func parseTerm() {
	parseFactor()
	tokenType := currentToken.tokenType
	for tokenType == tokenMultiply || tokenType == tokenDiv {
		switch tokenType {
		case tokenMultiply:
			match(tokenMultiply)
			parseFactor()
		case tokenDiv:
			match(tokenDiv)
			parseFactor()
		}
		tokenType = currentToken.tokenType
	}
}

func parseFactor() {
	switch currentToken.tokenType {
	case tokenLeftParen:
		match(tokenLeftParen)
		parseExp()
		if !match(tokenRightParen) {

		}
	case tokenNumber:
		fmt.Printf("match number:%s\n", currentToken.lexeme)
		match(tokenNumber)

	case tokenId:
		match(tokenId)
		fmt.Printf("match id:%s\n", currentToken.lexeme)
	default:
		fmt.Printf("parse factor default\n")

	}
}

func matchStr(lexeme string) bool {
	if lexeme == currentToken.lexeme {
		currentToken = GetToken()
		return true
	}
	return false
}
func match(token int) bool {
	if token == currentToken.tokenType {
		currentToken = GetToken()
		return true
	}
	return false
}
