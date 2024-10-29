package lib

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"slices"
	"strings"
)

func GenNewLine() string {
	return "\n"
}

func GenBuildTags(predicate bool) string {
	return fmt.Sprintf("%s\n\n", BuildDirective(predicate))
}

func GenTmpBuildTags(predicate bool) string {
	return fmt.Sprintf("%s\n\n", TmpBuildDirective(predicate))
}

func GenImport(pkgName, pkgPath string) string {
	if len(pkgName) <= 0 {
		pkgName = "."
	}
	return fmt.Sprintf("import %s \"%s\"\n", pkgName, pkgPath)
}

func GenReturn(expr string) string {
	return fmt.Sprintf("return %s", expr)
}

func GenFuncLit(funcType string, body string) string {
	return fmt.Sprintf("%s {\n%s\n}", funcType, body)
}

type ReplaceBlock struct {
	Old Extent
	New string
}

func ReadExtent(reader io.ReaderAt, extent *Extent) (string, error) {
	bs := make([]byte, extent.End.Offset-extent.Start.Offset)
	n, err := reader.ReadAt(bs, int64(extent.Start.Offset))
	if err != nil {
		return "", fmt.Errorf("reader.ReadAt() failed: %w", err)
	}
	if n < len(bs) {
		return "", fmt.Errorf("the lenght of read bytes is not enough")
	}
	return string(bs), nil
}

func ReadExtentList(reader io.ReaderAt, extents []*Extent) ([]string, error) {
	var ret []string
	for _, e := range extents {
		str, err := ReadExtent(reader, e)
		if err != nil {
			return nil, fmt.Errorf("readExtent() failed: %w", err)
		}
		ret = append(ret, str)
	}
	return ret, nil
}

func GenerateSyntax[Syntax interface {
	fmt.Stringer
	comparable
}](info *FileInfo[Syntax], writer io.Writer,
	proc func(file *os.File, addImports map[string]string) ([]*ReplaceBlock, error)) error {
	file, err := os.Open(info.Path)
	if err != nil {
		return fmt.Errorf("os.Open() failed: %w", err)
	}
	defer file.Close()

	addImports := make(map[string]string)
	for p, n := range info.Imports {
		if path.Base(p) == n {
			n = " "
		}
		addImports[p] = n
	}
	blocks, err := proc(file, addImports)
	if err != nil {
		return err
	}

	slices.SortFunc(blocks, func(a, b *ReplaceBlock) int {
		return a.Old.Start.Offset - b.Old.Start.Offset
	})

	if _, err := writer.Write([]byte(GenTmpBuildTags(true))); err != nil {
		return fmt.Errorf("writer.Write() failed: %w", err)
	}
	nextOffset := int64(info.ImportExtent.Start.Offset)
	if info.BuildTag != nil {
		_, err := file.Seek(int64(info.BuildTag.End.Offset+2), io.SeekStart)
		if err != nil {
			return fmt.Errorf("file.Seek() failed: %w", err)
		}
		nextOffset = int64(info.ImportExtent.Start.Offset - info.BuildTag.End.Offset - 1)
	}
	if _, err := io.CopyN(writer, file, nextOffset); err != nil {
		return fmt.Errorf("io.CopyN() failed: %w", err)
	}
	if _, err := file.Seek(int64(info.ImportExtent.End.Offset+1), io.SeekStart); err != nil {
		return fmt.Errorf("file.Seek() failed: %w", err)
	}
	var importPaths []string
	for p := range addImports {
		importPaths = append(importPaths, p)
	}
	slices.SortFunc(importPaths, func(a, b string) int {
		if len(a) == len(b) {
			return strings.Compare(a, b)
		}
		return len(a) - len(b)
	})
	for _, p := range importPaths {
		n := addImports[p]
		if _, err := writer.Write([]byte(GenNewLine() + GenImport(n, p))); err != nil {
			return fmt.Errorf("writer.Write() failed: %w", err)
		}
	}
	lastOffset := info.ImportExtent.End.Offset + 1
	for _, b := range blocks {
		if _, err := io.CopyN(writer, file, int64(b.Old.Start.Offset-lastOffset)); err != nil {
			return fmt.Errorf("io.CopyN() failed: %w", err)
		}
		if _, err := writer.Write([]byte(b.New)); err != nil {
			return fmt.Errorf("writer.Write() failed: %w", err)
		}
		if _, err := file.Seek(int64(b.Old.End.Offset-b.Old.Start.Offset), 1); err != nil {
			return fmt.Errorf("file.Seek() failed: %w", err)
		}
		lastOffset = b.Old.End.Offset
	}
	if _, err := io.Copy(writer, file); err != nil {
		return fmt.Errorf("io.Copy() failed: %w", err)
	}
	return nil
}

func replaceBuildTags(tags map[string]*Extent, old string, new string) error {
	for file := range tags {
		if old == prodBuildTag && strings.HasSuffix(file, old+".go") {
			if err := os.Remove(file); err != nil {
				return fmt.Errorf("os.Remove(): %w", err)
			}
			continue
		}
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("os.ReadFile(): %w", err)
		}
		content = bytes.Replace(content, []byte(old), []byte(new), 1)
		if err := os.WriteFile(file, content, 0644); err != nil {
			return fmt.Errorf("os.WriteFile() %w", err)
		}
	}
	return nil
}
