package cachectl

import (
	"fmt"
	"runtime"
)

func PrintCachectlVersion() {
	fmt.Printf(`cachectl %s
Compiler: %s %s
Copyright (C) 2014 Tatsuhiko Kubo <cubicdaiya@gmail.com>
`,
		Version,
		runtime.Compiler,
		runtime.Version())
}
