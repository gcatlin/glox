package main

// A File is a handle for a file belonging to a FileSet.
// A File has a name, size, and line offset table.
//
type File struct {
	Name  string // file name as provided to AddFile
	Base  int    // Pos value range for this file is [base...base+size]
	Size  int    // file size as provided to AddFile
	Lines []int  // lines contains the offset of the first character for each line (the first entry is always 0)
}

// AddLine adds the line offset for a new line.
// The line offset must be larger than the offset for the previous line
// and smaller than the file size; otherwise the line offset is ignored.
//
func (f *File) AddLine(offset int) {
	if i := len(f.Lines); (i == 0 || f.Lines[i-1] < offset) && offset < f.Size {
		f.Lines = append(f.Lines, offset)
	}
}

// SetLines sets the line offsets for a file and reports whether it succeeded.
// The line offsets are the offsets of the first character of each line;
// for instance for the content "ab\nc\n" the line offsets are {0, 3}.
// An empty file has an empty line offset table.
// Each line offset must be larger than the offset for the previous line
// and smaller than the file size; otherwise SetLines fails and returns false.
// Callers must not mutate the provided slice after SetLines returns.
//
func (f *File) SetLines(lines []int) bool {
	// verify validity of lines table
	size := f.Size
	for i, offset := range lines {
		if i > 0 && offset <= lines[i-1] || size <= offset {
			return false
		}
	}

	// set lines table
	f.Lines = lines
	return true
}

// SetLinesForContent sets the line offsets for the given file content.
func (f *File) SetLinesForContent(content []byte) {
	var lines []int
	line := 0
	for offset, b := range content {
		if line >= 0 {
			lines = append(lines, line)
		}
		line = -1
		if b == '\n' {
			line = offset + 1
		}
	}

	// set lines table
	f.Lines = lines
}

// Pos returns the Pos value for the given file offset;
// the offset must be <= f.Size().
// f.Pos(f.Offset(p)) == p.
//
func (f *File) Pos(offset int) Pos {
	if offset > f.Size {
		panic("illegal file offset")
	}
	return Pos(f.Base + offset)
}

// Offset returns the offset for the given file position p;
// p must be a valid Pos value in that file.
// f.Offset(f.Pos(offset)) == offset.
func (f *File) Offset(p Pos) int {
	if int(p) < f.Base || int(p) > f.Base+f.Size {
		panic("illegal Pos value")
	}
	return int(p) - f.Base
}

// Line returns the Line value for the given file position p;
// p must be a Pos value in that file or NoPos.
//
func (f *File) Line(p Pos) (line Line) {
	pos := f.Position(p)
	line.Filename = pos.Filename
	line.Offset = f.Lines[pos.Line]
	line.Line = pos.Line
	if len(f.Lines) > line.Line {
		line.Length = f.Lines[line.Line+1]
	} else {
		line.Length = f.Size - line.Offset
	}
	return
}

// unpack returns the filename and line and column number for a file offset.
func (f *File) unpack(offset int) (filename string, line, column int) {
	filename = f.Name
	if i := searchInts(f.Lines, offset); i >= 0 {
		line, column = i+1, offset-f.Lines[i]+1
	}
	return
}

func (f *File) position(p Pos) (pos Position) {
	offset := int(p) - f.Base
	pos.Offset = offset
	pos.Filename, pos.Line, pos.Column = f.unpack(offset)
	return
}

// Position returns the Position value for the given file position p.
func (f *File) Position(p Pos) (pos Position) {
	if p != NoPos {
		if int(p) < f.Base || f.Base+f.Size < int(p) {
			panic("illegal Pos value")
		}
		pos = f.position(p)
	}
	return
}

func searchInts(a []int, x int) int {
	i, j := 0, len(a)
	for i < j {
		h := i + (j-i)/2 // avoid overflow when computing h
		// i ≤ h < j
		if a[h] <= x {
			i = h + 1
		} else {
			j = h
		}
	}
	return i - 1
}
