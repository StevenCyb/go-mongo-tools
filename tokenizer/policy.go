package tokenizer

// PolicyType represent policy type values.
type PolicyType byte

const (
	// WhitelistPolicy declare a policy
	// to use whitelist logic.
	WhitelistPolicy PolicyType = 0
	// BlacklistPolicy declare a policy
	// to use blacklist logic.
	BlacklistPolicy PolicyType = 1
)

// Policy handles policy checks
// based on configuration.
type Policy struct {
	values     []string
	policyType PolicyType
}

// Allow check if a value is allowed.
func (p *Policy) Allow(value string) bool {
	for _, policyValue := range p.values {
		if policyValue == value {
			return p.policyType == WhitelistPolicy
		}
	}

	return p.policyType == BlacklistPolicy
}

// NewPolicy create a new policy instance with given arguments.
func NewPolicy(policyType PolicyType, values ...string) *Policy {
	return &Policy{
		policyType: policyType,
		values:     values,
	}
}
