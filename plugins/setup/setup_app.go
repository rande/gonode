package setup

import (
	"fmt"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/helper"
	"github.com/rande/gonode/test/fixtures"
	"github.com/zenazn/goji/web"
	"net/http"
)

func ConfigureServer(l *goapp.Lifecycle, conf *config.ServerConfig) {

	l.Prepare(func(app *goapp.App) error {
		mux := app.Get("goji.mux").(*web.Mux)
		manager := app.Get("gonode.manager").(*core.PgNodeManager)

		prefix := ""

		mux.Put(prefix+"/uninstall", func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")
			res.Header().Set("Access-Control-Allow-Origin", "*")

			prefix := conf.Databases["master"].Prefix

			manager.Db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS "%s_nodes"`, prefix))
			manager.Db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS "%s_nodes_audit"`, prefix))
			manager.Db.Exec(fmt.Sprintf(`DROP INDEX IF EXISTS "%s_uuid_idx"`, prefix))
			manager.Db.Exec(fmt.Sprintf(`DROP INDEX IF EXISTS "%s_uuid_current_idx"`, prefix))
			manager.Db.Exec(fmt.Sprintf(`DROP SEQUENCE IF EXISTS "%s_nodes_id_seq" CASCADE`, prefix))
			manager.Db.Exec(fmt.Sprintf(`DROP SEQUENCE IF EXISTS "%s_nodes_audit_id_seq" CASCADE`, prefix))

			helper.SendWithHttpCode(res, http.StatusOK, "Successfully delete tables!")
		})

		mux.Put(prefix+"/install", func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")
			res.Header().Set("Access-Control-Allow-Origin", "*")

			prefix := conf.Databases["master"].Prefix
			tx, _ := manager.Db.Begin()

			// Create my table
			tx.Exec(fmt.Sprintf(`CREATE SEQUENCE "%s_nodes_id_seq" INCREMENT 1 MINVALUE 0 MAXVALUE 2147483647 START 1 CACHE 1`, prefix))
			tx.Exec(fmt.Sprintf(`CREATE TABLE "%s_nodes" (
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
				"slug" CHARACTER VARYING( 256 ) COLLATE "default" NOT NULL,
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

			tx.Exec(fmt.Sprintf(`CREATE INDEX "%s_uuid_idx" ON "%s_nodes" USING btree( "uuid" ASC NULLS LAST )`, prefix, prefix))
			tx.Exec(fmt.Sprintf(`CREATE INDEX "%s_uuid_current_idx" ON "%s_nodes" USING btree( "uuid" ASC NULLS LAST, "current" ASC NULLS LAST )`, prefix, prefix))

			// Create Index
			tx.Exec(fmt.Sprintf(`CREATE SEQUENCE "%s_nodes_audit_id_seq" INCREMENT 1 MINVALUE 0 MAXVALUE 2147483647 START 1 CACHE 1`, prefix))
			tx.Exec(fmt.Sprintf(`CREATE TABLE "%s_nodes_audit" (
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
				"slug" CHARACTER VARYING( 256 ) COLLATE "default" NOT NULL,
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

			err := tx.Commit()

			if err != nil {
				helper.SendWithHttpCode(res, http.StatusInternalServerError, err.Error())
			} else {
				helper.SendWithHttpCode(res, http.StatusOK, "Successfully create tables!")
			}
		})

		return nil
	})

	l.Prepare(func(app *goapp.App) error {
		if !conf.Test {
			return nil
		}

		mux := app.Get("goji.mux").(*web.Mux)

		prefix := ""

		mux.Put(prefix+"/data/purge", func(res http.ResponseWriter, req *http.Request) {

			manager := app.Get("gonode.manager").(*core.PgNodeManager)

			prefix := conf.Databases["master"].Prefix

			tx, _ := manager.Db.Begin()
			manager.Db.Exec(fmt.Sprintf(`DELETE FROM "%s_nodes"`, prefix))
			manager.Db.Exec(fmt.Sprintf(`DELETE FROM "%s_nodes_audit"`, prefix))
			err := tx.Commit()

			if err != nil {
				helper.SendWithStatus("KO", err.Error(), res)
			} else {
				helper.SendWithStatus("OK", "Data purged!", res)
			}
		})

		mux.Put(prefix+"/data/load", func(res http.ResponseWriter, req *http.Request) {
			manager := app.Get("gonode.manager").(*core.PgNodeManager)
			nodes := manager.FindBy(manager.SelectBuilder(), 0, 10)

			if nodes.Len() != 0 {
				helper.SendWithStatus("KO", "Table contains data, purge the data first!", res)

				return
			}

			err := fixtures.LoadFixtures(manager, 100)

			if err != nil {
				helper.SendWithStatus("KO", err.Error(), res)
			} else {
				helper.SendWithStatus("OK", "Data loaded!", res)
			}
		})

		return nil
	})
}
