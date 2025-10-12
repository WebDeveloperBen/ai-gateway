package health

// ===================================
// |      Data Transfer Objects      |
// ===================================

type Health struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// ===================================
// |         Request Models          |
// ===================================

type GetRequest struct{}

// ===================================
// |         Response Models         |
// ===================================

type GetResponse struct {
	Body Health `json:"body"`
}
