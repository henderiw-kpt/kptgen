// Code generated by "mdtogo"; DO NOT EDIT.
package clonedocs

var CloneShort = `Clone copies the yaml files of a git repo to a target directory.`
var CloneLong = `
  kptgen clone GIT_REPO_URL TARGET_DIR [flags]

Args:

  GIT_REPO:
    The source repo of manifests.
  
  TARGET_DIR:
    The target diretcory the package would be created in
`
var CloneExamples = `

  # Creates a new Kpt Package from kubebuilder manifest files.
  $ kptgen clone https://github.com/henderiw-kpt/kptgen-templates ./blueprint/topology
`
