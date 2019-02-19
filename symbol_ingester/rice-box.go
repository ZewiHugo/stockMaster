package main

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    "createKeyspace.tmpl",
		FileModTime: time.Unix(1549026591, 0),

		Content: string("CREATE KEYSPACE IF NOT EXISTS {{.KeyspaceName}}\n    WITH REPLICATION = {\n        'class' : '{{.KeyspaceClass}}',\n        'replication_factor': {{.ReplicationFactor}}\n    };"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1548768061, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // "createKeyspace.tmpl"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`cassandra/template/`, &embedded.EmbeddedBox{
		Name: `cassandra/template/`,
		Time: time.Unix(1548768061, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"createKeyspace.tmpl": file2,
		},
	})
}
