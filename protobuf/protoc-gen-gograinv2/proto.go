package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"unicode"
	"unicode/utf8"

	google_protobuf "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

// code lifted from gogo proto
var isGoKeyword = map[string]bool{
	"break":       true,
	"case":        true,
	"chan":        true,
	"const":       true,
	"continue":    true,
	"default":     true,
	"else":        true,
	"defer":       true,
	"fallthrough": true,
	"for":         true,
	"func":        true,
	"go":          true,
	"goto":        true,
	"if":          true,
	"import":      true,
	"interface":   true,
	"map":         true,
	"package":     true,
	"range":       true,
	"return":      true,
	"select":      true,
	"struct":      true,
	"switch":      true,
	"type":        true,
	"var":         true,
}

// ProtoFile reprpesents a parsed proto file
type ProtoFile struct {
	PackageName string
	Namespace   string
	Messages    []*ProtoMessage
	Services    []*ProtoService
}

// ProtoMessage represents a parsed message in a proto file
type ProtoMessage struct {
	Name       string
	PascalName string
}

// ProtoService represents a parsed service in a proto file
type ProtoService struct {
	Name       string
	PascalName string
	Methods    []*ProtoMethod
}

// ProtoMethod represents a parsed method in a proto service
type ProtoMethod struct {
	Index        int
	Name         string
	PascalName   string
	Input        *ProtoMessage
	Output       *ProtoMessage
	InputStream  bool
	OutputStream bool
}

// ProtoAst transforms a FileDescriptor to an AST that can be used for code generation
func ProtoAst(file *google_protobuf.FileDescriptorProto) *ProtoFile {

	pkg := &ProtoFile{}
	pkg.Namespace = file.GetOptions().GetCsharpNamespace()

	// let us check the option go_package is defined in the file and use that one instead of the
	// default one
	var packageName string
	if file.GetOptions().GetGoPackage() != "" {
		packageName = cleanPackageName(file.GetOptions().GetGoPackage())
	} else {
		packageName = cleanPackageName(file.GetPackage())
	}

	// let us the go package name
	pkg.PackageName = packageName

	// build proto message dictionary
	messages := make(map[string]*ProtoMessage)
	for _, message := range file.GetMessageType() {
		m := &ProtoMessage{
			Name:       message.GetName(),
			PascalName: MakeFirstLowerCase(message.GetName()),
		}
		pkg.Messages = append(pkg.Messages, m)
		messages[m.Name] = m
		log.Printf("name=%s m=%v", m.Name, m)
	}

	mustGetMessage := func(m map[string]*ProtoMessage, name string) *ProtoMessage {
		rs, ok := m[name]
		if !ok {
			panic(fmt.Errorf("unknown name of ProtoMessage: %s", name))
		}
		return rs
	}

	for _, service := range file.GetService() {
		s := &ProtoService{}
		s.Name = service.GetName()
		s.PascalName = MakeFirstLowerCase(s.Name)
		pkg.Services = append(pkg.Services, s)

		for i, method := range service.GetMethod() {
			m := &ProtoMethod{}
			m.Index = i
			m.Name = method.GetName()
			m.PascalName = MakeFirstLowerCase(m.Name)
			//		m.InputStream = *method.ClientStreaming
			//		m.OutputStream = *method.ServerStreaming
			log.Printf("input:  type=%s package=%s", method.GetInputType(), file.GetPackage())
			log.Printf("output: type=%s package=%s", method.GetOutputType(), file.GetPackage())
			input := removePackagePrefix(method.GetInputType(), file.GetPackage())
			output := removePackagePrefix(method.GetOutputType(), file.GetPackage())
			m.Input = mustGetMessage(messages, input)
			m.Output = mustGetMessage(messages, output)
			s.Methods = append(s.Methods, m)
		}
	}
	return pkg
}

// MakeFirstLowerCase makes the first character in a string lower case
func MakeFirstLowerCase(s string) string {

	if len(s) < 2 {
		return strings.ToLower(s)
	}

	bts := []byte(s)

	lc := bytes.ToLower([]byte{bts[0]})
	rest := bts[1:]

	return string(bytes.Join([][]byte{lc, rest}, nil))
}

// cleanPackageName lifted from gogo generator
// https://github.com/gogo/protobuf/blob/master/protoc-gen-gogo/generator/generator.go#L695
func cleanPackageName(name string) string {
	name = strings.Map(badToUnderscore, name)
	// Identifier must not be keyword: insert _.
	if isGoKeyword[name] {
		name = "_" + name
	}
	// Identifier must not begin with digit: insert _.
	if r, _ := utf8.DecodeRuneInString(name); unicode.IsDigit(r) {
		name = "_" + name
	}
	return name
}

// badToUnderscore lifted from gogo generator
func badToUnderscore(r rune) rune {
	if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
		return r
	}
	return '_'
}

func validateProtoMessage(pm *ProtoMessage) error {
	if pm == nil {
		return fmt.Errorf("nil ProtoMessage")
	}
	if pm.Name == "" {
		return fmt.Errorf("empty name of ProtoMessage")
	}
	return nil
}
