package keyvalue

type EntryRequest struct {
	Name  string `json:"name" binding:"required,alpha,min=1,max=100"`
	Value string `json:"value" binding:"required,alpha,min=1,max=255"`
}

type EntryResponse struct {
	Name string `uri:"name" binding:"required,alpha,min=1,max=255"`
}
