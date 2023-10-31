// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package setup

import (
	"fmt"
	"net/http"

	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/embed"
	"github.com/rande/gonode/core/helper"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/test/fixtures"
	"github.com/zenazn/goji/web"
)

func Configure(l *goapp.Lifecycle, conf *config.Config) {

	l.Prepare(func(app *goapp.App) error {
		if !conf.Test {
			return nil
		}

		mux := app.Get("goji.mux").(*web.Mux)
		manager := app.Get("gonode.manager").(*base.PgNodeManager)

		prefix := ""

		mux.Post(prefix+"/setup/uninstall", func(res http.ResponseWriter, req *http.Request) {
			var err error
			res.Header().Set("Content-Type", "application/json")

			prefix := conf.Databases["master"].Prefix

			_, err = manager.Db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS "%s_nodes"`, prefix))
			helper.PanicOnError(err)
			_, err = manager.Db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS "%s_nodes_audit"`, prefix))
			helper.PanicOnError(err)
			_, err = manager.Db.Exec(fmt.Sprintf(`DROP INDEX IF EXISTS "%s_uuid_idx"`, prefix))
			helper.PanicOnError(err)
			_, err = manager.Db.Exec(fmt.Sprintf(`DROP INDEX IF EXISTS "%s_uuid_current_idx"`, prefix))
			helper.PanicOnError(err)
			_, err = manager.Db.Exec(fmt.Sprintf(`DROP SEQUENCE IF EXISTS "%s_nodes_id_seq" CASCADE`, prefix))
			helper.PanicOnError(err)
			_, err = manager.Db.Exec(fmt.Sprintf(`DROP SEQUENCE IF EXISTS "%s_nodes_audit_id_seq" CASCADE`, prefix))
			helper.PanicOnError(err)

			helper.SendWithHttpCode(res, http.StatusOK, "Successfully delete tables!")
		})

		mux.Post(prefix+"/setup/install", func(res http.ResponseWriter, req *http.Request) {
			var err error
			//var r sql.Result

			res.Header().Set("Content-Type", "application/json")

			prefix := conf.Databases["master"].Prefix

			// Create my table
			_, err = manager.Db.Exec(fmt.Sprintf(`CREATE SEQUENCE "%s_nodes_id_seq" INCREMENT 1 MINVALUE 0 MAXVALUE 2147483647 START 1 CACHE 1`, prefix))
			helper.PanicOnError(err)

			_, err = manager.Db.Exec(fmt.Sprintf(`CREATE TABLE "%s_nodes" (
				"id" INTEGER DEFAULT nextval('%s_nodes_id_seq'::regclass) NOT NULL UNIQUE,
				"uuid" UUid NOT NULL,
				"type" CHARACTER VARYING( 64 ) COLLATE "pg_catalog"."default" NOT NULL,
				"name" CHARACTER VARYING( 2044 ) COLLATE "pg_catalog"."default" DEFAULT ''::CHARACTER VARYING NOT NULL,
				"enabled" BOOLEAN DEFAULT 'true' NOT NULL,
				"current" BOOLEAN DEFAULT 'false' NOT NULL,
				"revision" INTEGER DEFAULT '1' NOT NULL,
				"version" INTEGER DEFAULT '1' NOT NULL,
				"status" INTEGER DEFAULT '0' NOT NULL,
				"deleted" BOOLEAN DEFAULT 'false' NOT NULL,
				"data" jsonb DEFAULT '{}'::jsonb NOT NULL,
				"meta" jsonb DEFAULT '{}'::jsonb NOT NULL,
				"modules" jsonb DEFAULT '{}'::jsonb NOT NULL,
				"access" text[] DEFAULT '{}' NOT NULL,
				"slug" CHARACTER VARYING( 256 ) COLLATE "default" NOT NULL,
				"path" CHARACTER VARYING( 2000 ) COLLATE "default" NOT NULL,
				"source" UUid,
				"set_uuid" UUid,
				"parent_uuid" UUid,
				"parents" UUid[],
				"created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL,
				"created_by" UUid NOT NULL,
				"updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL,
				"updated_by" UUid NOT NULL,
				"weight" INTEGER DEFAULT '0' NOT NULL,
				PRIMARY KEY ( "id" ),
				CONSTRAINT "%s_slug" UNIQUE( "parent_uuid","slug","revision" ),
				CONSTRAINT "%s_uuid" UNIQUE( "revision","uuid" )
			)`, prefix, prefix, prefix, prefix))
			helper.PanicOnError(err)

			_, err = manager.Db.Exec(fmt.Sprintf(`CREATE INDEX "%s_uuid_idx" ON "%s_nodes" USING btree( "uuid" ASC NULLS LAST )`, prefix, prefix))
			helper.PanicOnError(err)

			_, err = manager.Db.Exec(fmt.Sprintf(`CREATE INDEX "%s_uuid_current_idx" ON "%s_nodes" USING btree( "uuid" ASC NULLS LAST, "current" ASC NULLS LAST )`, prefix, prefix))
			helper.PanicOnError(err)

			_, err = manager.Db.Exec(fmt.Sprintf(`CREATE INDEX "%s_access_idx" ON "%s_nodes" USING GIN("access")`, prefix, prefix))
			helper.PanicOnError(err)

			// Create Index
			_, err = manager.Db.Exec(fmt.Sprintf(`CREATE SEQUENCE "%s_nodes_audit_id_seq" INCREMENT 1 MINVALUE 0 MAXVALUE 2147483647 START 1 CACHE 1`, prefix))
			helper.PanicOnError(err)
			_, err = manager.Db.Exec(fmt.Sprintf(`CREATE TABLE "%s_nodes_audit" (
				"id" INTEGER DEFAULT nextval('%s_nodes_audit_id_seq'::regclass) NOT NULL UNIQUE,
				"uuid" UUid NOT NULL,
				"type" CHARACTER VARYING( 64 ) COLLATE "pg_catalog"."default" NOT NULL,
				"name" CHARACTER VARYING( 2044 ) COLLATE "pg_catalog"."default" DEFAULT ''::CHARACTER VARYING NOT NULL,
				"enabled" BOOLEAN DEFAULT 'true' NOT NULL,
				"current" BOOLEAN DEFAULT 'false' NOT NULL,
				"revision" INTEGER DEFAULT '1' NOT NULL,
				"version" INTEGER DEFAULT '1' NOT NULL,
				"status" INTEGER DEFAULT '0' NOT NULL,
				"deleted" BOOLEAN DEFAULT 'false' NOT NULL,
				"data" jsonb DEFAULT '{}'::jsonb NOT NULL,
				"meta" jsonb DEFAULT '{}'::jsonb NOT NULL,
				"modules" jsonb DEFAULT '{}'::jsonb NOT NULL,
				"access" text[] DEFAULT '{}' NOT NULL,
				"slug" CHARACTER VARYING( 256 ) COLLATE "default" NOT NULL,
				"path" CHARACTER VARYING( 2000 ) COLLATE "default" NOT NULL,
				"source" UUid,
				"set_uuid" UUid,
				"parent_uuid" UUid,
				"parents" UUid[],
				"created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL,
				"created_by" UUid NOT NULL,
				"updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL,
				"updated_by" UUid NOT NULL,
				"weight" INTEGER DEFAULT '0' NOT NULL,
				PRIMARY KEY ( "id" )
			)`, prefix, prefix))
			helper.PanicOnError(err)

			if err != nil {
				helper.SendWithHttpCode(res, http.StatusInternalServerError, "create tables: "+err.Error())
			} else {
				helper.SendWithHttpCode(res, http.StatusOK, "Successfully create tables!")
			}
		})

		mux.Post(prefix+"/setup/data/purge", func(res http.ResponseWriter, req *http.Request) {

			manager := app.Get("gonode.manager").(*base.PgNodeManager)

			prefix := conf.Databases["master"].Prefix

			tx, _ := manager.Db.Begin()
			manager.Db.Exec(fmt.Sprintf(`DELETE FROM "%s_nodes"`, prefix))
			manager.Db.Exec(fmt.Sprintf(`DELETE FROM "%s_nodes_audit"`, prefix))
			err := tx.Commit()

			if err != nil {
				helper.SendWithHttpCode(res, http.StatusInternalServerError, err.Error())
			} else {
				helper.SendWithHttpCode(res, http.StatusOK, "Data purged!")
			}
		})

		mux.Post(prefix+"/setup/data/load", func(res http.ResponseWriter, req *http.Request) {
			manager := app.Get("gonode.manager").(*base.PgNodeManager)
			nodes := manager.FindBy(manager.SelectBuilder(base.NewSelectOptions()), 0, 10)

			if nodes.Len() != 0 {
				helper.SendWithHttpCode(res, http.StatusInternalServerError, "Table contains data, purge the data first!")
				return
			}

			err := fixtures.LoadFixtures(manager, 100)

			if err != nil {
				helper.SendWithHttpCode(res, http.StatusInternalServerError, err.Error())
			} else {
				helper.SendWithHttpCode(res, http.StatusOK, "Data loaded")
			}
		})

		return nil
	})

	l.Register(func(app *goapp.App) error {
		app.Get("gonode.embeds").(*embed.Embeds).Add("setup", GetEmbedFS())

		return nil
	})
}
