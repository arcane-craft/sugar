package lib

import (
	"fmt"
	"io"
	"os"
	"slices"
)

func GenNewLine() string {
	return "\n"
}

func GenBuildFlags(predicate bool) string {
	return fmt.Sprintf("%s\n\n", BuildDirective(predicate))
}

func GenTmpBuildFlags(predicate bool) string {
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
	blocks, err := proc(file, addImports)
	if err != nil {
		return err
	}

	slices.SortFunc(blocks, func(a, b *ReplaceBlock) int {
		return a.Old.Start.Offset - b.Old.Start.Offset
	})

	if _, err := writer.Write([]byte(GenTmpBuildFlags(true))); err != nil {
		return fmt.Errorf("writer.Write() failed: %w", err)
	}
	nextOffset := info.ImportExtent.End.Offset + 1
	if info.BuildFlag != nil {
		_, err := file.Seek(int64(info.BuildFlag.End.Offset+2), io.SeekStart)
		if err != nil {
			return fmt.Errorf("file.Seek() failed: %w", err)
		}
		nextOffset = info.ImportExtent.End.Offset - info.BuildFlag.End.Offset - 1
	}
	if _, err := io.CopyN(writer, file, int64(nextOffset)); err != nil {
		return fmt.Errorf("io.CopyN() failed: %w", err)
	}
	for p, n := range addImports {
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
