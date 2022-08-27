package enums

type Role string

const (
	ROLE_ADMIN       = Role("ROLE_ADMIN")
	ROLE_CUSTOMER    = Role("ROLE_CUSTOMER")
	ROLE_SUPER_ADMIN = Role("ROLE_SUPER_ADMIN")
)
