package model

import (
	"text/template"
	"strings"
	"unicode"
)



func TemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"lower":        strings.ToLower,
		"join":         func(ss []string, sep string) string { return strings.Join(ss, sep) },
		"title":        func(s string) string { if s == "" { return "" }; r := []rune(s); r[0] = unicode.ToUpper(r[0]); return string(r) },
		"toPascal":     toPascal,
		"protobufType": protobufType,
		"addIndex":     addIndex,
	}
}

func addIndex(index interface{}, offset int) int {
	switch v := index.(type) {
	case int:
		return v + offset
	case int64:
		return int(v) + offset
	default:
		return offset
	}
}

func protobufType(goType string) string {
	switch goType {
	case "string":
		return "string"
	case "int", "int32", "int64":
		return "int32"
	case "uint", "uint32", "uint64":
		return "uint32"
	case "bool":
		return "bool"
	case "float32", "float64":
		return "float"
	default:
		return "string" // fallback
	}
}

func toPascal(s string) string {
	if s == "" {
		return ""
	}
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	parts := strings.Fields(s)
	for i, part := range parts {
		if part == "" {
			continue
		}
		r := []rune(part)
		r[0] = unicode.ToUpper(r[0])
		parts[i] = string(r)
	}
	return strings.Join(parts, "")
}