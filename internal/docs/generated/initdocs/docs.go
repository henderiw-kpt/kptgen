// Code generated by "mdtogo"; DO NOT EDIT.
package initdocs

var InitShort = `Initialize an empty package.`
var InitLong = `
  kptgen init [DIR] [flags]

Args:

  DIR:
    init fails if DIR does not already exist. Defaults to the current working directory.

Flags:

  --description
    Short description of the package. (default "sample description")
  
  --keywords
    A list of keywords describing the package.
  
  --site
    Link to page with information about the package.
`
var InitExamples = `

  # Creates a new Kptfile with metadata in the cockroachdb directory.
  $ mkdir cockroachdb; kptgen init cockroachdb --keywords "cockroachdb,nosql,db"  \
      --description "my cockroachdb implementation"

  # Creates a new Kptfile without metadata in the current directory.
  $ kptgen init
`
