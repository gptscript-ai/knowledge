//
// Copyright 2020 FoxyUtils ehf. All rights reserved.
//
// This is a commercial product and requires a license to operate.
// A trial license can be obtained at https://unidoc.io
//
// DO NOT EDIT: generated by unitwist Go source code obfuscator.
//
// Use of this source code is governed by the UniDoc End User License Agreement
// terms that can be accessed at https://unidoc.io/eula/

// Package diskstore implements tempStorage interface
// by using disk as a storage
package diskstore ;import (_af "github.com/unidoc/unioffice/common/tempstorage";_a "io/ioutil";_ae "os";_g "strings";);

// Open opens file from disk according to a path
func (_gb diskStorage )Open (path string )(_af .File ,error ){return _ae .OpenFile (path ,_ae .O_RDWR ,0644);};

// RemoveAll removes all files in the directory
func (_aea diskStorage )RemoveAll (dir string )error {if _g .HasPrefix (dir ,_ae .TempDir ()){return _ae .RemoveAll (dir );};return nil ;};

// TempFile creates a new temp file by calling ioutil TempFile
func (_fd diskStorage )TempFile (dir ,pattern string )(_af .File ,error ){return _a .TempFile (dir ,pattern );};

// SetAsStorage sets temp storage as a disk storage
func SetAsStorage (){_f :=diskStorage {};_af .SetAsStorage (&_f )};

// Add is not applicable in the diskstore implementation
func (_e diskStorage )Add (path string )error {return nil };

// TempFile creates a new temp directory by calling ioutil TempDir
func (_gd diskStorage )TempDir (pattern string )(string ,error ){return _a .TempDir ("",pattern )};type diskStorage struct{};