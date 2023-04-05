package keyvalue

type PostEntry struct {
	Name  string `json:"name" binding:"required,alpha,min=1,max=100"`
	Value string `json:"value" binding:"required,alpha,min=1,max=255"`
}

type GetEntry struct {
	Name string `json:"name" uri:"name" binding:"required,alpha,min=1,max=255"`
}

type GetEntryResponse struct {
	Value string `json:"value" uri:"value"`
}
