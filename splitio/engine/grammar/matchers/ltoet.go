package matchers

import (
	"github.com/splitio/go-client/splitio/engine/grammar/matchers/datatypes"
)

// LessThanOrEqualToMatcher will match if two numbers or two datetimes are equal
type LessThanOrEqualToMatcher struct {
	Matcher
	ComparisonDataType string
	ComparisonValue    interface{}
}

// Match will match if the comparisonValue is less than or equal to the matchingValue
func (m *LessThanOrEqualToMatcher) Match(key string, attributes map[string]interface{}, bucketingKey *string) bool {

	matchingRaw, err := m.matchingKey(key, attributes)
	if err != nil {
		return false
	}

	matchingValue, okMatching := matchingRaw.(int64)
	comparisonValue, okComparison := m.ComparisonValue.(int64)
	if !okMatching || !okComparison {
		return false
	}

	switch m.ComparisonDataType {
	case datatypes.Number:
	case datatypes.Datetime:
		matchingValue = datatypes.ZeroTimeTS(matchingValue)
		comparisonValue = datatypes.ZeroTimeTS(comparisonValue)
	default:
		return false
	}

	return matchingValue <= comparisonValue
}

// NewLessThanOrEqualToMatcher returns a pointer to a new instance of LessThanOrEqualToMatcher
func NewLessThanOrEqualToMatcher(negate bool, cmpVal int64, cmpType string, attributeName *string) *LessThanOrEqualToMatcher {
	return &LessThanOrEqualToMatcher{
		Matcher: Matcher{
			negate:        negate,
			attributeName: attributeName,
		},
		ComparisonValue:    cmpVal,
		ComparisonDataType: cmpType,
	}
}
