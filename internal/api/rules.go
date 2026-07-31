package api

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/chankei613/ai-output-validator/internal/db"
)

// evaluateRule はAI出力(output)が1つのRuleを満たすかを判定する。
// 全ルールは決定論的（LLM-as-judge等の非決定的採点はv0.1.0スコープ外）。
func evaluateRule(rule db.Rule, output string) db.RuleResult {
	res := db.RuleResult{Type: rule.Type, Value: rule.Value}

	switch rule.Type {
	case "contains":
		res.Passed = strings.Contains(output, rule.Value)
		if !res.Passed {
			res.Message = fmt.Sprintf("output does not contain %q", rule.Value)
		}

	case "not_contains":
		res.Passed = !strings.Contains(output, rule.Value)
		if !res.Passed {
			res.Message = fmt.Sprintf("output unexpectedly contains %q", rule.Value)
		}

	case "regex":
		re, err := regexp.Compile(rule.Value)
		if err != nil {
			res.Message = fmt.Sprintf("invalid regex: %s", err)
			break
		}
		res.Passed = re.MatchString(output)
		if !res.Passed {
			res.Message = fmt.Sprintf("output does not match /%s/", rule.Value)
		}

	case "min_length":
		n, err := strconv.Atoi(rule.Value)
		if err != nil {
			res.Message = fmt.Sprintf("invalid min_length value: %s", rule.Value)
			break
		}
		res.Passed = len(output) >= n
		if !res.Passed {
			res.Message = fmt.Sprintf("output length %d is below minimum %d", len(output), n)
		}

	case "max_length":
		n, err := strconv.Atoi(rule.Value)
		if err != nil {
			res.Message = fmt.Sprintf("invalid max_length value: %s", rule.Value)
			break
		}
		res.Passed = len(output) <= n
		if !res.Passed {
			res.Message = fmt.Sprintf("output length %d exceeds maximum %d", len(output), n)
		}

	case "json_valid":
		res.Passed = json.Valid([]byte(output))
		if !res.Passed {
			res.Message = "output is not valid JSON"
		}

	case "json_key_exists":
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(output), &parsed); err != nil {
			res.Message = "output is not a valid JSON object"
			break
		}
		_, res.Passed = parsed[rule.Value]
		if !res.Passed {
			res.Message = fmt.Sprintf("JSON object is missing key %q", rule.Value)
		}

	default:
		res.Message = fmt.Sprintf("unknown rule type %q", rule.Type)
	}

	return res
}

// evaluateCase は1テストケース分のRule群をAND評価する。全ルールが合格なら合格。
func evaluateCase(rules []db.Rule, output string) (bool, []db.RuleResult) {
	results := make([]db.RuleResult, 0, len(rules))
	passed := true
	for _, rule := range rules {
		rr := evaluateRule(rule, output)
		results = append(results, rr)
		if !rr.Passed {
			passed = false
		}
	}
	return passed, results
}
