package sqlio

type Schema struct {
	Version      string
	Create       []string
	Folder       string
	ListEffects  string
	ReadEffect   string
	FindEffect   string
	InsertEffect string
	UpdateEffect string
}

var SchemaMap = make(map[string]*Schema)
var Schemas = []*Schema{
	{Version: "0.1",
		Create: []string{
			"DROP VIEW IF EXISTS palettes;",
			"DROP VIEW IF EXISTS folders;",
			"DROP TABLE IF EXISTS effects;",
			"DROP TABLE IF EXISTS colors;",
			`CREATE TABLE effects (
folder VARCHAR(80) NOT NULL,
title VARCHAR(80) NOT NULL,
effect TEXT,
PRIMARY KEY (folder, title)
);`,
			"CREATE INDEX effect_title ON effects (title);",
			`CREATE VIEW folders(folder, title) AS
SELECT folder, title
FROM effects
WHERE title = '..'
ORDER BY folder;
`,
		},
		Folder:       "SELECT folder FROM folders;",
		ListEffects:  "SELECT title FROM effects WHERE folder = '%s' ORDER BY folder;",
		ReadEffect:   "SELECT * FROM effects WHERE folder = '%s' AND title = '%s'",
		FindEffect:   "SELECT title FROM effects WHERE folder = '%s' AND title = '%s';",
		UpdateEffect: "UPDATE effects SET effect = '%s' WHERE folder = '%s' AND title = '%s'",
		InsertEffect: "INSERT INTO effects (folder, title, effect) VALUES('%s', '%s', '%s')",
	},
}

func init() {
}
