package types

type KeyType string

var (
	// Null key type.
	Null KeyType = "NULL"

	// PrimaryKey The Primary key
	PrimaryKey KeyType = "PRIMARY_KEY"
)

func (k KeyType) Number() IndexType {
	switch k {
	case Null:
		return 0
	case PrimaryKey:
		return 1
	default:
		return 0
	}
}
