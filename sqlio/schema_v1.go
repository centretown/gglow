package sqlio

import "gglow/iohandler"

var schema_v1 = &Schema{
	Version:  iohandler.Version{Major: 1, Minor: 0, Patch: 0, Extension: 0},
	AlterSQL: []string{},
	CreateSQL: []string{
		`CREATE TABLE version (
major SMALLINT,
minor SMALLINT,
patch SMALLINT,
extension SMALLINT);`,

		`CREATE TABLE content (
uuid VARCHAR(40) NOT NULL,
folder VARCHAR(80) NOT NULL,
title VARCHAR(80) NOT NULL,
category VARCHAR(40) NOT NULL,
content TEXT,
PRIMARY KEY (uuid));`,
		"CREATE INDEX content_category ON content (category,folder,title);",
		"CREATE INDEX content_folder ON content (folder,title);",
		"CREATE INDEX content_title ON content (title);",
		"CREATE INDEX content_category_title ON content (category,title);",

		`CREATE TABLE tags (
title VARCHAR(80) NOT NULL,
uuid VARCHAR(40) NOT NULL,
PRIMARY KEY (title,uuid));`,
		"CREATE INDEX tags_uuid ON tags (uuid,title);",
	},
	DropSQL: []string{
		"DROP TABLE IF EXISTS version;",
		"DROP TABLE IF EXISTS content;",
		"DROP TABLE IF EXISTS tags;",
	},
	selectFolderSQL: "SELECT folder FROM content WHERE title = '..' AND category = 'effect' ORDER BY folder;",
	listEffectsSQL:  "SELECT title FROM content WHERE folder = '%s' AND category = 'effect' ORDER BY title;",
	readEffectSQL:   "SELECT folder,title,effect FROM effects WHERE folder = '%s' AND title = '%s'",
	existsEffectSQL: "SELECT title FROM effects WHERE folder = '%s' AND title = '%s';",
	updateEffectSQL: "UPDATE effects SET effect = '%s' WHERE folder = '%s' AND title = '%s'",
	insertEffectSQL: "INSERT INTO effects (folder, title, effect) VALUES('%s', '%s', '%s')",
}
