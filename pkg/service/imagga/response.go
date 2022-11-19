package imagga

// Response
// is for image tagging response.
type Response struct {
	Result struct {
		Tags []struct {
			Confidence float64 `json:"confidence"`
			Tag        struct {
				En string `json:"en"`
			} `json:"tag"`
		} `json:"tags"`
	} `json:"result"`
}
