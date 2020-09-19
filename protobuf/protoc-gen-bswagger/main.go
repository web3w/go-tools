package main

import (
	"flag"
	"fmt"
	ge "github.com/gisvr/go-tools/protobuf/protoc-gen-bswagger/generator"
	"os"

	"github.com/gisvr/go-tools/protobuf/pkg/gen"
	"github.com/gisvr/go-tools/protobuf/pkg/generator"
)

func main() {
	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.Parse()
	if *versionFlag {
		fmt.Println(generator.Version)
		os.Exit(0)
	}

	g := ge.NewSwaggerGenerator()
	gen.Main(g)
}
