package by

var (
	// DeclarationOrder the default sort order.
	// The flags will be printed in the same order as they have been defined.
	DeclarationOrder Comparer
	// LongNameAscending sort by long name in ascending order.
	LongNameAscending = StringComparer{Ascending: true, Field: LongName}
	// LongNameDescending sort by long name in descending order.
	LongNameDescending = StringComparer{Field: LongName}
	// ShortNameAscending sort by short name in ascending order.
	ShortNameAscending = StringComparer{Ascending: true, Field: ShortName}
	// ShortNameDescending sort by short name in descending order.
	ShortNameDescending = StringComparer{Field: ShortName}
	// KeyAscending sort by key in ascending order.
	KeyAscending = StringComparer{Ascending: true, Field: Key}
	// KeyDescending sort by key in descending order.
	KeyDescending = StringComparer{Field: Key}
	// UsageAscending sort by usage in ascending order.
	UsageAscending = StringComparer{Ascending: true, Field: Usage}
	// UsageDescending sort by usage in descending order.
	UsageDescending = StringComparer{Field: Usage}
	// RequiredFirst put the required flags first.
	RequiredFirst = BooleanComparer{Ascending: true, Field: IsRequired}
	// RequiredLast put the required flags at the end.
	RequiredLast = BooleanComparer{Ascending: false, Field: IsRequired}
	// DeprecatedFirst put the deprecated flags first.
	DeprecatedFirst = BooleanComparer{Ascending: true, Field: IsDeprecated}
	// DeprecatedLast put the deprecated flags at the end.
	DeprecatedLast = BooleanComparer{Ascending: false, Field: IsDeprecated}
)
