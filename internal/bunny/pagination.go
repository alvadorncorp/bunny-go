package bunny

const (
	DefaultPage    = 0
	DefaultPerPage = 10
)

type OptionalValue[T any] struct {
	param *T
}

type PageParams struct {
	Page    OptionalValue[int]
	PerPage OptionalValue[int]
}

type Page[T any] struct {
	Items        []T  `json:"Items"`
	CurrentPage  int  `json:"CurrentPage"`
	TotalItems   int  `json:"TotalItems"`
	HasMoreItems bool `json:"HasMoreItems"`
}

func Optional[T any](param *T) OptionalValue[T] {
	return OptionalValue[T]{param}
}

func (o OptionalValue[T]) ValueOrDefault(defaultValue T) T {
	if o.param != nil {
		return *o.param
	}

	return defaultValue
}
