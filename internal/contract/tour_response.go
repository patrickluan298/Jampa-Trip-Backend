package contract

// TourResponse - resposta com dados do passeio
type TourResponse struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Dates         []string `json:"dates"`
	DepartureTime string   `json:"departure_time"`
	ArrivalTime   string   `json:"arrival_time"`
	MaxPeople     int      `json:"max_people"`
	Description   string   `json:"description"`
	Images        []string `json:"images"`
	Price         float64  `json:"price"`
	CompanyID     int      `json:"company_id"`
	CompanyName   string   `json:"company_name"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
}

// CreateTourResponse - resposta de criação de passeio
type CreateTourResponse struct {
	Success bool         `json:"success"`
	Tour    TourResponse `json:"tour"`
}

// UpdateTourResponse - resposta de atualização de passeio
type UpdateTourResponse struct {
	Success bool         `json:"success"`
	Tour    TourResponse `json:"tour"`
}

// ListToursResponse - resposta de listagem de passeios
type ListToursResponse struct {
	Success    bool               `json:"success"`
	Tours      []TourResponse     `json:"tours"`
	Pagination PaginationResponse `json:"pagination"`
}

// GetMyToursResponse - resposta de listagem de passeios da empresa
type GetMyToursResponse struct {
	Success    bool               `json:"success"`
	Tours      []MyTourResponse   `json:"tours"`
	Pagination PaginationResponse `json:"pagination"`
}

// MyTourResponse - resposta com dados do passeio da empresa (inclui contagem de reservas)
type MyTourResponse struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	Dates             []string `json:"dates"`
	DepartureTime     string   `json:"departure_time"`
	ArrivalTime       string   `json:"arrival_time"`
	MaxPeople         int      `json:"max_people"`
	Description       string   `json:"description"`
	Images            []string `json:"images"`
	Price             float64  `json:"price"`
	CreatedAt         string   `json:"created_at"`
	ReservationsCount int      `json:"reservations_count"`
}

// DeleteTourResponse - resposta de exclusão de passeio
type DeleteTourResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// PaginationResponse - resposta de paginação
type PaginationResponse struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	TotalItems   int `json:"total_items"`
	ItemsPerPage int `json:"items_per_page"`
}
