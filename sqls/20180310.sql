CREATE TABLE public.mstr_user (
  id SERIAL4 NOT NULL,
  name VARCHAR(64) NOT NULL,
  username VARCHAR(32) NOT NULL,
  password VARCHAR(96) NOT NULL,
  schema_name VARCHAR(32) NOT NULL,
  google_maps_key VARCHAR(128) DEFAULT NULL,
  bing_maps_key VARCHAR(128) DEFAULT NULL,
  login_attempts SMALLINT DEFAULT 0,
  last_login_time TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
  last_login_ip VARCHAR(48) DEFAULT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(id),
  CONSTRAINT user_uq_username UNIQUE(username),
  CONSTRAINT user_uq_schema_name UNIQUE(schema_name)
);

CREATE TABLE public.mstr_table (
  id SERIAL4 NOT NULL,
  user_id INTEGER NOT NULL,
  name VARCHAR(96) NOT NULL,
  schema_name VARCHAR(64) NOT NULL,
  table_name VARCHAR(64) NOT NULL,
  status SMALLINT NOT NULL,
  api_access BOOLEAN DEFAULT '0',
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(id),
  CONSTRAINT table_uq_schema_table UNIQUE(schema_name, table_name),
  CONSTRAINT table_fk_user FOREIGN KEY(user_id) REFERENCES mstr_user(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE public.mstr_layer (
  id SERIAL4 NOT NULL,
  table_id INTEGER NOT NULL,
  geometry_column VARCHAR(64) NOT NULL,
  query TEXT NOT NULL,
  infowindow TEXT NOT NULL,
  style TEXT NOT NULL,
  options TEXT,
  hash VARCHAR(64) NOT NULL,
  update_hash VARCHAR(64) NOT NULL,
  geometry_type VARCHAR(16) NOT NULL DEFAULT 'geometry',
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(id),
  CONSTRAINT layer_fk_table FOREIGN KEY(table_id) REFERENCES mstr_table(id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT layer_uq_hash UNIQUE(hash)
);

CREATE TABLE public.mstr_map (
  id         SERIAL4 NOT NULL,
  user_id    INTEGER NOT NULL,
  name       VARCHAR(64) NOT NULL,
  hash       VARCHAR(64) NOT NULL,
  base_layer VARCHAR(32) NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
  PRIMARY KEY(id),
  CONSTRAINT map_fk_user FOREIGN KEY(user_id) REFERENCES mstr_user(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE public.mstr_map_layer (
  map_id      INTEGER NOT NULL,
  layer_id    INTEGER NOT NULL,
  layer_index SMALLINT NOT NULL,
  created_at  TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
  updated_at  TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
  PRIMARY KEY(map_id, layer_id),
  CONSTRAINT map_layer_fk_map FOREIGN KEY(map_id) REFERENCES mstr_map(id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT map_layer_fk_layer FOREIGN KEY(layer_id) REFERENCES mstr_layer(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE public.mstr_export (
  id SERIAL4 NOT NULL,
  user_id INTEGER NOT NULL,
  table_id INTEGER NOT NULL,
  status INTEGER NOT NULL,
  filename VARCHAR(96) NOT NULL,
  hash VARCHAR(128) NOT NULL,
  extension   VARCHAR(16) NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE,
  updated_at TIMESTAMP WITHOUT TIME ZONE,
  PRIMARY KEY(id),
  UNIQUE(hash),
  CONSTRAINT export_fk_user FOREIGN KEY(user_id) REFERENCES mstr_user(id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT export_fk_table FOREIGN KEY(table_id) REFERENCES mstr_table(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE OR REPLACE FUNCTION public.naksha_update_geom_webmercator() RETURNS trigger AS $$
  BEGIN
  NEW.the_geom_webmercator = ST_Transform(NEW.the_geom, 3857);
    RETURN NEW;
  END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION public.naksha_update_updated_at() RETURNS trigger AS $$
  BEGIN
  NEW.updated_at = NOW();
    RETURN NEW;
  END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION public.naksha_prepare_table(id INTEGER, table_name VARCHAR(64), hash VARCHAR(64)) RETURNS void AS $$
DECLARE
  col_type VARCHAR;
  style VARCHAR;
  geometry_column VARCHAR;
  query VARCHAR;
  naksha_geom_type VARCHAR;
  infowindow VARCHAR;
BEGIN
  EXECUTE 'ALTER TABLE ' || table_name::regclass || ' DROP COLUMN IF EXISTS the_geom_webmercator';
  EXECUTE 'ALTER TABLE ' || table_name::regclass || ' ADD COLUMN the_geom_webmercator Geometry(Geometry, 3857)';

  EXECUTE 'UPDATE ' || table_name::regclass || ' SET the_geom_webmercator = ST_Transform(the_geom, 3857)';

  EXECUTE 'ALTER TABLE ' || table_name::regclass || ' DROP COLUMN IF EXISTS created_at';
  EXECUTE 'ALTER TABLE ' || table_name::regclass || ' DROP COLUMN IF EXISTS updated_at';
  EXECUTE 'ALTER TABLE ' || table_name::regclass || ' ADD COLUMN created_at TIMESTAMP WITHOUT TIME ZONE';
  EXECUTE 'ALTER TABLE ' || table_name::regclass || ' ADD COLUMN updated_at TIMESTAMP WITHOUT TIME ZONE';

  EXECUTE 'CREATE TRIGGER naksha_trg_update_webmercator BEFORE INSERT OR UPDATE ON ' || table_name::regclass
	  || ' FOR EACH ROW EXECUTE PROCEDURE public.naksha_update_geom_webmercator()';

  EXECUTE 'CREATE TRIGGER naksha_trg_update_updated_at BEFORE UPDATE ON ' || table_name::regclass
	  || ' FOR EACH ROW EXECUTE PROCEDURE public.naksha_update_updated_at()';

  EXECUTE 'SELECT ST_GeometryType(the_geom) FROM ' || table_name::regclass || ' LIMIT 1' INTO col_type;
  col_type := UPPER(col_type);
  IF col_type = 'ST_MULTIPOLYGON' OR col_type = 'ST_POLYGON' THEN
    style := '<Rule><PolygonSymbolizer fill="#000000" fill-opacity="0.75" /><LineSymbolizer stroke="#ffffff" stroke-width="0.5" stroke-opacity="1.0" /></Rule>';
    naksha_geom_type := 'polygon';
  ELSIF col_type = 'ST_MULTILINESTRING' OR col_type = 'ST_LINESTRING' THEN
    style := '<Rule><LineSymbolizer stroke="#ffffff" stroke-width="4" stroke-opacity="1.0" /></Rule>';
    naksha_geom_type := 'linestring';
  ELSIF col_type = 'ST_MULTIPOINT' OR col_type = 'ST_POINT' THEN
    style := '<Rule><MarkersSymbolizer fill="#000000" stroke="#ffffff" opacity="0.75" stroke-width="1" stroke-opacity="1.0" width="10" height="10" marker-type="ellipse" /></Rule>';
    naksha_geom_type := 'point';
  ELSE
    style := '<Rule>'
          || '<PolygonSymbolizer fill="#000000" fill-opacity="0.75" />'
          || '<LineSymbolizer stroke="#ffffff" stroke-width="0.5" stroke-opacity="1.0" />'
          || '<MarkersSymbolizer fill="#000000" stroke="#ffffff" opacity="0.75" stroke-width="1" stroke-opacity="1.0" width="10" height="10" marker-type="ellipse" />'
          || '</Rule>';
    naksha_geom_type := 'unknown';
  END IF;

  geometry_column := 'the_geom_webmercator';
  query := 'SELECT * FROM ' || table_name::regclass;
  infowindow := '{"fields": []}';

  INSERT INTO public.mstr_layer (table_id, hash, geometry_column, query, style, geometry_type, infowindow, update_hash)
         VALUES (id, hash, geometry_column, query, style, naksha_geom_type, infowindow, hash);
  
  EXECUTE 'UPDATE ' || table_name::regclass || ' SET created_at = current_timestamp, updated_at = current_timestamp';
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION public.naksha_prepare_csv(schema_name VARCHAR(64), table_name VARCHAR(64)) RETURNS VOID AS $$
DECLARE
  cnt INTEGER;
  max_id INTEGER;
  seq_name VARCHAR;
  table_col_name VARCHAR;
  schema_table VARCHAR;
BEGIN
  schema_table := schema_name || '.' || table_name;
  EXECUTE 'SELECT COUNT(*) AS cnt FROM information_schema.columns WHERE table_schema = ''' || schema_name || ''' AND table_name = ''' || table_name || ''' AND column_name = ''naksha_id''' INTO cnt;
  IF cnt = 0 THEN
    EXECUTE 'ALTER TABLE ' || schema_table::regclass || ' ADD COLUMN naksha_id SERIAL4 NOT NULL';
  ELSE
    EXECUTE 'SELECT MAX(naksha_id) FROM ' || schema_table::regclass INTO max_id;
    max_id = max_id + 1;
    seq_name := schema_table || '_naksha_id_seq';
    table_col_name := schema_table || '.naksha_id';
    EXECUTE 'DROP SEQUENCE IF EXISTS ' || seq_name;
    EXECUTE 'CREATE SEQUENCE ' || seq_name || ' START WITH ' || max_id || ' OWNED BY ' || table_col_name;
    EXECUTE 'ALTER TABLE ' || schema_table::regclass || ' ALTER COLUMN naksha_id SET DEFAULT nextval(''' || seq_name || '''::regclass)';
    EXECUTE 'ALTER TABLE ' || schema_table::regclass || ' ALTER COLUMN naksha_id SET NOT NULL';
  END IF;
  EXECUTE 'ALTER TABLE ' || schema_table::regclass || ' ADD PRIMARY KEY(naksha_id)';
END;
$$ LANGUAGE plpgsql;

