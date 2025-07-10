package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	VALIDATOR_TAG = "valid"
	SEP           = ","
	ARG_SEP       = "="

	// 内置校验器标签
	MAX      = "max"
	MIN      = "min"
	REQUIRED = "required"
	EMAIL    = "email"

	EMAIL_FROMAT = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

type validatorError struct {
	Field string
	Tag   string
	Value any
	Msg   string
}

func (ve *validatorError) Error() string {
	return ve.Msg
}

type ValidatorErrors []validatorError

func (ves ValidatorErrors) Error() string {
	msgs := []string{}
	for _, ve := range ves {
		msgs = append(msgs, ve.Error())
	}
	return strings.Join(msgs, ",")
}

type validatorFunc func(value any, param string) bool

type Validator struct {
	validators map[string]validatorFunc
}

// 注册内置校验器
func (va *Validator) registerBuiltinValidators() {
	// 检测字段非空
	va.validators[REQUIRED] = func(value any, param string) bool {
		return !isEmpty(value)
	}

	// 检测字段长度下限
	va.validators[MIN] = func(value any, param string) bool {
		minValue, _ := strconv.Atoi(param)
		switch t := value.(type) {
		case string:
			return len(t) >= minValue
		case int:
			return t >= minValue
		}
		return false
	}
	// 检测字段长度上限
	va.validators[MAX] = func(value any, param string) bool {
		maxValue, _ := strconv.Atoi(param)
		switch t := value.(type) {
		case string:
			return len(t) <= maxValue
		case int:
			return t <= maxValue
		}
		return false
	}

	// 检测邮箱是否符合格式
	va.validators[EMAIL] = func(value any, param string) bool {
		email, ok := value.(string)
		if !ok {
			return false
		}
		emailRegex := regexp.MustCompile(EMAIL_FROMAT)
		return emailRegex.MatchString(email)
	}
}

func isEmpty(value any) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Array, reflect.Slice, reflect.Map:
		return v.Len() == 0
	case reflect.Ptr:
		return v.IsNil()
	default:
		return false
	}
}

type validatorRule struct {
	name  string
	param string
}

func (va *Validator) parseValidatorTag(tag string) []validatorRule {
	var rules []validatorRule
	ruleItems := strings.Split(tag, SEP)
	for _, item := range ruleItems {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if strings.Contains(item, "=") {
			part := strings.SplitN(item, ARG_SEP, 2)
			rules = append(rules, validatorRule{
				name:  part[0],
				param: part[1],
			})
		} else {
			rules = append(rules, validatorRule{
				name:  item,
				param: "",
			})
		}
	}
	return rules
}

func (va *Validator) checkField(value any, fieldName string, tag string) []validatorError {
	var errs []validatorError
	rules := va.parseValidatorTag(tag)
	for _, rule := range rules {
		validFunc, found := va.validators[rule.name]
		if !found {
			continue
		}
		if !validFunc(value, rule.param) {
			errs = append(errs, validatorError{
				Field: fieldName,
				Value: value,
				Tag:   tag,
				Msg:   fmt.Sprintf("|Field: %s, Rule: %s, Param: %s|", fieldName, rule.name, rule.param),
			})
		}
	}
	return errs
}

func NewValidator() *Validator {
	va := &Validator{
		validators: make(map[string]validatorFunc),
	}
	va.registerBuiltinValidators()
	return va
}

func (va *Validator) Check(s any) error {
	sValue := reflect.ValueOf(s)
	if sValue.Kind() == reflect.Ptr {
		sValue = sValue.Elem()
	}
	if sValue.Kind() != reflect.Struct {
		return fmt.Errorf("validator! [type] -expected: struct -fact: %s", sValue.Kind().String())
	}
	var errs ValidatorErrors
	sType := sValue.Type()
	for i := 0; i < sValue.NumField(); i++ {
		svField := sValue.Field(i)
		stField := sType.Field(i)

		tag := stField.Tag.Get(VALIDATOR_TAG)
		if tag == "" {
			continue
		}
		fieldErrs := va.checkField(svField.Interface(), stField.Name, tag)
		if fieldErrs != nil {
			errs = append(errs, fieldErrs...)
		}
	}
	if len(errs) != 0 {
		return errs
	}
	return nil
}
