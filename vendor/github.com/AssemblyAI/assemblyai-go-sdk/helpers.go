package assemblyai

func String(v string) *string {
	return &v
}

func Int64(v int64) *int64 {
	return &v
}

func Float64(v float64) *float64 {
	return &v
}

func Bool(v bool) *bool {
	return &v
}

func ToString(p *string) (v string) {
	if p == nil {
		return v
	}
	return *p
}

func ToInt64(p *int64) (v int64) {
	if p == nil {
		return v
	}
	return *p
}

func ToFloat64(p *float64) (v float64) {
	if p == nil {
		return v
	}
	return *p
}

func ToBool(p *bool) (v bool) {
	if p == nil {
		return v
	}
	return *p
}
