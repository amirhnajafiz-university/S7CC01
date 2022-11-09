package imagga

type Response struct {
	Result struct {
		Tags []struct {
			Confidence string `json:"confidence"`
			Tag        struct {
				En string `json:"en"`
			} `json:"tag"`
		} `json:"tags"`
	} `json:"result"`
}
