package model

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/*
var tmplFS embed.FS

func GenerateCode(config *ModelConfig) error {
	dirs := []string{
		// filepath.Join(config.OutputPath, "internal", "type"),
		filepath.Join(config.OutputPath, "internal", "service"),
		filepath.Join(config.OutputPath, "internal", "service" , "dto"),
		filepath.Join(config.OutputPath, "internal", "api", "endpoints"),
		filepath.Join(config.OutputPath, "internal", "api", "transports"),
		filepath.Join(config.OutputPath, "internal", "api", "transports" , "http"),
		filepath.Join(config.OutputPath, "internal", "api", "transports" , "grpc"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	if err := generateModel(config); err != nil {
		return err
	}

	if err := generateProto(config); err != nil {
		return err
	}

	// if err := generateType(config); err != nil {
	// 	return err
	// }

	if err := generateRepository(config); err != nil {
		return err
	}

	if err := generateService(config); err != nil {
		return err
	}

	if err := generateEndpoint(config); err != nil {
		return err
	}

	if config.GenerateHTTP {
		if err := generateTransportHTTP(config); err != nil {
			return err
		}

		if err := generateTransportHTTPTest(config); err != nil {
			return err
		}
	}
	if config.GenerategRPC {
		if err := generateTransportgRPC(config); err != nil {
			return err
		}

		if err := generateTransportGRPCTest(config); err != nil {
			return err
		}
	}

	if err := generateRoutes(config); err != nil {
		return err
	}

	if config.GenerateTests {
		if err := generateServiceTest(config); err != nil {
			return err
		}
		if err := generateAPITest(config); err != nil {
			return err
		}
	}

	fmt.Printf("✅ Code generated successfully in %s\n", config.OutputPath)
	return nil
}

// func generateType(config *ModelConfig) error {
// 	root, err := GetProjectRoot()
// 	if err != nil {
// 		return err
// 	}
// 	tmplPath := filepath.Join(root, "pkg", "generator", "model", "templates", "type.go.tmpl")

// 	tmpl, err := template.New(filepath.Base(tmplPath)).Funcs(TemplateFuncMap()).ParseFiles(tmplPath)
// 	if err != nil {
// 		return err
// 	}

// 	f, err := os.Create(filepath.Join(config.OutputPath, "pkg", "type", strings.ToLower(config.ModelName)+".go"))
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	return tmpl.Execute(f, config)
// }

func generateService(config *ModelConfig) error {
	tmplContent, err := tmplFS.ReadFile("templates/service.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read embedded template service.go.tmpl: %w", err)
	}

	tmpl, err := template.New("service.go.tmpl").Funcs(TemplateFuncMap()).Parse(string(tmplContent))
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(config.OutputPath, "internal", "service", strings.ToLower(config.ModelName)+"_service.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, config) 
	if err != nil {
		return err
	}


	tmplDtoContent, err := tmplFS.ReadFile("templates/dto.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read embedded template dto.go.tmpl: %w", err)
	}

	tmplDto, err := template.New("dto.go.tmpl").Funcs(TemplateFuncMap()).Parse(string(tmplDtoContent))
	if err != nil {
		return err
	}

	fDto, err := os.Create(filepath.Join(config.OutputPath, "internal", "service" , "dto" , strings.ToLower(config.ModelName)+"_dto.go"))
	if err != nil {
		return err
	}
	defer fDto.Close()

	err = tmplDto.Execute(fDto, config)
	if err != nil {
		return err
	}

	return nil
}

func generateEndpoint(config *ModelConfig) error {
	tmplContent, err := tmplFS.ReadFile("templates/endpoint.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read embedded template endpoint.go.tmpl: %w", err)
	}

	tmpl, err := template.New("endpoint.go.tmpl").Funcs(TemplateFuncMap()).Parse(string(tmplContent))
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(config.OutputPath, "internal", "api", "endpoints", strings.ToLower(config.ModelName)+"_endpoint.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, config)
}

func generateTransportHTTP(config *ModelConfig) error {
	tmplContent, err := tmplFS.ReadFile("templates/transport_http.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read embedded template transport_http.go.tmpl: %w", err)
	}

	tmpl, err := template.New("transport_http.go.tmpl").Funcs(TemplateFuncMap()).Parse(string(tmplContent))
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(config.OutputPath, "internal", "api", "transports" , "http", strings.ToLower(config.ModelName)+"_http.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, config)
}

func generateTransportgRPC(config *ModelConfig) error {
	tmplContent, err := tmplFS.ReadFile("templates/transport_grpc.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read embedded template transport_grpc.go.tmpl: %w", err)
	}

	tmpl, err := template.New("transport_grpc.go.tmpl").Funcs(TemplateFuncMap()).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse gRPC transport template: %w", err)
	}

	transportsDir := filepath.Join(config.OutputPath, "internal", "api", "transports" , "grpc")
	os.MkdirAll(transportsDir, 0755)

	f, err := os.Create(filepath.Join(transportsDir, strings.ToLower(config.ModelName)+"_grpc.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, config)
}

func generateProto(config *ModelConfig) error {
	if !config.GenerateHTTP && !config.GenerategRPC {
		return nil
	}

	tmplContent, err := tmplFS.ReadFile("templates/proto.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read embedded template proto.go.tmpl: %w", err)
	}

	tmpl, err := template.New("proto.go.tmpl").Funcs(TemplateFuncMap()).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse proto template: %w", err)
	}

	protoDir := filepath.Join(config.OutputPath, "api", "proto", "v1")
	os.MkdirAll(protoDir, 0755)

	f, err := os.Create(filepath.Join(protoDir, strings.ToLower(config.ModelName)+".proto"))
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, config)
}

func generateRoutes(config *ModelConfig) error {
	path := filepath.Join(config.OutputPath, "internal", "api" , "transports" , "http", "routes.go")
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		tmplContent, err := tmplFS.ReadFile("templates/routes.go.tmpl")
		if err != nil {
			return fmt.Errorf("failed to read embedded template routes.go.tmpl: %w", err)
		}

		tmpl, err := template.New("routes.go.tmpl").Funcs(TemplateFuncMap()).Parse(string(tmplContent))
		if err != nil {
			return err
		}
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
		return tmpl.Execute(f, config)
	} else {
		fmt.Println("⚠️  routes.go already exists — manual update required for now.")
		return nil
	}
}

func generateServiceTest(config *ModelConfig) error {
	tmplContent, err := tmplFS.ReadFile("templates/service_test.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read embedded template service_test.go.tmpl: %w", err)
	}

	tmpl, err := template.New("service_test.go.tmpl").Funcs(TemplateFuncMap()).Parse(string(tmplContent))
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(config.OutputPath, "internal", "service", strings.ToLower(config.ModelName)+"_service_test.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, config)
}

func generateAPITest(config *ModelConfig) error {
	tmplContent, err := tmplFS.ReadFile("templates/api_test.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read embedded template api_test.go.tmpl: %w", err)
	}

	tmpl, err := template.New("api_test.go.tmpl").Funcs(TemplateFuncMap()).Parse(string(tmplContent))
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(config.OutputPath, "internal", "api", "endpoints", strings.ToLower(config.ModelName)+"_endpoint_test.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, config)
}

func generateModel(config *ModelConfig) error {
	tmplContent, err := tmplFS.ReadFile("templates/model.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read embedded template model.go.tmpl: %w", err)
	}

	tmpl, err := template.New("model.go.tmpl").Funcs(TemplateFuncMap()).Parse(string(tmplContent))
	if err != nil {
		return err
	}

	modelsDir := filepath.Join(config.OutputPath, "internal", "models")
	os.MkdirAll(modelsDir, 0755)

	f, err := os.Create(filepath.Join(modelsDir, strings.ToLower(config.ModelName)+".go"))
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, config)
}

func generateRepository(config *ModelConfig) error {
	tmplContent, err := tmplFS.ReadFile("templates/repository.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read embedded template repository.go.tmpl: %w", err)
	}

	tmpl, err := template.New("repository.go.tmpl").Funcs(TemplateFuncMap()).Parse(string(tmplContent))
	if err != nil {
		return err
	}

	repoDir := filepath.Join(config.OutputPath, "internal", "repositories")
	os.MkdirAll(repoDir, 0755)

	f, err := os.Create(filepath.Join(repoDir, strings.ToLower(config.ModelName)+"_repository.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, config)
}

func generateTransportGRPCTest(config *ModelConfig) error {
	if !config.GenerategRPC {
		return nil
	}

	tmplContent, err := tmplFS.ReadFile("templates/transport_grpc_test.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read embedded template transport_grpc_test.go.tmpl: %w", err)
	}

	tmpl, err := template.New("transport_grpc_test.go.tmpl").Funcs(TemplateFuncMap()).Parse(string(tmplContent))
	if err != nil {
		return err
	}

	transportsDir := filepath.Join(config.OutputPath, "internal", "api", "transports" , "grpc")
	os.MkdirAll(transportsDir, 0755)

	f, err := os.Create(filepath.Join(transportsDir, strings.ToLower(config.ModelName)+"_grpc_test.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, config)
}

func generateTransportHTTPTest(config *ModelConfig) error {
	if !config.GenerateHTTP {
		return nil
	}

	tmplContent, err := tmplFS.ReadFile("templates/transport_http_test.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read embedded template transport_http_test.go.tmpl: %w", err)
	}

	tmpl, err := template.New("transport_http_test.go.tmpl").Funcs(TemplateFuncMap()).Parse(string(tmplContent))
	if err != nil {
		return err
	}

	transportsDir := filepath.Join(config.OutputPath, "internal", "api", "transports", "http")
	os.MkdirAll(transportsDir, 0755)

	f, err := os.Create(filepath.Join(transportsDir, strings.ToLower(config.ModelName)+"_http_test.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, config)
}