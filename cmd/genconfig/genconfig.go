package main

import (
	"bufio"
	"fmt"
	"naksha/helper"
	"os"
	"os/user"
	"path"
	"strings"
)

type config_data struct {
	db_name       string
	db_admin_user string
	db_admin_pass string
	db_map_user   string
	db_map_pass   string
	db_api_user   string
	db_api_pass   string
	domain_name   string
}

var server_port = "31001"
var tiler_port = "31002"

var config = `
<?xml version="1.0" encoding="UTF-8"?>
<naksha>
  <database>
    <host>127.0.0.1</host>
    <name>{db_name}</name>
    <sslmode>disable</sslmode>
    <adminuser>{admin_user}</adminuser>
    <adminpass>{admin_pass}</adminpass>
    <user>{map_user}</user>
    <pass>{map_pass}</pass>
    <apiuser>{api_user}</apiuser>
    <apipass>{api_pass}</apipass>
  </database>
  <assets_dir>{assets_dir}</assets_dir>
  <views_dir>{views_dir}</views_dir>
  <tmp_dir>{tmp_dir}</tmp_dir>
  <host>127.0.0.1</host>
  <port>{server_port}</port>
  <max_request_body_size>8388608</max_request_body_size>
  <auth_key>{auth_key}</auth_key>
  <enc_key>{enc_key}</enc_key>
  <tiler_domain>{tiler_domain}</tiler_domain>
  <!-- trailing slash is required -->
  <mapnik_input_path>/usr/lib/mapnik/3.0/input/</mapnik_input_path>
</naksha>
`

var sql = `
create database {db_name};
\connect {db_name};
revoke all on schema public from public;
create extension postgis;
create user {admin_user} with password '{admin_pass}';
create user {map_user} with password '{map_pass}';
create user {api_user} with password '{api_pass}';
grant all privileges on database {db_name} to {admin_user};
grant all privileges on schema public to {admin_user};
grant usage on schema public to {map_user};
alter default privileges for user {admin_user} in schema public grant select, insert, update, delete, trigger on tables to {map_user};
alter default privileges for user {admin_user} in schema public grant usage, select, update on sequences to {map_user};
alter default privileges for user {admin_user} in schema public grant execute on functions to {map_user};
grant {map_user} to {admin_user};
`

var supervisor_conf = `
[program:naksha_server]
command={server_dir}/bin/server -config {server_dir}/config/config.xml
autostart=true
user={username}
stderr_logfile={server_dir}/tmp/logs/naksha.error
stdout_logfile={server_dir}/tmp/logs/naksha.out

[fcgi-program:naksha_tiler]
command={tiler_dir}/naksha_tiler.fcgi {server_dir}/config/config.xml
socket=tcp://localhost:{tiler_port}
autostart=true
user={username}
stderr_logfile={server_dir}/tmp/logs/tiler.error
stdout_logfile={server_dir}/tmp/logs/tiler.out
`

var apache_config = `
<VirtualHost *:80>
    ServerName {domain_name}
    DocumentRoot {server_dir}

    ProxyRequests Off
    ProxyPreserveHost On
    ProxyPass "/lyr/" "fcgi://127.0.0.1:{tiler_port}/"

    ProxyPass /assets/ !
    ProxyPass /lyr/ !
    ProxyPass "/" "http://127.0.0.1:{server_port}/"
    ProxyPassReverse "/" "http://127.0.0.1:{server_port}/"

    <Directory {server_dir}>
        Options +ExecCGI
        Allow from all
        Require all granted
    </Directory>

    RedirectMatch "/assets/([0-9]+)/(.*)" "/assets/$2"

    ErrorLog ${APACHE_LOG_DIR}/naksha-error.log
    LogLevel warn
    CustomLog ${APACHE_LOG_DIR}/naksha-access.log combined
</VirtualHost>
`

func main() {
	var tmp string

	cd := config_data{}
	scanner := bufio.NewScanner(os.Stdin)

	tmp = get_input(scanner, "Domain (or public ip address) from which application will be accessed")
	cd.domain_name = clean_domain(tmp)

	tmp = get_input_default(scanner, "Database name", "naksha_db")
	cd.db_name = tmp

	tmp = get_input_default(scanner, "DB admin user", "mn_admin_user")
	cd.db_admin_user = tmp

	tmp = get_input(scanner, "DB admin user password")
	cd.db_admin_pass = tmp

	tmp = get_input_default(scanner, "DB map user", "mn_map_user")
	cd.db_map_user = tmp

	tmp = get_input(scanner, "DB map user password")
	cd.db_map_pass = tmp

	tmp = get_input_default(scanner, "DB api user", "mn_api_user")
	cd.db_api_user = tmp

	tmp = get_input(scanner, "DB api user password")
	cd.db_api_pass = tmp

	fmt.Println("=============================================")
	fmt.Println("Database name = " + cd.db_name)
	fmt.Println("DB admin user = " + cd.db_admin_user)
	fmt.Println("DB admin pass = " + cd.db_admin_pass)
	fmt.Println("DB map user = " + cd.db_map_user)
	fmt.Println("DB map pass = " + cd.db_map_pass)
	fmt.Println("DB api user = " + cd.db_api_user)
	fmt.Println("DB api pass = " + cd.db_api_pass)
	fmt.Println("=============================================")

	errs := make([]string, 0)
	cwd, err := os.Getwd()
	if err != nil {
		errs = append(errs, fmt.Sprintf("%s", err))
	}
	user, err := user.Current()
	if err != nil {
		errs = append(errs, fmt.Sprintf("%s", err))
	}

	if len(errs) > 0 {
		fmt.Printf("Error(s):\n%v\n", strings.Join(errs, "\n"))
	} else {
		cwd = strings.TrimSpace(cwd)
		write_config_file(cd, cwd)
		write_sql_file(cd, cwd)
		write_supervisor_file(user, cwd)
		write_apache_file(user, cwd, cd.domain_name)
	}
}

func write_config_file(cd config_data, app_dir string) {
	str := config

	auth_key := helper.RandomKey(64)
	enc_key := helper.RandomKey(32)

	str = strings.Replace(str, "{db_name}", cd.db_name, 1)
	str = strings.Replace(str, "{admin_user}", cd.db_admin_user, 1)
	str = strings.Replace(str, "{admin_pass}", cd.db_admin_pass, 1)
	str = strings.Replace(str, "{map_user}", cd.db_map_user, 1)
	str = strings.Replace(str, "{map_pass}", cd.db_map_pass, 1)
	str = strings.Replace(str, "{api_user}", cd.db_api_user, 1)
	str = strings.Replace(str, "{api_pass}", cd.db_api_pass, 1)
	str = strings.Replace(str, "{assets_dir}", app_dir+"/assets", 1)
	str = strings.Replace(str, "{views_dir}", app_dir+"/views", 1)
	str = strings.Replace(str, "{tmp_dir}", app_dir+"/tmp", 1)
	str = strings.Replace(str, "{auth_key}", auth_key, 1)
	str = strings.Replace(str, "{enc_key}", enc_key, 1)
	str = strings.Replace(str, "{tiler_domain}", cd.domain_name, 1)
	str = strings.Replace(str, "{server_port}", server_port, -11)

	config_path := path.Join(app_dir, "config.xml")
	write_to_file(config_path, str)
	fmt.Println("config file " + config_path + " written")
}

func write_sql_file(cd config_data, app_dir string) {
	str := sql

	str = strings.Replace(str, "{db_name}", cd.db_name, -1)
	str = strings.Replace(str, "{admin_user}", cd.db_admin_user, -1)
	str = strings.Replace(str, "{admin_pass}", cd.db_admin_pass, -1)
	str = strings.Replace(str, "{map_user}", cd.db_map_user, -1)
	str = strings.Replace(str, "{map_pass}", cd.db_map_pass, -1)
	str = strings.Replace(str, "{api_user}", cd.db_api_user, -1)
	str = strings.Replace(str, "{api_pass}", cd.db_api_pass, -1)

	sql_file := path.Join(app_dir, "db.sql")
	write_to_file(sql_file, str)
	fmt.Println("sql file " + sql_file + " written")
}

func write_supervisor_file(usr *user.User, app_dir string) {
	tiler_path := path.Join(app_dir, "..", "tiler")
	str := supervisor_conf
	str = strings.Replace(str, "{server_dir}", app_dir, -1)
	str = strings.Replace(str, "{tiler_dir}", tiler_path, -1)
	str = strings.Replace(str, "{username}", usr.Username, -1)
	str = strings.Replace(str, "{tiler_port}", tiler_port, -1)
	supervisor_file := path.Join(app_dir, "supervisor.naksha.conf")
	write_to_file(supervisor_file, str)
	fmt.Println("supervisor config file " + supervisor_file + " written")
}

func write_apache_file(usr *user.User, app_dir string, domain_name string) {
	domain_parts := strings.Split(domain_name, "/")
	str := apache_config
	str = strings.Replace(str, "{server_dir}", app_dir, -1)
	// domain name is prefixed with http(s):// and suffixed with /
	// For apache vhost config, we don't need that suffix and prefix
	str = strings.Replace(str, "{domain_name}", domain_parts[2], -1)
	str = strings.Replace(str, "{server_port}", server_port, -1)
	str = strings.Replace(str, "{tiler_port}", tiler_port, -1)
	apache_file := path.Join(app_dir, "apache.vhost.conf")
	write_to_file(apache_file, str)
	fmt.Println("apache config file " + apache_file + " written")
}

func clean_domain(domain string) string {
	domain = strings.TrimSpace(domain)
	domain = strings.ToLower(domain)
	if strings.HasPrefix(domain, "http://") || strings.HasPrefix(domain, "https://") {
		// nothing to do
	} else {
		domain = "http://" + domain
	}
	domain = strings.TrimRight(domain, "/")
	domain = domain + "/"

	return domain
}

func write_to_file(filepath, contents string) {
	f, err := os.Create(filepath)
	defer f.Close()
	if err != nil {
		fmt.Println("Could not open file for writing. Error = ", err)
		os.Exit(1)
	}
	_, err = f.WriteString(contents)
	if err != nil {
		fmt.Println("Could not write to file. Error = ", err)
		os.Exit(1)
	}
	f.Sync()
}

func get_input_default(scanner *bufio.Scanner, input_name string, default_value string) string {
	input_text := ""

	for input_text == "" {
		fmt.Printf("Enter " + input_name + ". Default(" + default_value + "): ")
		scanner.Scan()
		input_text = scanner.Text()
		input_text = strings.TrimSpace(input_text)
		if len(input_text) == 0 {
			input_text = default_value
		}
	}

	return input_text
}

func get_input(scanner *bufio.Scanner, input_name string) string {
	input_text := ""

	for input_text == "" {
		fmt.Printf("Enter " + input_name + ": ")
		scanner.Scan()
		input_text = scanner.Text()
		input_text = strings.TrimSpace(input_text)
		if len(input_text) == 0 {
			fmt.Println("Invalid value")
		}
	}

	return input_text
}
