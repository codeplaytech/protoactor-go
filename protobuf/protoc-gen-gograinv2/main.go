package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"

	google_protobuf "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	plugin "github.com/gogo/protobuf/protoc-gen-gogo/plugin"
	"github.com/gogo/protobuf/vanity/command"
)

func main() {
	req := command.Read()
	resp := generateCode(req, "_protoactor.go", true)
	command.Write(resp)
}

func removePackagePrefix(name string, pname string) string {
	return strings.Replace(name, "."+pname+".", "", 1)
}

func generateCode(req *plugin.CodeGeneratorRequest, filenameSuffix string, goFmt bool) *plugin.CodeGeneratorResponse {

	targetFiles := map[string]struct{}{}
	for _, name := range req.FileToGenerate {
		targetFiles[name] = struct{}{}
	}

	response := &plugin.CodeGeneratorResponse{}
	for _, f := range req.GetProtoFile() {
		showFileDescriptorProto(f)
		name := f.GetName()
		if _, isTarget := targetFiles[name]; !isTarget {
			continue
		}

		if !strings.HasSuffix(name, ".proto") {
			panic(fmt.Errorf("proto file must be with .proto suffix: %s, ", name))
		}

		s := generate(f)

		// we only generate grains for proto files containing valid service definition
		if len(f.GetService()) > 0 {
			fileName := strings.TrimSuffix(name, ".proto") + filenameSuffix
			r := &plugin.CodeGeneratorResponse_File{
				Content: &s,
				Name:    &fileName,
			}

			response.File = append(response.File, r)
		}
	}

	return response
}

func generate(file *google_protobuf.FileDescriptorProto) string {

	pkg := ProtoAst(file)

	t := template.New("grain")
	t, err := t.Parse(code)
	if err != nil {
		panic(err)
	}

	var doc bytes.Buffer
	if err := t.Execute(&doc, pkg); err != nil {
		panic(err)
	}
	s := doc.String()
	return s
}

func showFileDescriptorProto(file *google_protobuf.FileDescriptorProto) {
	log.Printf("===== file ======")
	log.Printf("name=%v", file.GetName())
	log.Printf("package=%v", file.GetPackage())
	log.Printf("dependency=%v", file.Dependency)
	log.Printf("PublicDependency=%v", file.PublicDependency)

	log.Printf("---- messges ----")
	for i, m := range file.GetMessageType() {
		log.Printf("\t%d: %s", i+1, m.GetName())
	}
	log.Printf("=================")
}
