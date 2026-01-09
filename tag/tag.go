package tag

import (
	"slices"
	"strings"
)

type Rule struct {
	Name   string
	Params []string
}

type Config struct {
	RuleSplitter       rune
	NameParamSeparator rune
	ParamsSeparator    rune
}

func DefaultConfig() Config {
	return Config{
		RuleSplitter:       '|',
		NameParamSeparator: '=',
		ParamsSeparator:    ',',
	}
}

type Parser struct {
	cfg Config
}

func NewParser(cfg Config) *Parser {
	if cfg == (Config{}) {
		cfg = DefaultConfig()
	}
	return &Parser{cfg: cfg}
}

func Parse(tag string) []Rule {
	return NewParser(DefaultConfig()).Parse(tag)
}

func (p *Parser) Parse(tag string) []Rule {
	if tag == "" {
		return nil
	}

	rules := make([]Rule, 0)
	currentRule := ""
	inParam := false

	for i := 0; i < len(tag); i++ {
		ch := tag[i]

		if ch == byte(p.cfg.NameParamSeparator) {
			inParam = true
			currentRule += string(ch)
		} else if ch == byte(p.cfg.RuleSplitter) {
			if inParam {
				nextPart := ""
				for j := i + 1; j < len(tag); j++ {
					if tag[j] == byte(p.cfg.RuleSplitter) {
						break
					}
					nextPart += string(tag[j])
				}

				if !slices.Contains([]byte(nextPart), byte(p.cfg.NameParamSeparator)) && !isValidatorName(nextPart) {
					currentRule += string(ch)
				} else {
					inParam = false
					if currentRule != "" {
						rules = append(rules, p.parseRule(currentRule))
						currentRule = ""
					}
				}
			} else {
				if currentRule != "" {
					rules = append(rules, p.parseRule(currentRule))
					currentRule = ""
				}
			}
		} else {
			currentRule += string(ch)
		}
	}

	if currentRule != "" {
		rules = append(rules, p.parseRule(currentRule))
	}

	return rules
}

func (p *Parser) parseRule(ruleStr string) Rule {
	ruleStr = strings.TrimSpace(ruleStr)

	if before, after, ok := strings.Cut(ruleStr, string(p.cfg.NameParamSeparator)); ok {
		name := strings.TrimSpace(before)
		raw := strings.TrimSpace(after)
		parts := []string{}
		if raw != "" {
			for _, param := range strings.Split(raw, string(p.cfg.ParamsSeparator)) {
				tp := strings.TrimSpace(param)
				if tp != "" {
					parts = append(parts, tp)
				}
			}
		}
		return Rule{
			Name:   name,
			Params: parts,
		}
	}

	return Rule{
		Name:   ruleStr,
		Params: []string{},
	}
}

func isValidatorName(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}

	for i, ch := range s {
		if i == 0 {
			if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')) {
				return false
			}
		} else {
			if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_') {
				return false
			}
		}
	}

	return true
}
