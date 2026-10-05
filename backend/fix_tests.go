package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	rootDir := "./internal"
	
	// Replace `, nil)` with `, &storage.MockStorageService{})` in usecase.NewReportUsecase calls
	re := regexp.MustCompile(`(usecase\.NewReportUsecase\s*\(\s*[^,]+,\s*[^,]+,\s*[^,)]+),\s*nil\s*\)`)

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			
			newContent := re.ReplaceAll(content, []byte("$1, &storage.MockStorageService{})"))
			
			if !bytes.Equal(content, newContent) {
				// We also need to add import for "github.com/dbd-guard/backend/internal/infrastructure/storage"
				importStr := "\"github.com/dbd-guard/backend/internal/infrastructure/storage\""
				if !bytes.Contains(newContent, []byte(importStr)) {
					// Hacky way to inject import
					newContent = bytes.Replace(newContent, []byte("import ("), []byte("import (\n\t"+importStr), 1)
				}
				
				fmt.Printf("Updated %s\n", path)
				ioutil.WriteFile(path, newContent, 0644)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
