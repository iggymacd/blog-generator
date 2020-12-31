package generator

import (
	"bytes"
	"testing"
	"text/template"
)

// func TestModelGenerator_Generate(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		g       *ModelGenerator
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := tt.g.Generate(); (err != nil) != tt.wantErr {
// 				t.Errorf("ModelGenerator.Generate() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestWriteModelImport(t *testing.T) {
// 	type args struct {
// 		// entity *Entity
// 		config *ModelConfig
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantW   string
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "first",
// 			args: args{
// 				// &Entity{
// 				// 	Name:       "User",
// 				// 	Meta:       &EntityMeta{},
// 				// 	Attributes: []*Attribute{},
// 				// },
// 				&ModelConfig{
// 					Entity: &Entity{
// 						Name: "User",
// 						Meta: &EntityMeta{},
// 						Attributes: []*Attribute{
// 							{
// 								Name: "gender",
// 								Meta: &AttributeMeta{
// 									AttributeType: "enum",
// 									DisplayName:   "Gender",
// 								},
// 							},
// 							{
// 								Name: "vital_status",
// 								Meta: &AttributeMeta{
// 									AttributeType: "enum",
// 									DisplayName:   "Vital Status",
// 								},
// 							},
// 						},
// 					},
// 					Template:    &template.Template{},
// 					Destination: "",
// 					ModelName:   "",
// 					Writer:      &ModelWriter{},
// 				},
// 			},
// 			wantW: `import 'package:freezed_annotation/freezed_annotation.dart';

// import '../../../../core/types/gender.dart';
// import '../../../../core/types/vital_status.dart';
// import '../../../../core/util/mapper.dart';
// import '../../domain/entities/character.dart';

// part 'user_model.freezed.dart';
// part 'user_model.g.dart';
// `,
// 			wantErr: false,
// 		},
// 	}

// 	targpath := "/lib/features/home/data/models"
// 	basepath := "/lib/core/types"

// 	relpath, _ := filepath.Rel(targpath, basepath)
// 	fmt.Println("Relative Path:", relpath)

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			w := &bytes.Buffer{}
// 			if err := WriteModelImport(w, tt.args.config); (err != nil) != tt.wantErr {
// 				t.Errorf("WriteModelImport() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if gotW := w.String(); gotW != tt.wantW {
// 				t.Errorf("WriteModelImport() = %v, want %v", gotW, tt.wantW)
// 			}
// 		})
// 	}
// }

// func TestModelWriter_WriteModel(t *testing.T) {
// 	type args struct {
// 		path             string
// 		classDeclaration string
// 		importStatements string
// 		staticFields     string
// 		fields           string
// 		transformer      string
// 		t                *template.Template
// 	}
// 	tests := []struct {
// 		name    string
// 		mw      *ModelWriter
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "Generic template with no template data.",
// 			mw:&ModelWriter{
// 				ModelName: "User",
// 			},
// 			args:args{
// 				path:             "",
// 				classDeclaration: "",
// 				importStatements: "",
// 				staticFields:     "",
// 				fields:           "",
// 				transformer:      "",
// 				t:                &template.Template{},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := tt.mw.WriteModel(tt.args.path, tt.args.classDeclaration, tt.args.importStatements, tt.args.staticFields, tt.args.fields, tt.args.transformer, tt.args.t); (err != nil) != tt.wantErr {
// 				t.Errorf("ModelWriter.WriteModel() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

func TestModelWriter_WriteModel(t *testing.T) {
	templateContent := `import 'package:freezed_annotation/freezed_annotation.dart';
{{range .Dependencies}}{{"\n"}}import '{{.}}';{{end}}
import '../../../../core/util/mapper.dart';
import '../../domain/entities/{{toSnake .Name}}.dart';

part '{{toSnake .Name}}_model.freezed.dart';
part '{{toSnake .Name}}_model.g.dart';

@freezed
abstract class {{toCamel .Name}}Model with _${{toCamel .Name}}Model {
	const factory {{toCamel .Name}}Model({
	{{range .Attributes}}{{fieldWithDecoration .}}{{end}}
	@required String id,
	@required String name,
	@JsonKey(
		fromJson: Mapper.statusInType,
		toJson: Mapper.statusInString,
		name: 'status')
	@required
		VitalStatus vitalStatus,
	@JsonKey(fromJson: Mapper.genderInType, toJson: Mapper.genderInString)
	@required
		Gender gender,
	@required String type,
	@required String species,
	@required String image,
	}) = _CharacterModel;

	factory CharacterModel.fromJson(Map<String, dynamic> json) =>
		_$CharacterModelFromJson(json);

	factory CharacterModel.fromEntity(Character character) => CharacterModel(
		id: character.id,
		name: character.name,
		vitalStatus: character.vitalStatus,
		gender: character.gender,
		type: character.type,
		species: character.species,
		image: character.image,
		);
}

extension CharacterModelX on CharacterModel {
	Character toEntity() => Character(
		id: id,
		name: name,
		vitalStatus: vitalStatus,
		gender: gender,
		type: type,
		species: species,
		image: image,
		);
}
`
	template, err := template.New("user").Funcs(funcs).Parse(templateContent)
	if err != nil {
		t.Errorf("Could not parse template %s", templateContent)
	}
	type args struct {
		path             string
		classDeclaration string
		importStatements string
		staticFields     string
		fields           string
		transformer      string
		config           *ModelConfig
		// t                *template.Template
	}
	tests := []struct {
		name    string
		mw      *ModelWriter
		args    args
		wantW   string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Generate template with no data.",
			mw: &ModelWriter{
				ModelName: "user",
			},
			args: args{
				path:             "",
				classDeclaration: "",
				importStatements: "",
				staticFields:     "",
				fields:           "",
				transformer:      "",
				config: &ModelConfig{
					Entity: &Entity{
						Name: "user",
						Meta: &EntityMeta{
							Description: "",
							DisplayName: "User",
							Permissions: []string{},
							Attributes:  []*Attribute{},
							Path:        "/lib/features/home/data/models",
						},
						Attributes: []*Attribute{
							{
								Name: "id",
								Meta: &AttributeMeta{
									AttributeType: "string",
									DisplayName:   "Id",
									// Path: "lib/core/types",
								},
							},
							{
								Name: "gender",
								Meta: &AttributeMeta{
									AttributeType: "multi",
									DisplayName:   "Gender",
									Association: &Association{
										Type:        "multi",
										Cardinality: "one to many",
										Entity: &Entity{
											Name: "gender",
											Meta: &EntityMeta{
												Description: "",
												DisplayName: "Gender",
												Permissions: []string{},
												Attributes:  []*Attribute{},
												Path:        "/lib/core/types",
											},
											Attributes: []*Attribute{},
										},
									},
									// Path: "lib/core/types",
								},
							},
							{
								Name: "vital_status",
								Meta: &AttributeMeta{
									AttributeType: "single",
									DisplayName:   "Vital Status",
									Association: &Association{
										Type:        "single",
										Cardinality: "many to one",
										Entity: &Entity{
											Name: "vital_status",
											Meta: &EntityMeta{
												Description: "",
												DisplayName: "Vital Status",
												Permissions: []string{},
												Attributes:  []*Attribute{},
												Path:        "/lib/core/types",
											},
											Attributes: []*Attribute{},
										},
									},
									// Path: "lib/core/types",
								},
							},
						},
						// Dependencies: []string{},
					},
					Template:    template,
					Destination: "",
					ModelName:   "",
				},
				// t:                &template.Template{},
			},
			wantW: `import 'package:freezed_annotation/freezed_annotation.dart';

import '../../../../core/types/gender.dart';
import '../../../../core/types/vital_status.dart';
import '../../../../core/util/mapper.dart';
import '../../domain/entities/user.dart';

part 'user_model.freezed.dart';
part 'user_model.g.dart';

@freezed
abstract class UserModel with _$UserModel {
	const factory UserModel({
	@required String id,
	@required String name,
	@JsonKey(
		fromJson: Mapper.statusInType,
		toJson: Mapper.statusInString,
		name: 'status')
	@required
		VitalStatus vitalStatus,
	@JsonKey(fromJson: Mapper.genderInType, toJson: Mapper.genderInString)
	@required
		Gender gender,
	@required String type,
	@required String species,
	@required String image,
	}) = _CharacterModel;

	factory CharacterModel.fromJson(Map<String, dynamic> json) =>
		_$CharacterModelFromJson(json);

	factory CharacterModel.fromEntity(Character character) => CharacterModel(
		id: character.id,
		name: character.name,
		vitalStatus: character.vitalStatus,
		gender: character.gender,
		type: character.type,
		species: character.species,
		image: character.image,
		);
}

extension CharacterModelX on CharacterModel {
	Character toEntity() => Character(
		id: id,
		name: name,
		vitalStatus: vitalStatus,
		gender: gender,
		type: type,
		species: species,
		image: image,
		);
}
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := tt.mw.WriteModel(tt.args.path, tt.args.classDeclaration, tt.args.importStatements, tt.args.staticFields, tt.args.fields, tt.args.transformer, w, tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("ModelWriter.WriteModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ModelWriter.WriteModel() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
