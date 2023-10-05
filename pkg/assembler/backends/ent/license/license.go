// Code generated by ent, DO NOT EDIT.

package license

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the license type in the database.
	Label = "license"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldInline holds the string denoting the inline field in the database.
	FieldInline = "inline"
	// FieldListVersion holds the string denoting the list_version field in the database.
	FieldListVersion = "list_version"
	// EdgeDeclaredInCertifyLegals holds the string denoting the declared_in_certify_legals edge name in mutations.
	EdgeDeclaredInCertifyLegals = "declared_in_certify_legals"
	// EdgeDiscoveredInCertifyLegals holds the string denoting the discovered_in_certify_legals edge name in mutations.
	EdgeDiscoveredInCertifyLegals = "discovered_in_certify_legals"
	// Table holds the table name of the license in the database.
	Table = "licenses"
	// DeclaredInCertifyLegalsTable is the table that holds the declared_in_certify_legals relation/edge. The primary key declared below.
	DeclaredInCertifyLegalsTable = "certify_legal_declared_licenses"
	// DeclaredInCertifyLegalsInverseTable is the table name for the CertifyLegal entity.
	// It exists in this package in order to avoid circular dependency with the "certifylegal" package.
	DeclaredInCertifyLegalsInverseTable = "certify_legals"
	// DiscoveredInCertifyLegalsTable is the table that holds the discovered_in_certify_legals relation/edge. The primary key declared below.
	DiscoveredInCertifyLegalsTable = "certify_legal_discovered_licenses"
	// DiscoveredInCertifyLegalsInverseTable is the table name for the CertifyLegal entity.
	// It exists in this package in order to avoid circular dependency with the "certifylegal" package.
	DiscoveredInCertifyLegalsInverseTable = "certify_legals"
)

// Columns holds all SQL columns for license fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldInline,
	FieldListVersion,
}

var (
	// DeclaredInCertifyLegalsPrimaryKey and DeclaredInCertifyLegalsColumn2 are the table columns denoting the
	// primary key for the declared_in_certify_legals relation (M2M).
	DeclaredInCertifyLegalsPrimaryKey = []string{"certify_legal_id", "license_id"}
	// DiscoveredInCertifyLegalsPrimaryKey and DiscoveredInCertifyLegalsColumn2 are the table columns denoting the
	// primary key for the discovered_in_certify_legals relation (M2M).
	DiscoveredInCertifyLegalsPrimaryKey = []string{"certify_legal_id", "license_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
)

// OrderOption defines the ordering options for the License queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByInline orders the results by the inline field.
func ByInline(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldInline, opts...).ToFunc()
}

// ByListVersion orders the results by the list_version field.
func ByListVersion(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldListVersion, opts...).ToFunc()
}

// ByDeclaredInCertifyLegalsCount orders the results by declared_in_certify_legals count.
func ByDeclaredInCertifyLegalsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newDeclaredInCertifyLegalsStep(), opts...)
	}
}

// ByDeclaredInCertifyLegals orders the results by declared_in_certify_legals terms.
func ByDeclaredInCertifyLegals(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newDeclaredInCertifyLegalsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByDiscoveredInCertifyLegalsCount orders the results by discovered_in_certify_legals count.
func ByDiscoveredInCertifyLegalsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newDiscoveredInCertifyLegalsStep(), opts...)
	}
}

// ByDiscoveredInCertifyLegals orders the results by discovered_in_certify_legals terms.
func ByDiscoveredInCertifyLegals(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newDiscoveredInCertifyLegalsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newDeclaredInCertifyLegalsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(DeclaredInCertifyLegalsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, DeclaredInCertifyLegalsTable, DeclaredInCertifyLegalsPrimaryKey...),
	)
}
func newDiscoveredInCertifyLegalsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(DiscoveredInCertifyLegalsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, DiscoveredInCertifyLegalsTable, DiscoveredInCertifyLegalsPrimaryKey...),
	)
}