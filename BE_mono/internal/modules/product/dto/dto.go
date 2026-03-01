package dto

type AdminUpsertRequest struct {
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Stock       *int   `json:"stock"`
	Status      string `json:"status"`
	ImageURL    string `json:"image_url"`
}
