package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

func RemoveModule(name string) {
	modulePath := filepath.Join("internal", "modules", name)

	// 1. Remove module folder
	if _, err := os.Stat(modulePath); os.IsNotExist(err) {
		fmt.Println("[WARN] Module folder not found:", modulePath)
	} else {
		if err := os.RemoveAll(modulePath); err != nil {
			panic(err)
		}
		fmt.Println("[OK] Removed module folder:", modulePath)
	}

	// 2. Update modules.go
	err := RemoveModuleFromRegistry(name)
	if err != nil {
		panic(err)
	}

	fmt.Println("[DONE] Module removed:", name)
}
