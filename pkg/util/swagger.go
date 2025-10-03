package util

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Função principal para mesclar o Swagger a partir do index.yaml
func ParseSwagger(indexPath string) {
	basePath := filepath.Dir(indexPath)
	swagger, err := loadYAMLFile(indexPath)
	if err != nil {
		fmt.Printf("Erro ao carregar index.yaml: %v\n", err)
		return
	}

	// Resolve referências recursivamente
	if err = resolveRefs(swagger, basePath); err != nil {
		fmt.Printf("Erro ao resolver referências: %v\n", err)
		return
	}

	// Gera o arquivo final consolidado
	output, err := yaml.Marshal(swagger)
	if err != nil {
		fmt.Printf("Erro ao serializar YAML consolidado: %v\n", err)
		return
	}

	if err = os.WriteFile(filepath.Join(basePath, "swagger.yaml"), output, 0644); err != nil {
		fmt.Printf("Erro ao salvar swagger.yaml: %v\n", err)
	}
}

// Função para carregar um arquivo YAML como map
func loadYAMLFile(filePath string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler arquivo %s: %v", filePath, err)
	}

	var content map[string]interface{}
	if err = yaml.Unmarshal(data, &content); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse do YAML %s: %v", filePath, err)
	}

	return content, nil
}

// Função recursiva para resolver referências $ref no YAML
func resolveRefs(content map[string]interface{}, basePath string) error {
	for _, value := range content {
		if v, ok := value.(map[string]interface{}); ok {
			if ref, exists := v["$ref"]; exists {
				refPath := filepath.Join(basePath, ref.(string))

				// Carrega o arquivo referenciado
				resolved, err := loadYAMLFile(refPath)
				if err != nil {
					return fmt.Errorf("falha ao resolver referência %s: %v", refPath, err)
				}

				// Substitui a referência pelo conteúdo resolvido
				delete(v, "$ref")
				for k, val := range resolved {
					v[k] = val
				}
			} else {
				// Resolução recursiva para sub-níveis
				if err := resolveRefs(v, basePath); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
