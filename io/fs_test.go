package io_test

import (
	"testing"
	"testing/fstest"

	"github.com/alexfalkowski/go-service/v3/io"
	"github.com/alexfalkowski/go-service/v3/time"
	"github.com/stretchr/testify/require"
)

func TestSubReadsWithinDirectory(t *testing.T) {
	fsys := fstest.MapFS{
		"assets/file.txt": {Data: []byte(" contents \n")},
		"file.txt":        {Data: []byte("outside")},
	}

	sub, err := io.Sub(fsys, "assets")
	require.NoError(t, err)

	data, err := io.ReadFile(sub, "file.txt")
	require.NoError(t, err)
	require.Equal(t, " contents \n", string(data))

	_, err = io.Sub(fsys, "../assets")
	require.ErrorIs(t, err, io.ErrInvalid)
}

func TestReadDirAndGlobReturnMatchingEntries(t *testing.T) {
	fsys := fstest.MapFS{
		"dir/b.txt":     {},
		"dir/a.txt":     {},
		"dir/sub/c.txt": {},
		"outside.txt":   {},
	}

	entries, err := io.ReadDir(fsys, "dir")
	require.NoError(t, err)

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	require.Equal(t, []string{"a.txt", "b.txt", "sub"}, names)

	matches, err := io.Glob(fsys, "dir/*.txt")
	require.NoError(t, err)
	require.Equal(t, []string{"dir/a.txt", "dir/b.txt"}, matches)
}

func TestStatAndLstatDistinguishSymbolicLinks(t *testing.T) {
	fsys := fstest.MapFS{
		"file.txt": {Data: []byte("hello"), Mode: 0o640},
		"link.txt": {Data: []byte("file.txt"), Mode: io.ModeSymlink | 0o777},
	}

	info, err := io.Stat(fsys, "link.txt")
	require.NoError(t, err)
	require.Equal(t, io.FileMode(0o640), info.Mode())
	require.Equal(t, int64(5), info.Size(), "stat must describe the link target")

	link, err := io.Lstat(fsys, "link.txt")
	require.NoError(t, err)
	require.Equal(t, io.ModeSymlink|0o777, link.Mode())

	target, err := io.ReadLink(fsys, "link.txt")
	require.NoError(t, err)
	require.Equal(t, "file.txt", target)
}

func TestFileInfoConversionAndFormatting(t *testing.T) {
	fsys := fstest.MapFS{
		"file.txt": {
			Data:    []byte("hello"),
			Mode:    0o640,
			ModTime: time.Unix(0).UTC(),
		},
	}
	info, err := io.Stat(fsys, "file.txt")
	require.NoError(t, err)

	entry := io.FileInfoToDirEntry(info)
	require.Equal(t, "file.txt", entry.Name())
	entryInfo, err := entry.Info()
	require.NoError(t, err)
	require.Equal(t, info, entryInfo)

	require.Equal(t, "- file.txt", io.FormatDirEntry(entry))
	require.Equal(t, "-rw-r----- 5 1970-01-01 00:00:00 file.txt", io.FormatFileInfo(info))
}

func TestValidPath(t *testing.T) {
	for _, tt := range []struct {
		name  string
		path  string
		valid bool
	}{
		{name: "root", path: ".", valid: true},
		{name: "nested file", path: "assets/file.txt", valid: true},
		{name: "parent traversal", path: "../file.txt"},
		{name: "absolute path", path: "/file.txt"},
		{name: "empty element", path: "assets//file.txt"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.valid, io.ValidPath(tt.path))
		})
	}
}

func TestWalkDirHonorsSkipSentinels(t *testing.T) {
	fsys := fstest.MapFS{
		"root/a.txt":           {},
		"root/skip/hidden.txt": {},
		"root/z.txt":           {},
		"outside.txt":          {},
	}

	for _, tt := range []struct {
		name string
		skip error
		want []string
	}{
		{
			name: "directory",
			skip: io.ErrSkipDir,
			want: []string{"root", "root/a.txt", "root/skip", "root/z.txt"},
		},
		{
			name: "all",
			skip: io.ErrSkipAll,
			want: []string{"root", "root/a.txt", "root/skip"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var paths []string
			err := io.WalkDir(fsys, "root", func(path string, _ io.DirEntry, err error) error {
				if err != nil {
					return err
				}

				paths = append(paths, path)
				if path == "root/skip" {
					return tt.skip
				}

				return nil
			})
			require.NoError(t, err)
			require.Equal(t, tt.want, paths)
		})
	}
}
