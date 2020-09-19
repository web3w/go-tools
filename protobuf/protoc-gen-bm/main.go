package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gisvr/go-tools/protobuf/pkg/gen"
	"github.com/gisvr/go-tools/protobuf/pkg/generator"
	bmgen "github.com/gisvr/go-tools/protobuf/protoc-gen-bm/generator"
)

func main() {
	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.Parse()
	if *versionFlag {
		fmt.Println(generator.Version)
		os.Exit(0)
	}

	g := bmgen.BmGenerator()
	gen.Main(g)
}
