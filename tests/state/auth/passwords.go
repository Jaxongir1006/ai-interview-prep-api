package auth

// Pre-computed bcrypt hashes with cost=4 for fast test execution.
// These passwords can be used in tests without incurring hashing overhead.

const (
	TestPassword1 = "TestPassword_1"
	TestPassword2 = "TestPassword_2"
	TestPassword3 = "TestPassword_3"
	TestPassword4 = "TestPassword_4"
	TestPassword5 = "TestPassword_5"
)

// GetPrecomputedHash returns the pre-computed bcrypt hash for a known test password.
// Returns the hash and true if the password is known, empty string and false otherwise.
func GetPrecomputedHash(password string) (string, bool) {
	var precomputedHashes = map[string]string{
		TestPassword1: "$2a$04$AEA..u1LDK.VVXYVWX0H6uoJjVpQwUJdRsyqd3GdIfLtZEIF5qBh6",
		TestPassword2: "$2a$04$WiCfQfxe8G6u1xaooj701upIvuC9XZZuL1J2ROkCZi9I9HZr7Hxem",
		TestPassword3: "$2a$04$Gee5sJReXDvi6cRHuEh2iu.8wMz0sMh5wM3S30FCwZ4ydOlaWSC6K",
		TestPassword4: "$2a$04$HZEyGrbvQadSgqn1n00YQuu9jhqIxb8s2G46/JcbnJ4Ce3rs53VGq",
		TestPassword5: "$2a$04$3DVtTM1TEMiB.d07POx3cOowMS1SU76m4j3tWxNEpxS/1CCMRlrQS",
	}

	h, ok := precomputedHashes[password]
	return h, ok
}
