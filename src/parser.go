package main

import (
	"fmt"
)

//TODO: use a struct may be better
var currentToken int
var currentLexeme int
var parseErr error

func Parse() {
	currentToken = GetToken()
	parseStmtSequence()
}

func parseStmtSequence() {
	for parseErr == nil && currentToken != tokenEOF {
		switch currentToken {
		case tokenKeywordIf:
			match(tokenKeywordIf)
			parseIfStmt()

		case tokenKeywordRead:
			parseRepeatStmt()

		case tokenKeywordRead:
			match(tokenKewordRead)
			parseReadStmt()

		case tokenKeywordWrite:
			match(tokenKeywordWrite)
			parseWriteStmt()

		case tokenId:
			match(tokenId)
			if !match(tokenAssign) {
				parseErr = -1
				return
			}
			parseAssignStmt()
		default:
			//TODO: error
		}
	}
}

func parseStmt() {

}

func parseIfStmt() {
	parseExp()
	switch currentToken {
	case tokenKeywordThen:
	case tokenKeywordElse:
	default:
	}
}

func parseRepeatStmt() {

}

func parseAssignStmt() {
	parseExp()
}

func parseReadStmt() {
	switch currentToken {
	case tokenId:
		match(tokenId)
	default:
		//TODO: error
	}
}

func parseWriteStmt() {

}

func parseExp() {
	parseSimpleExp()
	switch currentToken {

	}
}

func parseSimpleExp() {
	parseTerm()
}

func parseTerm() {
	parseFactor()
}

func parseFactor() {
	switch currentToken {
	case tokenLeftParen:
		match(tokenLeftParen)
		parseExp()
	}
}

func match(token int) bool {
	if token == currentToken {
		currentToken = GetToken()
		return true
	}
	return false
}
