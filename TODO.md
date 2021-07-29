## file

edit.go: the comment character used by default

Copy with options: (overwrite)

type option int

const (
	_ option = iota
	ARCHIVE // preserve all attributes

	BACKUP  // make a backup of each existing destination file

	ATTRIBUTES // only copy attributes
	HLINK     // make hard link instead of copying
	SLINK     // make symbolic links instead of copying

	DEREFERENCE // always follow symbolic links in SOURCE

	RECURSIVE // copy directories recursively

	UPDATE // copy only when the source file is newer than the destination file
	       // or when the destination file is missing

)
