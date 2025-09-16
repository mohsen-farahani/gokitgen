package model

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Field struct {
	Name           string
	Type           string
	TypeIsEnum     bool
	TypeIsRelation bool
	IsNullable     bool
	Validation     []string
	GormTag        string // e.g., "default:0", "index", "unique"
	Comment        string
}

type Enum struct {
	Name   string
	Values []string
}

type ModelConfig struct {
	ModelName      string
	ModulePath 	   string
	Fields         []Field
	Enums          []Enum
	GenerateHTTP   bool
	GenerategRPC   bool
	GenerateTests  bool
	OutputPath     string
}

func RunWizard() *ModelConfig {
	reader := bufio.NewReader(os.Stdin)
	config := &ModelConfig{}

	fmt.Print("📝 Enter model name (e.g., Order): ")
	modelName, _ := reader.ReadString('\n')
	config.ModelName = strings.TrimSpace(modelName)
	if config.ModelName == "" {
		fmt.Println("❌ Model name is required.")
		os.Exit(1)
	}

	fmt.Print("📦 Enter module path (e.g., github.com/your_project): ")
	modulePath, _ := reader.ReadString('\n')
	config.ModulePath = strings.TrimSpace(modulePath)
	if config.ModulePath == "" {
		fmt.Println("⚠️  Module path is required for imports. Using 'your-module' as fallback.")
		config.ModulePath = "your-module"
	}

	config.Enums = askEnums(reader)

	config.Fields = askFields(reader, config.Enums)

	config.GenerateHTTP, config.GenerategRPC = askTransportType(reader)

	config.GenerateTests = askGenerateTests(reader)

	// you can use ./generated
	config.OutputPath = "./"

	return config
}

func askEnums(reader *bufio.Reader) []Enum {
	var enums []Enum

	for {
		fmt.Print("🎨 Add enum? (y/n): ")
		yn, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(yn)) != "y" {
			break
		}

		fmt.Print("  Enum name (e.g., OrderStatus): ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)
		if name == "" {
			fmt.Println("⚠️  Enum name is required. Skipping.")
			continue
		}

		var values []string
		for {
			fmt.Print("  Add value (e.g., PENDING) or press Enter to finish: ")
			val, _ := reader.ReadString('\n')
			val = strings.TrimSpace(val)
			if val == "" {
				break
			}
			values = append(values, val)
		}

		if len(values) == 0 {
			fmt.Println("⚠️  At least one value is required. Skipping enum.")
			continue
		}

		enums = append(enums, Enum{Name: name, Values: values})
	}

	return enums
}

func askFields(reader *bufio.Reader, enums []Enum) []Field {
	var fields []Field

	enumNames := make(map[string]bool)
	for _, e := range enums {
		enumNames[e.Name] = true
	}

	for {
		fmt.Print("➕ Add field? (y/n): ")
		yn, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(yn)) != "y" {
			break
		}

		fmt.Print("  Field name (e.g., Side, Market, Amount): ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)
		if name == "" {
			fmt.Println("⚠️  Field name is required. Skipping.")
			continue
		}

		fmt.Print("  Field type (e.g., string, int, uint, bool, or enum name like OrderStatus, or Ref:Market for relation): ")
		typ, _ := reader.ReadString('\n')
		typ = strings.TrimSpace(typ)
		if typ == "" {
			fmt.Println("⚠️  Field type is required. Skipping.")
			continue
		}

		isEnum := enumNames[typ]
		isRelation := false

		if strings.HasPrefix(typ, "Ref:") {
			isRelation = true
			typ = strings.TrimPrefix(typ, "Ref:")
		}

		fmt.Print("  Nullable? (y/n): ")
		nullable, _ := reader.ReadString('\n')
		isNullable := strings.TrimSpace(strings.ToLower(nullable)) == "y"

		var validations []string
		for {
			fmt.Print("  Add validation? (e.g., required, email, min=1, or press Enter to skip): ")
			val, _ := reader.ReadString('\n')
			val = strings.TrimSpace(val)
			if val == "" {
				break
			}
			validations = append(validations, val)
		}

		var gormTag string
		if !isRelation && !isEnum {
			fmt.Print("  Add GORM tag? (e.g., default:0, index, unique, or press Enter to skip): ")
			tag, _ := reader.ReadString('\n')
			gormTag = strings.TrimSpace(tag)
		}

		fmt.Print("  Add comment? (e.g., Order side type, or press Enter to skip): ")
		comment, _ := reader.ReadString('\n')
		comment = strings.TrimSpace(comment)

		fields = append(fields, Field{
			Name:           name,
			Type:           typ,
			TypeIsEnum:     isEnum,
			TypeIsRelation: isRelation,
			IsNullable:     isNullable,
			Validation:     validations,
			GormTag:        gormTag,
			Comment:        comment,
		})
	}

	return fields
}

func askTransportType(reader *bufio.Reader) (bool, bool) {
	fmt.Print("🌐 Generate API for (1=HTTP, 2=gRPC, 3=Both): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		return true, false
	case "2":
		return false, true
	case "3":
		return true, true
	default:
		fmt.Println("⚠️  Invalid choice. Defaulting to HTTP.")
		return true, false
	}
}

func askGenerateTests(reader *bufio.Reader) bool {
	fmt.Print("🧪 Generate tests? (y/n): ")
	yn, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(yn)) == "y"
}