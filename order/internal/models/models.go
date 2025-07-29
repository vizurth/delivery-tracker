package models

type OrderCreateRequest struct {
	ID     int    `json:"id"`
	UserId int    `json:"user_id"`
	Status string `json:"status"`
}

type OrderUpdateRequest struct {
	Status string `json:"status"`
}

type OrderGet struct {
	OrderCreateRequest
}

type OrderAcceptRequest struct {
	ID        int    `json:"id"`
	UserId    int    `json:"user_id"`
	CourierId string `json:"courier_id"`
	Status    string `json:"status"`
}
