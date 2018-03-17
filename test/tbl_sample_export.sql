--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.6
-- Dumped by pg_dump version 9.5.6
-- -- Used for testing the export package

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

SET search_path = public, pg_catalog;

DROP TRIGGER IF EXISTS naksha_trg_update_webmercator ON public.tbl_sample_export;
DROP TRIGGER IF EXISTS naksha_trg_update_updated_at ON public.tbl_sample_export;
ALTER TABLE IF EXISTS ONLY public.tbl_sample_export DROP CONSTRAINT IF EXISTS tbl_sample_export_pkey;
ALTER TABLE IF EXISTS public.tbl_sample_export ALTER COLUMN naksha_id DROP DEFAULT;
DROP SEQUENCE IF EXISTS public.tbl_sample_export_naksha_id_seq;
DROP TABLE IF EXISTS public.tbl_sample_export;
SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: tbl_sample_export; Type: TABLE; Schema: public;
--

CREATE TABLE tbl_sample_export (
    naksha_id integer NOT NULL,
    the_geom geometry(Geometry,4326),
    the_geom_webmercator geometry(Geometry,3857),
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    name character varying
);


--
-- Name: tbl_sample_export_naksha_id_seq; Type: SEQUENCE; Schema: public;
--

CREATE SEQUENCE tbl_sample_export_naksha_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tbl_sample_export_naksha_id_seq; Type: SEQUENCE OWNED BY; Schema: public;
--

ALTER SEQUENCE tbl_sample_export_naksha_id_seq OWNED BY tbl_sample_export.naksha_id;


--
-- Name: naksha_id; Type: DEFAULT; Schema: public;
--

ALTER TABLE ONLY tbl_sample_export ALTER COLUMN naksha_id SET DEFAULT nextval('tbl_sample_export_naksha_id_seq'::regclass);


--
-- Data for Name: tbl_sample_export; Type: TABLE DATA; Schema: public;
--

COPY tbl_sample_export (naksha_id, the_geom, the_geom_webmercator, created_at, updated_at, name) FROM stdin;
1	0105000020E61000000100000001020000000900000001000060D3F552406AF920F36CD03640010000602B16534094C7D5B4C3133740010000604B325340886108B9CC3237400100006073555340CC92F6BD4C6137400100006093715340A78FF0A467C8374001000060D37C53407DF6EB543758384000000060AB865340358A8ADF54B4384001000060BB945340358482E2A5293940010000609BA553406989610614803940	0105000020110F000001000000010200000009000000B19CBF0D5B1A604164DAA412CCE84341087FA342D33560419BAD86F9E62644411D5AFA38B64D6041B2E7EEED90434441F8EBE6EC916B6041D73E8BDC8F6E44410CC73DE3748360412DABE6B51BCE4441168560DF028D6041A10FCD7FDF5345415AEBFE1B5F956041ECBD055DDDA94541E7582A9750A160414A479530BE174641F2755E91A5AF60418E96BC0FF5684641	2017-10-14 16:52:50	2017-10-14 11:29:50.711717	Road 1
2	0105000020E61000000100000001020000000600000001000060BBC153402AA6CE61AA4E354001000060A3D953404A986AEBA5A7354001000060A3D953409033B1668834364001000060A3D95340A43157D6F48C36400100006023F05340DFE78718FADF36400000006093F853400FDB7231F8373740	0105000020110F0000010000000102000000060000000751B58788C760412E639FA146874241D8247F3FD6DB604177B2C6807DD84241D8247F3FD6DB6041E9B79B4C7A594341D8247F3FD6DB60412F07C32BB1AA4341EBA0C437F2EE604171F7D80C21F7434171AFDEB41CF66041B64600EC57484441	2017-10-14 16:53:12	2017-10-14 11:29:53.705869	Road 2
3	0105000020E610000001000000010200000004000000010000608B9753404E774F377A99314001000060FB7253408637156EEF09324001000060735553403D3970F01C7A3240010000609B4B5340C2986B6F0F3F3340	0105000020110F00000100000001020000000400000068083316B4A3604199E6F3F645603E41CE1EC2A2A6846041487DCDA5EB283F41F8EBE6EC916B6041F513A75491F13F41B08548B0356360419145D6628BA94041	2017-10-14 16:53:28	2017-10-14 11:29:56.373833	Road 3
4	0105000020E61000000100000001020000000500000001000060D36E52404BA696B0112834400100006003995240FC4026A5DCAB344001000060B3AC52402BCDE05E391A354001000060C3BA5240862E25CD309D35400100006053DF5240A92ACCD8F9143540	0105000020110F00000100000001020000000500000099683D78664F5F41499AD20DBF7B4141DCF9415B0F975F41B4E184DD2DF34141FA92BB4D80B85F4108ADF1B480574241086E124463D05F4173F4A384EFCE4241A1207A153F076041044EE0B6B9524241	2017-10-14 16:53:44	2017-10-14 11:29:59.557499	Road 4
5	0105000020E6100000010000000102000000050000000100006073715540AF32A79A458C344001000060C38A55406EF52CFFF6EA3440000000607B785540D8CF7BC654733540010000602B5F554030F31DCE72D1354000000060C330554033FFAEAB5B493640	0105000020110F000001000000010200000005000000848D6A33313662419BA71CE983D6414117B9B8AAB04B6241E15555C6812C42414A4400F1293C624150FC1894B7A84241B718B279AA26624198AA5171B5FE4241D5C8A2C940FF6141F933E144966C4341	2017-10-14 16:54:09	2017-10-14 11:30:02.514497	Road 5
\.


--
-- Name: tbl_sample_export_naksha_id_seq; Type: SEQUENCE SET; Schema: public;
--

SELECT pg_catalog.setval('tbl_sample_export_naksha_id_seq', 5, true);


--
-- Name: tbl_sample_export_pkey; Type: CONSTRAINT; Schema: public;
--

ALTER TABLE ONLY tbl_sample_export
    ADD CONSTRAINT tbl_sample_export_pkey PRIMARY KEY (naksha_id);


--
-- Name: naksha_trg_update_updated_at; Type: TRIGGER; Schema: public;
--

CREATE TRIGGER naksha_trg_update_updated_at BEFORE UPDATE ON tbl_sample_export FOR EACH ROW EXECUTE PROCEDURE naksha_update_updated_at();


--
-- Name: naksha_trg_update_webmercator; Type: TRIGGER; Schema: public;
--

CREATE TRIGGER naksha_trg_update_webmercator BEFORE INSERT OR UPDATE ON tbl_sample_export FOR EACH ROW EXECUTE PROCEDURE naksha_update_geom_webmercator();


--
-- PostgreSQL database dump complete
--
