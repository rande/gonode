

SQL requirements
================

Nodes Tables

```sql

    CREATE SEQUENCE "public"."nodes_id_seq" INCREMENT 1 MINVALUE 0 MAXVALUE 2147483647 START 0 CACHE 1;

    CREATE TABLE "public"."nodes" ( 
        "id" INTEGER DEFAULT nextval('nodes_id_seq'::regclass) NOT NULL UNIQUE, 
        "uuid" UUid NOT NULL, 
        "type" CHARACTER VARYING( 64 ) COLLATE "pg_catalog"."default" NOT NULL, 
        "name" CHARACTER VARYING( 2044 ) COLLATE "pg_catalog"."default" DEFAULT ''::CHARACTER VARYING NOT NULL, 
        "enabled" BOOLEAN DEFAULT 'true' NOT NULL, 
        "current" BOOLEAN DEFAULT 'false' NOT NULL, 
        "revision" INTEGER DEFAULT '1' NOT NULL,
        "status" INTEGER DEFAULT '0' NOT NULL, 
        "deleted" BOOLEAN DEFAULT 'false' NOT NULL, 
        "data" jsonb DEFAULT '{}'::jsonb NOT NULL, 
        "meta" jsonb DEFAULT '{}'::jsonb NOT NULL, 
        "slug" CHARACTER VARYING( 256 ) COLLATE "pg_catalog"."default" NOT NULL, 
        "source" UUid,
        "set_uuid" UUid, 
        "parent_uuid" UUid, 
        "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL, 
        "created_by" UUid NOT NULL, 
        "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL, 
        "updated_by" UUid NOT NULL, 
     PRIMARY KEY ( "id" )
    , CONSTRAINT "slug" UNIQUE( "parent_uuid","slug","revision" ), CONSTRAINT "uuid" UNIQUE( "revision","uuid" ) );
    
    CREATE INDEX "uuid_idx" ON "public"."nodes" USING btree( "uuid" ASC NULLS LAST );
    CREATE INDEX "uuid_current_idx" ON "public"."nodes" USING btree( "uuid" ASC NULLS LAST, "current" ASC NULLS LAST );
    -- CREATE INDEX "nodes_expr_idx" ON "public"."nodes" USING btree( "(data ->> 'Login'::text)" ASC NULLS LAST );
    -- CREATE INDEX "node_login_index" ON "public"."nodes" USING btree( "(data ->> 'Login'::text)" ASC NULLS LAST );    
    -- CREATE INDEX "node_media_size_index" ON "public"."nodes" USING btree( "(data ->> 'Width'::text)" ASC NULLS LAST, "(data ->> 'Height'::text)" ASC NULLS LAST );


Audit Log node table

    ```sql
    CREATE SEQUENCE "public"."nodes_audit_id_seq" INCREMENT 1 MINVALUE 0 MAXVALUE 2147483647 START 0 CACHE 1;
    CREATE TABLE "public"."nodes_audit" ( 
        "id" INTEGER DEFAULT nextval('nodes_id_seq'::regclass) NOT NULL UNIQUE, 
        "uuid" UUid NOT NULL, 
        "type" CHARACTER VARYING( 64 ) COLLATE "pg_catalog"."default" NOT NULL, 
        "name" CHARACTER VARYING( 2044 ) COLLATE "pg_catalog"."default" DEFAULT ''::CHARACTER VARYING NOT NULL, 
        "enabled" BOOLEAN DEFAULT 'true' NOT NULL, 
        "current" BOOLEAN DEFAULT 'false' NOT NULL, 
        "revision" INTEGER DEFAULT '1' NOT NULL, 
        "status" INTEGER DEFAULT '0' NOT NULL,
        "deleted" BOOLEAN DEFAULT 'false' NOT NULL, 
        "data" jsonb DEFAULT '{}'::jsonb NOT NULL, 
        "meta" jsonb DEFAULT '{}'::jsonb NOT NULL, 
        "slug" CHARACTER VARYING( 256 ) COLLATE "pg_catalog"."default" NOT NULL, 
        "source" UUid,
        "set_uuid" UUid, 
        "parent_uuid" UUid, 
        "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL, 
        "created_by" UUid NOT NULL, 
        "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL, 
        "updated_by" UUid NOT NULL, 
     PRIMARY KEY ( "id" )
    );
    
    
{"Type":"core.user","Name":"The user 12","Slug":"the-user-12","Data":{"Name":"","Login":"user12","Password":"{plain}user12"},"Meta":{},"Status":0,"Weight":0,"Revision":1,"CreatedAt":"2014-11-22T00:28:33.141819Z","UpdatedAt":"2014-11-22T00:28:33.141822Z","Enabled":false,"Deleted":false,"Parents":null,"Uuid":"b50b7a12-9529-4ca0-af3c-a84d54d47083","UpdatedBy":"00000000-0000-1000-0000-000000000000","CreatedBy":"00000000-0000-1000-0000-000000000000","ParentUuid":"11111111-1111-1111-1111-111111111111","SetUuid":"11111111-1111-1111-1111-111111111111","Source":"11111111-1111-1111-1111-111111111111"}