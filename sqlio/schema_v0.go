package sqlio

import "gglow/iohandler"

var schema_v0 = &Schema{
	Version:  iohandler.Version{Major: 0, Minor: 0, Patch: 0, Extension: 0},
	AlterSQL: []string{},
	CreateSQL: []string{
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
	DropSQL: []string{
		"DROP VIEW IF EXISTS palettes;",
		"DROP VIEW IF EXISTS folders;",
		"DROP TABLE IF EXISTS effects;",
		"DROP TABLE IF EXISTS colors;",
	},
	// selectFolderSQL: "SELECT title FROM effects WHERE folder = '%s' ORDER BY title;",
	selectFolderSQL: "SELECT folder,title FROM effects WHERE title = '..' ORDER BY folder;",
	listEffectsSQL:  "SELECT title,folder FROM effects WHERE folder = '%s' ORDER BY folder,title;",
	readEffectSQL:   "SELECT folder,title,effect FROM effects WHERE folder = '%s' AND title = '%s'",
	existsEffectSQL: "SELECT title FROM effects WHERE folder = '%s' AND title = '%s';",
	updateEffectSQL: "UPDATE effects SET effect = '%s' WHERE folder = '%s' AND title = '%s'",
	insertEffectSQL: "INSERT INTO effects (folder, title, effect) VALUES('%s', '%s', '%s')",
}
