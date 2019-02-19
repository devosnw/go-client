package client

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/splitio/go-toolkit/datastructures/set"

	"github.com/splitio/go-client/splitio/engine/evaluator"

	"github.com/splitio/go-toolkit/logging"
)

// InputValidation struct is responsible for cheking any input of treatment and
// track methods.

// MaxLength constant to check the length of the splits
const MaxLength = 250

// RegExpEventType constant that EventType must match
const RegExpEventType = "^[a-zA-Z0-9][-_.:a-zA-Z0-9]{0,79}$"

type inputValidation struct {
	logger logging.LoggerInterface
}

func parseIfNumeric(value interface{}, operation string) (string, error) {
	f, float := value.(float64)
	i, integer := value.(int)
	i32, integer32 := value.(int32)
	i64, integer64 := value.(int64)

	if float {
		if math.IsNaN(f) || math.IsInf(f, -1) || math.IsInf(f, 1) || math.IsInf(f, 0) {
			return "", errors.New(operation + ": you passed an invalid key, key must be a non-empty string")
		}
		return strconv.FormatFloat(f, 'f', -1, 64), nil
	}
	if integer {
		return strconv.Itoa(i), nil
	}
	if integer32 {
		return strconv.FormatInt(int64(i32), 10), nil
	}
	if integer64 {
		return strconv.FormatInt(i64, 10), nil
	}
	return "", errors.New(operation + ": you passed an invalid key, key must be a non-empty string")
}

func (i *inputValidation) checkWhitespaces(value string, operation string) string {
	trimmed := strings.TrimSpace(value)
	if strings.TrimSpace(value) != value {
		i.logger.Warning(fmt.Sprintf(operation+": split name '%s' has extra whitespace, trimming", value))
	}
	return trimmed
}

func checkIsEmptyString(value string, name string, operation string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New(operation + ": you passed an empty " + name + ", " + name + " must be a non-empty string")
	}
	return nil
}

func checkIsNotValidLength(value string, name string, operation string) error {
	if len(value) > MaxLength {
		return errors.New(operation + ": " + name + " too long - must be " + strconv.Itoa(MaxLength) + " characters or less")
	}
	return nil
}

func checkIsValidString(value string, name string, operation string) error {
	err := checkIsEmptyString(value, name, operation)
	if err != nil {
		return err
	}
	return checkIsNotValidLength(value, name, operation)
}

func checkValidKeyObject(matchingKey string, bucketingKey *string, operation string) (string, *string, error) {
	if bucketingKey == nil {
		return "", nil, errors.New(operation + ": you passed a nil bucketingKey, bucketingKey must be a non-empty string")
	}

	err := checkIsValidString(matchingKey, "matchingKey", operation)
	if err != nil {
		return "", nil, err
	}

	err = checkIsValidString(*bucketingKey, "bucketingKey", operation)
	if err != nil {
		return "", nil, err
	}

	return matchingKey, bucketingKey, nil
}

// ValidateTreatmentKey implements the validation for Treatment call
func (i *inputValidation) ValidateTreatmentKey(key interface{}, operation string) (string, *string, error) {
	if key == nil {
		return "", nil, errors.New(operation + ": you passed a nil key, key must be a non-empty string")
	}
	okey, ok := key.(*Key)
	if ok {
		return checkValidKeyObject(okey.MatchingKey, &okey.BucketingKey, operation)
	}
	var sMatchingKey string
	var err error
	sMatchingKey, ok = key.(string)
	if !ok {
		sMatchingKey, err = parseIfNumeric(key, operation)
		if err != nil {
			return "", nil, err
		}
		i.logger.Warning(fmt.Sprintf(operation+": key %s is not of type string, converting", key))
	}
	err = checkIsValidString(sMatchingKey, "key", "Treatment")
	if err != nil {
		return "", nil, err
	}

	return sMatchingKey, nil, nil
}

// ValidateFeatureName implements the validation for FetureName
func (i *inputValidation) ValidateFeatureName(featureName string) (string, error) {
	err := checkIsEmptyString(featureName, "featureName", "Treatment")
	if err != nil {
		return "", err
	}
	return i.checkWhitespaces(featureName, "Treatment"), nil
}

func checkEventType(eventType string) error {
	err := checkIsEmptyString(eventType, "event type", "Track")
	if err != nil {
		return err
	}
	var r = regexp.MustCompile(RegExpEventType)
	if !r.MatchString(eventType) {
		return errors.New("Track: you passed " + eventType + ", event name must adhere to " +
			"the regular expression " + RegExpEventType + ". This means an event " +
			"name must be alphanumeric, cannot be more than 80 characters long, and can " +
			"only include a dash, underscore, period, or colon as separators of " +
			"alphanumeric characters")
	}
	return nil
}

func (i *inputValidation) checkTrafficType(trafficType string) (string, error) {
	err := checkIsEmptyString(trafficType, "traffic type", "Track")
	if err != nil {
		return "", err
	}
	toLower := strings.ToLower(trafficType)
	if toLower != trafficType {
		i.logger.Warning("Track: traffic type should be all lowercase - converting string to lowercase")
	}
	return toLower, nil
}

func checkValue(value interface{}) error {
	if value == nil {
		return nil
	}

	_, float := value.(float64)
	_, integer := value.(int)
	_, integer32 := value.(int32)
	_, integer64 := value.(int64)

	if float || integer || integer32 || integer64 {
		return nil
	}
	return errors.New("Track: value must be a number")
}

// ValidateTrackInputs implements the validation for Track call
func (i *inputValidation) ValidateTrackInputs(key string, trafficType string, eventType string, value interface{}) (string, string, string, interface{}, error) {
	err := checkIsValidString(key, "key", "Track")
	if err != nil {
		return "", trafficType, eventType, value, err
	}

	err = checkEventType(eventType)
	if err != nil {
		return key, trafficType, "", value, err
	}

	trafficType, err = i.checkTrafficType(trafficType)
	if err != nil {
		return key, "", eventType, value, err
	}

	err = checkValue(value)
	if err != nil {
		return key, trafficType, eventType, nil, err
	}

	return key, trafficType, eventType, value, nil
}

// ValidateManagerInputs implements the validation for Track call
func (i *inputValidation) ValidateManagerInputs(feature string) error {
	return checkIsEmptyString(feature, "split name", "Split")
}

// ValidateFeatureNames implements the validation for Treatments call
func (i *inputValidation) ValidateFeatureNames(features []string) ([]string, error) {
	var featuresSet = set.NewSet()
	if len(features) == 0 {
		return []string{}, errors.New("Treatments: features must be a non-empty array")
	}
	for _, feature := range features {
		f, err := i.ValidateFeatureName(feature)
		if err != nil {
			i.logger.Error(err.Error())
		} else {
			featuresSet.Add(f)
		}
	}
	if featuresSet.IsEmpty() {
		return []string{}, errors.New("Treatments: features must be a non-empty array")
	}
	f := make([]string, featuresSet.Size())
	for i, v := range featuresSet.List() {
		s, ok := v.(string)
		if ok {
			f[i] = s
		}
	}
	return f, nil
}

func (i *inputValidation) GenerateControlTreatments(features []string) map[string]string {
	treatments := make(map[string]string)
	filtered, err := i.ValidateFeatureNames(features)
	if err != nil {
		return treatments
	}
	for _, feature := range filtered {
		treatments[feature] = evaluator.Control
	}
	return treatments
}