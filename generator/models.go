package generator

import "text/template"

// ModelData holds the data for the model definition
type ModelData struct {
	Name string
	// Date       string
	// Short      string
	// Link       string
	// TimeToRead string
	ImportStatements, StaticFields, Fields, Transformer template.Template
}

// // Field holds the data for each field
// type Field struct {
// 	Name string
// }

// ModelGenerator Object
type ModelGenerator struct {
	Config *ModelConfig
}

// ModelConfig holds the configuration for the model definition
type ModelConfig struct {
	ModelData              *ModelData
	Template               *template.Template
	Destination, ModelName string
	// IsIndex                bool
	Writer *IndexWriter
}
