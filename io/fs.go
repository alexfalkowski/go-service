package io

import "io/fs"

// ErrClosed aliases [fs.ErrClosed].
var ErrClosed = fs.ErrClosed

// ErrExist aliases [fs.ErrExist].
var ErrExist = fs.ErrExist

// ErrInvalid aliases [fs.ErrInvalid].
var ErrInvalid = fs.ErrInvalid

// ErrNotExist aliases [fs.ErrNotExist].
var ErrNotExist = fs.ErrNotExist

// ErrPermission aliases [fs.ErrPermission].
var ErrPermission = fs.ErrPermission

// ErrSkipAll aliases [fs.SkipAll].
var ErrSkipAll = fs.SkipAll

// ErrSkipDir aliases [fs.SkipDir].
var ErrSkipDir = fs.SkipDir

// ModeAppend aliases [fs.ModeAppend].
const ModeAppend = fs.ModeAppend

// ModeCharDevice aliases [fs.ModeCharDevice].
const ModeCharDevice = fs.ModeCharDevice

// ModeDevice aliases [fs.ModeDevice].
const ModeDevice = fs.ModeDevice

// ModeDir aliases [fs.ModeDir].
const ModeDir = fs.ModeDir

// ModeExclusive aliases [fs.ModeExclusive].
const ModeExclusive = fs.ModeExclusive

// ModeIrregular aliases [fs.ModeIrregular].
const ModeIrregular = fs.ModeIrregular

// ModeNamedPipe aliases [fs.ModeNamedPipe].
const ModeNamedPipe = fs.ModeNamedPipe

// ModePerm aliases [fs.ModePerm].
const ModePerm = fs.ModePerm

// ModeSetgid aliases [fs.ModeSetgid].
const ModeSetgid = fs.ModeSetgid

// ModeSetuid aliases [fs.ModeSetuid].
const ModeSetuid = fs.ModeSetuid

// ModeSocket aliases [fs.ModeSocket].
const ModeSocket = fs.ModeSocket

// ModeSticky aliases [fs.ModeSticky].
const ModeSticky = fs.ModeSticky

// ModeSymlink aliases [fs.ModeSymlink].
const ModeSymlink = fs.ModeSymlink

// ModeTemporary aliases [fs.ModeTemporary].
const ModeTemporary = fs.ModeTemporary

// ModeType aliases [fs.ModeType].
const ModeType = fs.ModeType

// DirEntry aliases [fs.DirEntry].
type DirEntry = fs.DirEntry

// FS aliases [fs.FS] to keep read-only filesystem contracts on the go-service import path.
type FS = fs.FS

// File aliases [fs.File].
type File = fs.File

// FileInfo aliases [fs.FileInfo].
type FileInfo = fs.FileInfo

// FileMode aliases [fs.FileMode].
type FileMode = fs.FileMode

// GlobFS aliases [fs.GlobFS].
type GlobFS = fs.GlobFS

// PathError aliases [fs.PathError].
type PathError = fs.PathError

// ReadDirFS aliases [fs.ReadDirFS].
type ReadDirFS = fs.ReadDirFS

// ReadDirFile aliases [fs.ReadDirFile].
type ReadDirFile = fs.ReadDirFile

// ReadFileFS aliases [fs.ReadFileFS].
type ReadFileFS = fs.ReadFileFS

// ReadLinkFS aliases [fs.ReadLinkFS].
type ReadLinkFS = fs.ReadLinkFS

// StatFS aliases [fs.StatFS].
type StatFS = fs.StatFS

// SubFS aliases [fs.SubFS].
type SubFS = fs.SubFS

// WalkDirFunc aliases [fs.WalkDirFunc].
type WalkDirFunc = fs.WalkDirFunc

// FileInfoToDirEntry wraps [fs.FileInfoToDirEntry].
func FileInfoToDirEntry(info FileInfo) DirEntry {
	return fs.FileInfoToDirEntry(info)
}

// FormatDirEntry wraps [fs.FormatDirEntry].
func FormatDirEntry(dir DirEntry) string {
	return fs.FormatDirEntry(dir)
}

// FormatFileInfo wraps [fs.FormatFileInfo].
func FormatFileInfo(info FileInfo) string {
	return fs.FormatFileInfo(info)
}

// Glob wraps [fs.Glob].
func Glob(fsys FS, pattern string) ([]string, error) {
	return fs.Glob(fsys, pattern)
}

// Lstat wraps [fs.Lstat].
func Lstat(fsys FS, name string) (FileInfo, error) {
	return fs.Lstat(fsys, name)
}

// ReadDir wraps [fs.ReadDir].
func ReadDir(fsys FS, name string) ([]DirEntry, error) {
	return fs.ReadDir(fsys, name)
}

// ReadFile wraps [fs.ReadFile], preserving file contents without trimming whitespace.
func ReadFile(fsys FS, name string) ([]byte, error) {
	return fs.ReadFile(fsys, name)
}

// ReadLink wraps [fs.ReadLink].
func ReadLink(fsys FS, name string) (string, error) {
	return fs.ReadLink(fsys, name)
}

// Stat wraps [fs.Stat].
func Stat(fsys FS, name string) (FileInfo, error) {
	return fs.Stat(fsys, name)
}

// Sub wraps [fs.Sub].
func Sub(fsys FS, dir string) (FS, error) {
	return fs.Sub(fsys, dir)
}

// ValidPath wraps [fs.ValidPath].
func ValidPath(name string) bool {
	return fs.ValidPath(name)
}

// WalkDir wraps [fs.WalkDir].
func WalkDir(fsys FS, root string, fn WalkDirFunc) error {
	return fs.WalkDir(fsys, root, fn)
}
