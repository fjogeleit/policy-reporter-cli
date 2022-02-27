package cli

type Grouping = string

const (
	CategoryGrouping Grouping = "category"
	ResourceGrouping Grouping = "resource"
	PolicyGrouping   Grouping = "policy"
	NoneGroup        Grouping = "none"
)
