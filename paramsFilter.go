package slf

// ParamsFilter describes function, used to filter params list
type ParamsFilter func([]Param) []Param

// denyAllP does not allow any param
func denyAllP([]Param) []Param { return []Param{} }

// allowAllP allows all packets
func allowAllP(p []Param) []Param { return p }

// wlParams contains whitelist filtering
type wlParams struct {
	allowed map[string]bool
}

func (w wlParams) filter(p []Param) []Param {
	if len(p) == 0 {
		return []Param{}
	}

	n := []Param{}
	for _, param := range p {
		if _, ok := w.allowed[param.GetKey()]; ok {
			n = append(n, param)
		}
	}

	return n
}

// wlParams contains blacklist filtering
type blParams struct {
	allowed map[string]bool
}

func (w blParams) filter(p []Param) []Param {
	if len(p) == 0 {
		return []Param{}
	}

	n := []Param{}
	for _, param := range p {
		if _, ok := w.allowed[param.GetKey()]; !ok {
			n = append(n, param)
		}
	}

	return n
}

// NewWhiteListParamsFilter return params filter, that allows only params
// with names, passed to this function
func NewWhiteListParamsFilter(allowedNames []string) ParamsFilter {
	if len(allowedNames) == 0 {
		return denyAllP
	}

	f := wlParams{allowed: make(map[string]bool, len(allowedNames))}
	for _, name := range allowedNames {
		f.allowed[name] = true
	}

	return f.filter
}

// NewBlackListParamsFilter returns params filter, that clears params
// with name, passed to this function
func NewBlackListParamsFilter(allowedNames []string) ParamsFilter {
	if len(allowedNames) == 0 {
		return allowAllP
	}

	f := blParams{allowed: make(map[string]bool, len(allowedNames))}
	for _, name := range allowedNames {
		f.allowed[name] = true
	}

	return f.filter
}
