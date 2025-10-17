package repository

import (
	"github.com/jampa_trip/internal/model"
	"github.com/jampa_trip/internal/query"
	"gorm.io/gorm"
)

// CompanyRepository - objeto de contexto
type CompanyRepository struct {
	DB *gorm.DB
}

// CompanyRepositoryNew - construtor do objeto
func CompanyRepositoryNew(DB *gorm.DB) *CompanyRepository {
	return &CompanyRepository{
		DB: DB,
	}
}

// GetByID - busca uma empresa pelo ID
func (receiver *CompanyRepository) GetByID(id int) (*model.Company, error) {
	row := &model.Company{}

	err := receiver.DB.Raw(query.GetCompanyByID, id).Row().Scan(
		&row.ID,
		&row.Name,
		&row.Email,
		&row.Password,
		&row.CNPJ,
		&row.Phone,
		&row.Address,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return row, nil
}

// GetByEmail - busca uma empresa pelo email
func (receiver *CompanyRepository) GetByEmail(email string) (*model.Company, error) {
	row := &model.Company{}

	err := receiver.DB.Raw(query.GetCompanyByEmail, email).Row().Scan(
		&row.ID,
		&row.Name,
		&row.Email,
		&row.Password,
		&row.CNPJ,
		&row.Phone,
		&row.Address,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return row, nil
}

// Create - cria uma nova empresa
func (receiver *CompanyRepository) Create(company *model.Company) error {
	err := receiver.DB.Raw(query.CreateCompany,
		company.Name,
		company.Email,
		company.Password,
		company.CNPJ,
		company.Phone,
		company.Address,
		company.CreatedAt,
		company.UpdatedAt,
	).Row().Scan(&company.ID)

	return err
}

// Update - atualiza os campos enviados no map
func (receiver *CompanyRepository) Update(id int, updates map[string]interface{}) error {
	result := receiver.DB.Model(&model.Company{}).Where("id = ?", id).Updates(updates)
	return result.Error
}

// List - busca todas as empresas
func (receiver *CompanyRepository) List(filtros *model.Company) ([]*model.Company, error) {
	rows, err := receiver.DB.Raw(query.ListAllCompanies, filtros.Name, filtros.Name, filtros.Email, filtros.Email,
		filtros.CNPJ, filtros.CNPJ, filtros.Phone, filtros.Phone, filtros.Address, filtros.Address).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []*model.Company
	for rows.Next() {
		company := &model.Company{}
		err := rows.Scan(
			&company.ID,
			&company.Name,
			&company.Email,
			&company.CNPJ,
			&company.Phone,
			&company.Address,
			&company.CreatedAt,
			&company.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}

	return companies, nil
}

// EmailExiste - verifica se o email já está cadastrado
func (r *CompanyRepository) EmailExiste(email string) (bool, error) {
	var count int64
	err := r.DB.Model(&model.Company{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// EmailExisteParaOutraEmpresa - verifica se o email já está cadastrado para outra empresa
func (r *CompanyRepository) EmailExisteParaOutraEmpresa(email string, id int) (bool, error) {
	var count int64
	err := r.DB.Model(&model.Company{}).Where("email = ? AND id != ?", email, id).Count(&count).Error
	return count > 0, err
}
