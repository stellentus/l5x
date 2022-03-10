package l5x

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/stellentus/go-plc"
)

type Type interface {
	// PlcName is the name of the type, as used by the PLC (e.g. "my_struct" or "INT")
	PlcName() string

	// GoName is the name of the type, as used by go (e.g. "MyStruct" or "int16")
	GoName() string

	// GoTypeString is the go definition of the type (e.g. "struct {Val uint8}" or "").
	// For basic types, it's an empty string.
	GoTypeString() string
}

type NamedType struct {
	PlcName string
	GoName  string
	Type    Type
}

type TypeList []Type

// NamedTypeDeclaration is used to declare a named type inside a struct.
// It automatically embeds if the variable name matches the type name.
func NamedTypeDeclaration(nt NamedType) string {
	typeName := nt.Type.GoName()
	if nt.GoName == typeName {
		return nt.GoName
	}
	return nt.GoName + " " + typeName
}

func TypeDefinition(ty Type) string {
	if ty.GoTypeString() == "" {
		// This type shouldn't be defined
		return ""
	}
	return "type " + makeValidIdentifier(ty.GoName()) + " " + ty.GoTypeString()
}

var (
	typeBOOL   = basicType{"BOOL", "bool"}
	typeSINT   = basicType{"SINT", "int8"}
	typeINT    = basicType{"INT", "int16"}
	typeDINT   = basicType{"DINT", "int32"}
	typeLINT   = basicType{"LINT", "int64"}
	typeUSINT  = basicType{"USINT", "uint8"}
	typeUINT   = basicType{"UINT", "uint16"}
	typeUDINT  = basicType{"UDINT", "uint32"}
	typeULINT  = basicType{"ULINT", "uint64"}
	typeREAL   = basicType{"REAL", "float32"}
	typeLREAL  = basicType{"LREAL", "float64"}
	typeSTRING = basicType{"STRING", "string"}
	typeBYTE   = basicType{"BYTE", "byte"}    // might not be correct; might be an array
	typeWORD   = basicType{"WORD", "uint16"}  // might not be correct; might be an array
	typeDWORD  = basicType{"DWORD", "uint32"} // might not be correct; might be an array
	typeLWORD  = basicType{"LWORD", "uint64"} // might not be correct; might be an array
)

func NewTypeList() TypeList {
	return TypeList{
		typeBOOL,
		typeSINT,
		typeINT,
		typeDINT,
		typeLINT,
		typeUSINT,
		typeUINT,
		typeUDINT,
		typeULINT,
		typeREAL,
		typeLREAL,
		typeSTRING,
		typeBYTE,
		typeWORD,
		typeDWORD,
		typeLWORD,
	}
}

func (tl TypeList) WithPlcName(name string) (Type, error) {
	for _, ty := range tl {
		if ty.PlcName() == name {
			return ty, nil
		}
	}
	return nil, errors.New("DataType '" + name + "' couldn't be found")
}

func (tl TypeList) WriteDefinitions(wr io.Writer) error {
	for _, ty := range tl {
		str := TypeDefinition(ty)
		if str == "" {
			continue
		}
		_, err := wr.Write([]byte(str + "\n\n"))
		if err != nil {
			return err
		}
	}
	return nil
}

// AddControlLogixTypes adds a couple of types that are (I think) built into ControlLogix
func (tl *TypeList) AddControlLogixTypes() error {
	timer, err := DataType{
		Name: "TIMER", // Reference: https://twcontrols.com/lessons/rslogix-500-training-timers-ton-tof-and-rto?rq=RsLogix%20500%20Training%20-%20timers
		Members: []Member{
			{Name: "PRE", DataType: "DINT"}, // "Preset": the target value.
			{Name: "ACC", DataType: "DINT"}, // "Accumulated": the current value.
			{Name: "EN", DataType: "BOOL"},  // "Enable Bit": true if the timer is running.
			{Name: "TT", DataType: "BOOL"},  // "Timer Timing Bit": true if (EN && ACC < PRE).
			{Name: "DN", DataType: "BOOL"},  // "Done Bit": true if (EN && ACC > PRE).
		},
	}.AsType(*tl)
	if err != nil {
		return fmt.Errorf("could not register ControlLogix TIMER because %w", err)
	}
	*tl = append(*tl, timer)

	counter, err := DataType{
		Name: "COUNTER", // Reference: https://twcontrols.com/lessons/rslogix-500-training-counters-ctu-and-ctd
		Members: []Member{
			{Name: "PRE", DataType: "DINT"}, // "Preset": the target value.
			{Name: "ACC", DataType: "DINT"}, // "Accumulated": the current value.
			{Name: "CU", DataType: "BOOL"},  // "Count Up Enable Bit": true if the counter is counting up.
			{Name: "CD", DataType: "BOOL"},  // "Count Down Enable Bit": true if the timer is counting down.
			{Name: "DN", DataType: "BOOL"},  // "Done": true if ACC >= PRE (regardless of Count Up or Count Down).
			{Name: "OV", DataType: "BOOL"},  // "Cout Up Overflow Bit": the timer rolled over from +32767 to -32768.
			{Name: "UN", DataType: "BOOL"},  // "Dount Down Underflow Bit": the timer rolled over from -32768 to +32767.
		},
	}.AsType(*tl)
	if err != nil {
		return fmt.Errorf("could not register ControlLogix COUNTER because %w", err)
	}
	*tl = append(*tl, counter)

	pidEnhanced, err := pid_enhancedDT.AsType(*tl)
	if err != nil {
		return fmt.Errorf("could not register ControlLogix PID_ENHANCED because %w", err)
	}
	*tl = append(*tl, pidEnhanced)

	// The MESSAGE type has some data, but this will add it as `struct{}` since we don't use it
	*tl = append(*tl, structType{safeGoName: safeGoName{"MESSAGE", "MESSAGE"}})

	return nil
}

func (mb Member) AsNamedType(knownTypes TypeList) (NamedType, error) {
	if mb.DataType == "BOOL" && mb.Dimension > 1 {
		// When the BOOL type is stored in arrays, it appears to use 32-bit storage.
		// Maybe byte-sized access would also work. 🤷
		if mb.Dimension%32 != 0 {
			// I've only seen examples with multiples of 32.
			return NamedType{}, fmt.Errorf("%w is a BOOL array, but not a multiple of 32 (%d)", ErrUnknownType, mb.Dimension)
		}
		return knownTypes.newNamedType(mb.Name, "UDINT", []int{mb.Dimension / 32})
	}

	return knownTypes.newNamedType(mb.Name, mb.DataType, []int{mb.Dimension})
}

func (par Parameter) AsNamedType(knownTypes TypeList) (NamedType, error) {
	return knownTypes.newNamedType(par.Name, par.DataType, nil)
}

func newNamedType(name string, ty Type) (NamedType, error) {
	nt := NamedType{
		GoName: name,
		Type:   ty,
	}
	if !isPublicGoIdentifier(name) {
		valid := makeValidIdentifier(nt.GoName)
		if valid == "" {
			return NamedType{}, fmt.Errorf("couldn't create valid identifier for '%s'", name)
		}
		nt.PlcName = nt.GoName
		nt.GoName = valid
	}
	return nt, nil
}

// SetAsProgram changes the PlcName to indicate this is a program tag.
func (nt *NamedType) SetAsProgram() {
	if nt.PlcName == "" {
		nt.PlcName = nt.GoName
	}
	nt.PlcName = "Program:" + nt.PlcName
}

func (tl TypeList) newNamedType(name, dataType string, dims []int) (NamedType, error) {
	var nt NamedType
	for _, ty := range tl {
		if ty.PlcName() == dataType {
			var err error
			nt, err = newNamedType(name, ty)
			if err != nil {
				return NamedType{}, err
			}
			break
		}
	}
	if nt.Type == nil {
		return NamedType{}, ErrUnknownType
	}

	if len(dims) == 0 || len(dims) == 1 && dims[0] <= 1 { // not an array
		return nt, nil
	}

	for idx := len(dims) - 1; idx >= 0; idx-- {
		if dims[idx] <= 1 {
			return NamedType{}, fmt.Errorf("couldn't create NamedType with dimensions %v", dims)
		}
		nt.Type = arrayType{
			elementInfo: nt.Type,
			count:       dims[idx],
		}
	}

	return nt, nil
}

func (dt DataType) AsType(knownTypes TypeList) (Type, error) {
	if dt.Class != ClassUser {
		return nil, fmt.Errorf("Unknown class type")
	}

	switch dt.Family {
	case DataTypeFamilyString:
		return parseString(dt.Name, dt.Members)

	case DataTypeFamilyNone:
		return parseStruct(dt.Name, dt.Members, knownTypes)

	default:
		return nil, fmt.Errorf("Unknown data family type")
	}
}

func (aoi AddOnInstrDef) AsType(knownTypes TypeList) (Type, error) {
	sti, err := newStructType(aoi.Name, nil)
	if err != nil {
		return nil, err
	}
	for _, param := range aoi.Parameters {
		if param.DataType == "BIT" {
			// Not yet implemented, so ignore it
			continue
		}
		nm, err := param.AsNamedType(knownTypes)
		if err != nil {
			if errors.Is(err, ErrUnknownType) {
				err = errUnknownTypeSpecific{aoi.Name, param.DataType}
			}
			return nil, err
		}
		sti.members = append(sti.members, nm)
	}
	if len(sti.members) == 0 {
		return nil, fmt.Errorf("Add-on Instruction '%s' produced no parameters (%d)", aoi.Name, len(aoi.Parameters))
	}
	return sti, nil
}

func parseStruct(name string, membs []Member, knownTypes TypeList) (Type, error) {
	sti, err := newStructType(name, nil)
	if err != nil {
		return nil, err
	}
	for _, memb := range membs {
		if memb.DataType == "BIT" {
			// Not yet implemented, so ignore it
			continue
		}
		nm, err := memb.AsNamedType(knownTypes)
		if err != nil {
			if errors.Is(err, ErrUnknownType) {
				err = errUnknownTypeSpecific{name, memb.DataType}
			}
			return nil, err
		}
		sti.members = append(sti.members, nm)
	}
	if len(sti.members) == 0 {
		return nil, fmt.Errorf("DataType '%s' produced no members (%d)", name, len(membs))
	}

	return sti, nil
}

func parseString(name string, memb []Member) (Type, error) {
	if len(memb) != 2 {
		return nil, fmt.Errorf("StringFamily '%s' had %d members instead of 2", name, len(memb))
	}

	if memb[0].Name != "LEN" {
		return nil, fmt.Errorf("StringFamily '%s' LEN member is missing: %s", name, memb[0].Name)
	}
	if memb[0].DataType != "DINT" {
		return nil, fmt.Errorf("StringFamily '%s' LEN.DataType is incorrect: %v", name, memb[0].DataType)
	}
	if memb[0].Dimension != 0 {
		return nil, fmt.Errorf("StringFamily '%s' LEN.Dimension is incorrect: %d", name, memb[0].Dimension)
	}
	if memb[0].Radix != RadixDecimal {
		return nil, fmt.Errorf("StringFamily '%s' LEN.Radix is incorrect: %v", name, memb[0].Radix)
	}

	if memb[1].Name != "DATA" {
		return nil, fmt.Errorf("StringFamily '%s' DATA member is missing: %v", name, memb[1].Name)
	}
	if memb[1].DataType != "SINT" {
		return nil, fmt.Errorf("StringFamily '%s' DATA.DataType is incorrect: %v", name, memb[1].DataType)
	}
	if memb[1].Dimension <= 0 {
		return nil, fmt.Errorf("StringFamily '%s' DATA.Dimension is invalid: %d", name, memb[1].Dimension)
	}
	if memb[1].Radix != RadixASCII {
		return nil, fmt.Errorf("StringFamily '%s' DATA.Radix is incorrect: %v", name, memb[1].Radix)
	}

	return newStringType(name, memb[1].Dimension)
}

type basicType struct {
	plcName, goName string
}

func (bt basicType) PlcName() string      { return bt.plcName }
func (bt basicType) GoName() string       { return bt.goName }
func (bt basicType) GoTypeString() string { return "" }

type arrayType struct {
	elementInfo Type
	count       int
}

func (ati arrayType) PlcName() string { return "" }
func (ati arrayType) GoName() string  { return ati.GoTypeString() }
func (ati arrayType) GoTypeString() string {
	return fmt.Sprintf("[%d]%v", ati.count, ati.elementInfo.GoName())
}

type structType struct {
	safeGoName
	members []NamedType
}

func (sti structType) GoTypeString() string {
	strs := make([]string, len(sti.members))
	for i, member := range sti.members {
		tagSuffix := ""
		if member.PlcName != "" {
			tagSuffix = fmt.Sprintf(" `%s:\"%s\"`", plc.TagPrefix, member.PlcName)
		}
		strs[i] = fmt.Sprintf("\n\t%s%s", NamedTypeDeclaration(member), tagSuffix)
	}
	return fmt.Sprintf("struct {%s\n}", strings.Join(strs, ""))
}
func newStructType(name string, members []NamedType) (structType, error) {
	if members == nil {
		members = []NamedType{}
	}
	sgn, err := newSafeGoName(name)
	if err != nil {
		return structType{}, err
	}
	return structType{
		safeGoName: sgn,
		members:    members,
	}, nil
}

type stringType struct {
	safeGoName
	count int
}

func (sty stringType) GoTypeString() string {
	return fmt.Sprintf("struct {Len int16; Data [%d]int8}", sty.count)
}
func newStringType(name string, count int) (stringType, error) {
	sgn, err := newSafeGoName(name)
	if err != nil {
		return stringType{}, err
	}
	return stringType{
		safeGoName: sgn,
		count:      count,
	}, nil
}

type safeGoName struct {
	goN, plcN string
}

func (sgn safeGoName) PlcName() string { return sgn.plcN }
func (sgn safeGoName) GoName() string  { return sgn.goN }
func newSafeGoName(name string) (safeGoName, error) {
	sgn := safeGoName{goN: name, plcN: name}
	if !isPublicGoIdentifier(name) {
		sgn.goN = makeValidIdentifier(sgn.goN)
		if sgn.goN == "" {
			return safeGoName{}, fmt.Errorf("couldn't create safe go name for '%s'", name)
		}
	}
	return sgn, nil
}

// isValidGoIdentifier determines if str is a valid go identifier. According to the spec:
//     identifier = letter { letter | unicode_digit } .
//     letter        = unicode_letter | "_" .
//     unicode_letter = /* a Unicode code point classified as "Letter" */ .
//     unicode_digit  = /* a Unicode code point classified as "Number, decimal digit" */ .
func isValidGoIdentifier(str string) bool {
	if len(str) <= 0 {
		return false
	}
	isLetter := func(r rune) bool {
		if unicode.IsLetter(r) {
			return true
		}
		return r == '_'
	}
	for pos, r := range str {
		if isLetter(r) {
			continue // A letter is always allowed
		}
		if unicode.IsDigit(r) && pos > 0 {
			continue // After the first position, a digit is allowed
		}
		return false
	}
	return true
}

// isPublicGoIdentifier checks if str is public (i.e. starts with a capital)
// and a valid go identifier.
func isPublicGoIdentifier(str string) bool {
	if !isValidGoIdentifier(str) {
		return false
	}
	for _, r := range str {
		return unicode.IsUpper(r)
	}
	return false
}

// makeValidIdentifier returns a valid go identifer from 'str'.
// It returns "" if it can't make an identifier.
func makeValidIdentifier(str string) string {
	var valid strings.Builder
	for pos, r := range strings.TrimSpace(str) {
		if pos == 0 {
			valid.WriteRune(unicode.ToUpper(r))
			continue
		}
		if !isValidIdentiferRune(r) {
			valid.WriteRune('_')
			continue
		}
		valid.WriteRune(r)
	}
	if !isValidGoIdentifier(valid.String()) {
		return ""
	}
	return valid.String()
}

// isValidIdentiferRune checks if runes after the first one are valid.
func isValidIdentiferRune(r rune) bool {
	if unicode.IsLetter(r) {
		return true
	}
	if unicode.IsDigit(r) {
		return true
	}
	return r == '_'
}

var ErrUnknownType = errors.New("unknown type")

type errUnknownTypeSpecific struct {
	typ, requires string
}

func (err errUnknownTypeSpecific) Error() string {
	return fmt.Sprintf("Type '%s' requires type '%s'", err.typ, err.requires)
}

func (err errUnknownTypeSpecific) Unwrap() error {
	return ErrUnknownType
}
