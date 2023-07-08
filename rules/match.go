package rules

type Banner struct {
	Protocol string
	Port     string
	Header   string
	Body     string
	Response string
	Cert     string
	Title    string
	Hash     string
	Icon     string
	ICP      string
}

func (banner Banner) Search() []string {
	var products []string
	for _, fingerPrint := range FingerPrints {
		if isMatch := banner.Match(fingerPrint); isMatch {
			products = append(products, fingerPrint.ProductName)
		}
	}

	return removeDuplicate(products)
}

func removeDuplicate(p []string) []string {
	result := make([]string, 0, len(p))
	temp := map[string]struct{}{}
	for _, item := range p {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func (banner Banner) Match(fingerprint *FingerPrint) bool {
	//TODO 根据Expression去匹配Banner

	return true
}
