package util

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Validações de caracteres para os campos fornecidos
var (
	COD_01 = regexp.MustCompile("^[a-zA-Z0-9]+$")                                                                        // Letras e números
	COD_02 = regexp.MustCompile("^[0-9]+$")                                                                              // Apenas números
	COD_03 = regexp.MustCompile("^[a-zA-Z0-9@._-]+$")                                                                    // Letras, números, @, pontos, underlines e hífens
	COD_04 = regexp.MustCompile("^[0-9./-]+$")                                                                           // Números, pontos, barras e hífens
	COD_05 = regexp.MustCompile("^([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$")                                           // Formato de hora HH:MM:SS
	COD_06 = regexp.MustCompile(`^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01]) (0[0-9]|1\d|2[0-3]):([0-5]\d):([0-5]\d)$`) // Formato "YYYY-MM-DD HH:MM:SS"
	COD_07 = regexp.MustCompile(`^(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]).*$`)                    // Senha: pelo menos 1 maiúscula, 1 número e 1 caractere especial
)

// FormatarErroValidacao trata o erro das validações dos campos do body
func FormatarErroValidacao(err error) error {

	mensagensErro := map[string]string{
		"cannot be blank":           "Campos obrigatórios estando vazios",
		"must be in a valid format": "Campos com caracteres especiais ou com formato inválido",
	}

	camposInvalidos := make(map[string][]string)

	if tipoErro, ok := err.(validation.Errors); ok {
		for campo, erro := range tipoErro {
			erroMsg := erro.Error()
			for substring, mensagem := range mensagensErro {
				if strings.Contains(erro.Error(), substring) {
					erroMsg = mensagem
					break
				}
			}
			camposInvalidos[erroMsg] = append(camposInvalidos[erroMsg], campo)
		}
	}

	var builder strings.Builder
	for mensagem, campos := range camposInvalidos {
		builder.WriteString(fmt.Sprintf("%s: %s; ", mensagem, strings.Join(campos, ", ")))
	}

	mensagemFinal := strings.TrimSuffix(builder.String(), "; ")
	return errors.New(mensagemFinal)
}
