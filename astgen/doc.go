// Package astgen is a highly experimental package that converts astorg
// datasets into DPMA-based phone configuration files.
//
// TODO: Refactor the functions in this package to avoid ingestion of whole
// datasets. Taking a dataset as an argument hides the actual requirements
// of each generator. Fixing this will probably require that we redesign
// the generator functions to cumulatively apply the astorg phone, person
// and role objects in separate functions.
package astgen
