package main

import "sort"

// A FileSet represents a set of source files.
// Methods of file sets are synchronized; multiple goroutines
// may invoke them concurrently.
//
type FileSet struct {
	Base  int     // base offset for the next file
	Files []*File // list of files in the order added to the set
	last  *File   // cache of last file looked up
}

// NewFileSet creates a new file set.
func NewFileSet() *FileSet {
	return &FileSet{Base: 1} // 0 == NoPos
}

// AddFile adds a new file with a given filename, base offset, and file size
// to the file set s and returns the file. Multiple files may have the same
// name. The base offset must not be smaller than the FileSet's Base, and
// size must not be negative. As a special case, if a negative base is provided,
// the current value of the FileSet's Base is used instead.
//
// Adding the file will set the file set's Base value to base + size + 1
// as the minimum base value for the next file. The following relationship
// exists between a Pos value p for a given file offset offs:
//
//	int(p) = base + offs
//
// with offs in the range [0, size] and thus p in the range [base, base+size].
// For convenience, File.Pos may be used to create file-specific position
// values from a file offset.
//
func (s *FileSet) AddFile(filename string, base, size int) *File {
	if base < 0 {
		base = s.Base
	}
	if base < s.Base || size < 0 {
		panic("illegal base or size")
	}
	// base >= s.base && size >= 0
	f := &File{Name: filename, Base: base, Size: size, Lines: []int{0}}
	base += size + 1 // +1 because EOF also has a position
	if base < 0 {
		panic("token.Pos offset overflow (> 2G of source code in file set)")
	}
	// add the file to the file set
	s.Base = base
	s.Files = append(s.Files, f)
	s.last = f
	return f
}

func searchFiles(a []*File, x int) int {
	return sort.Search(len(a), func(i int) bool { return a[i].Base > x }) - 1
}

func (s *FileSet) file(p Pos) *File {
	// common case: p is in last file
	if f := s.last; f != nil && f.Base <= int(p) && int(p) <= f.Base+f.Size {
		return f
	}
	// p is not in last file - search all files
	if i := searchFiles(s.Files, int(p)); i >= 0 {
		f := s.Files[i]
		// f.Base <= int(p) by definition of searchFiles
		if int(p) <= f.Base+f.Size {
			s.last = f // race is ok - s.last is only a cache
			return f
		}
	}
	return nil
}

// File returns the file that contains the position p.
// If no such file is found (for instance for p == NoPos),
// the result is nil.
//
func (s *FileSet) File(p Pos) (f *File) {
	if p != NoPos {
		f = s.file(p)
	}
	return
}

// Position converts a Pos p in the fileset into a Position value.
func (s *FileSet) Position(p Pos) (pos Position) {
	if f := s.file(p); p != NoPos && f != nil {
		return f.position(p)
	}
	return
}
