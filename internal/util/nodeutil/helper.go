package nodeutil

import "strings"

type Type struct {
	TypeName       string
	PureType       string
	IsArray        bool
	IsPointer      bool
	IsArrayPointer bool
	IsInterface    bool
}

func GetType(tp string) Type {
	pureType := tp
	isArray := false
	isPointer := false
	isArrayPointer := false

	if len(tp) > 0 && tp[0] == '*' {
		isPointer = true
		pureType = tp[1:]
	}

	if len(pureType) > 0 && pureType[0] == '[' && pureType[1] == ']' {
		isArray = true
		pureType = pureType[2:]
	}

	if isArray && len(pureType) > 0 && pureType[0] == '*' {
		isArrayPointer = true
		pureType = pureType[1:]
	}

	typeName := pureType
	if dotIndex := len(pureType) - 1; dotIndex >= 0 && strings.Contains(pureType, ".") {
		splitType := strings.Split(pureType, ".")
		typeName = splitType[len(splitType)-1]
	}

	return Type{
		TypeName:       typeName,
		PureType:       pureType,
		IsArray:        isArray,
		IsPointer:      isPointer,
		IsArrayPointer: isArrayPointer,
	}
}

func GetShortIdentityByFullPath(tp string) (short string) {
	parts := strings.Split(tp, "/")
	if len(parts) > 0 {
		short = parts[len(parts)-1]
	}
	return short
}
