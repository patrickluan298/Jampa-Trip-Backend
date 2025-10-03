package util

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

// Validações de caracteres para os campos fornecidos
var (
	COD_01 = regexp.MustCompile(`^[a-zA-Z0-9]+$`)                                                                        // Letras e números
	COD_02 = regexp.MustCompile(`^[0-9]+$`)                                                                              // Apenas números
	COD_03 = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)                                      // Letras, números, @, pontos, underlines e hífens
	COD_04 = regexp.MustCompile(`^[0-9./-]+$`)                                                                           // Números, pontos, barras e hífens
	COD_05 = regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$`)                                           // Formato de hora HH:MM:SS
	COD_06 = regexp.MustCompile(`^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01]) (0[0-9]|1\d|2[0-3]):([0-5]\d):([0-5]\d)$`) // Formato "YYYY-MM-DD HH:MM:SS"
	COD_07 = regexp.MustCompile(`^.{8,}$`)                                                                               // Apenas comprimento mínimo para login
	COD_08 = regexp.MustCompile(`[A-Z]`)                                                                                 // Pelo menos 1 maiúscula
	COD_09 = regexp.MustCompile(`\d`)                                                                                    // Pelo menos 1 número
	COD_10 = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)                                                 // Pelo menos 1 caractere especial
	COD_11 = regexp.MustCompile(`^\(\d{2}\)\d{4,5}-\d{4}$`)                                                              // Formato de telefone (XX)XXXXX-XXXX ou (XX)XXXX-XXXX
	COD_12 = regexp.MustCompile(`^\d{2}\.\d{3}\.\d{3}\/\d{4}-\d{2}$`)                                                    // Formato de CNPJ XX.XXX.XXX/XXXX-XX
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

// ValidaSegurancaSenha - valida se a senha é forte
func ValidaSegurancaSenha(senha string) error {
	if !COD_08.MatchString(senha) {
		return errors.New("senha deve conter pelo menos 1 letra maiúscula")
	}
	if !COD_09.MatchString(senha) {
		return errors.New("senha deve conter pelo menos 1 número")
	}
	if !COD_10.MatchString(senha) {
		return errors.New("senha deve conter pelo menos 1 caractere especial")
	}
	return nil
}

// ValidaCPF - valida se o CPF é válido (formato e dígitos verificadores)
func ValidaCPF(cpf string) error {
	cpf = strings.ReplaceAll(cpf, ".", "")
	cpf = strings.ReplaceAll(cpf, "-", "")

	if len(cpf) != 11 {
		return errors.New("CPF deve ter 11 dígitos")
	}

	if strings.Count(cpf, string(cpf[0])) == 11 {
		return errors.New("CPF inválido: todos os dígitos são iguais")
	}

	// Validação do primeiro dígito verificador
	soma := 0
	for i := 0; i < 9; i++ {
		digito, _ := strconv.Atoi(string(cpf[i]))
		soma += digito * (10 - i)
	}
	resto := soma % 11
	primeiroDigito := 0
	if resto >= 2 {
		primeiroDigito = 11 - resto
	}

	digito1, _ := strconv.Atoi(string(cpf[9]))
	if digito1 != primeiroDigito {
		return errors.New("CPF inválido: primeiro dígito verificador incorreto")
	}

	// Validação do segundo dígito verificador
	soma = 0
	for i := 0; i < 10; i++ {
		digito, _ := strconv.Atoi(string(cpf[i]))
		soma += digito * (11 - i)
	}
	resto = soma % 11
	segundoDigito := 0
	if resto >= 2 {
		segundoDigito = 11 - resto
	}

	digito2, _ := strconv.Atoi(string(cpf[10]))
	if digito2 != segundoDigito {
		return errors.New("CPF inválido: segundo dígito verificador incorreto")
	}

	return nil
}

// ValidaCNPJ - valida se o CNPJ é válido (formato e dígitos verificadores)
func ValidaCNPJ(cnpj string) error {
	cnpj = strings.ReplaceAll(cnpj, ".", "")
	cnpj = strings.ReplaceAll(cnpj, "/", "")
	cnpj = strings.ReplaceAll(cnpj, "-", "")

	if len(cnpj) != 14 {
		return errors.New("CNPJ deve ter 14 dígitos")
	}

	if strings.Count(cnpj, string(cnpj[0])) == 14 {
		return errors.New("CNPJ inválido: todos os dígitos são iguais")
	}

	soma := 0
	peso := 5
	for i := 0; i < 12; i++ {
		digito, _ := strconv.Atoi(string(cnpj[i]))
		soma += digito * peso
		peso--
		if peso < 2 {
			peso = 9
		}
	}
	resto := soma % 11
	primeiroDigito := 0
	if resto >= 2 {
		primeiroDigito = 11 - resto
	}

	digito1, _ := strconv.Atoi(string(cnpj[12]))
	if digito1 != primeiroDigito {
		return errors.New("CNPJ inválido: primeiro dígito verificador incorreto")
	}

	soma = 0
	peso = 6
	for i := 0; i < 13; i++ {
		digito, _ := strconv.Atoi(string(cnpj[i]))
		soma += digito * peso
		peso--
		if peso < 2 {
			peso = 9
		}
	}
	resto = soma % 11
	segundoDigito := 0
	if resto >= 2 {
		segundoDigito = 11 - resto
	}

	digito2, _ := strconv.Atoi(string(cnpj[13]))
	if digito2 != segundoDigito {
		return errors.New("CNPJ inválido: segundo dígito verificador incorreto")
	}

	return nil
}

// GenerateToken - gera um token aleatório para autenticação
func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// CriptografarSenha - criptografa a senha
func CriptografarSenha(senha string) (string, error) {
	senhaHash, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		return "", WrapError("Erro ao criptografar senha", err, http.StatusInternalServerError)
	}
	return string(senhaHash), nil
}

// VerificaSenha - verifica se a senha fornecida corresponde ao hash armazenado
func VerificaSenha(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
