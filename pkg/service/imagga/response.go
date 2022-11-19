package imagga

type TagResponse struct {
	Result struct {
		Tags []struct {
			Confidence uint `json:"confidence"`
			Tag        struct {
				En string `json:"en"`
			} `json:"tag"`
		} `json:"tags"`
	} `json:"result"`
}

type UploadResponse struct {
	Result struct {
		UploadId string `json:"upload_id"`
	} `json:"result"`
	Status struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"status"`
}
