package dto

// CreateOrderRequest представляет запрос на создание нового заказа.
type CreateOrderRequest struct {
	MarkBox    string                   `json:"mark_box" binding:"required,min=1,max=10"`
	CustomerID string                   `json:"customer_id" binding:"required,min=1"`
	Items      []CreateOrderItemRequest `json:"items" binding:"required,min=1,dive"`
	Notes      string                   `json:"notes,omitempty"`
}

// CreateOrderItemRequest представляет элемент заказа в запросе на создание.
type CreateOrderItemRequest struct {
	Variety    string  `json:"variety" binding:"required,min=1,max=100"`
	Length     int     `json:"length" binding:"required,min=1,max=200"`
	BoxCount   float64 `json:"box_count" binding:"required,gt=0"`
	PackRate   int     `json:"pack_rate" binding:"required,min=1"`
	TotalStems int     `json:"total_stems" binding:"required,min=1"`
	FarmName   string  `json:"farm_name" binding:"required,min=1,max=100"`
	TruckName  string  `json:"truck_name" binding:"required,min=1,max=100"`
	Comments   string  `json:"comments,omitempty"`
	Price      float64 `json:"price,omitempty" binding:"gte=0"`
}
