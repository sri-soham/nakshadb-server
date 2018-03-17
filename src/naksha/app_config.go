package naksha

import (
	"encoding/xml"
	"os"
)

type config struct {
	DB                 database `xml:"database"`
	AssetsDir          string   `xml:"assets_dir"`
	ViewsDir           string   `xml:"views_dir"`
	TmpDir             string   `xml:"tmp_dir"`
	Port               string   `xml:"port"`
	AuthKey            string   `xml:"auth_key"`
	EncKey             string   `xml:"enc_key"`
	Host               string   `xml:"host"`
	MaxRequestBodySize int64    `xml:"max_request_body_size"`
	TilerDomain        string   `xml:"tiler_domain"`
}

type database struct {
	Host      string `xml:"host"`
	Name      string `xml:"name"`
	Sslmode   string `xml:"sslmode"`
	AdminUser string `xml:"adminuser"`
	AdminPass string `xml:"adminpass"`
	User      string `xml:"user"`
	Pass      string `xml:"pass"`
	ApiUser   string `xml:"apiuser"`
	ApiPass   string `xml:"apipass"`
}

type AppConfig struct {
	db_host               string
	db_name               string
	db_ssl_mode           string
	db_admin_user         string
	db_admin_pass         string
	db_user               string
	db_pass               string
	db_api_user           string
	db_api_pass           string
	assets_dir            string
	views_dir             string
	tmp_dir               string
	host                  string
	port                  string
	auth_key              []byte
	enc_key               []byte
	tiler_domain          string
	max_request_body_size int64
}

func (ac *AppConfig) DBHost() string {
	return ac.db_host
}

func (ac *AppConfig) DBName() string {
	return ac.db_name
}

func (ac *AppConfig) DBAdminUser() string {
	return ac.db_admin_user
}

func (ac *AppConfig) DBAdminPass() string {
	return ac.db_admin_pass
}

func (ac *AppConfig) DBUser() string {
	return ac.db_user
}

func (ac *AppConfig) DBPass() string {
	return ac.db_pass
}

func (ac *AppConfig) DBSslMode() string {
	return ac.db_ssl_mode
}

func (ac *AppConfig) DBApiUser() string {
	return ac.db_api_user
}

func (ac *AppConfig) DBApiPass() string {
	return ac.db_api_pass
}

func (ac *AppConfig) Port() string {
	return ac.port
}

func (ac *AppConfig) Host() string {
	return ac.host
}

func (ac *AppConfig) HostPort() string {
	return ac.host + ":" + ac.port
}

func (ac *AppConfig) AuthKey() []byte {
	return ac.auth_key
}

func (ac *AppConfig) EncKey() []byte {
	return ac.enc_key
}

func (ac *AppConfig) MaxRequestBodySize() int64 {
	return ac.max_request_body_size
}

func (ac *AppConfig) TemplatesDir() string {
	return ac.views_dir
}

func (ac *AppConfig) AssetsDir() string {
	return ac.assets_dir
}

func (ac *AppConfig) LogsDir() string {
	return ac.tmp_dir + "/logs"
}

func (ac *AppConfig) SessionsDir() string {
	return ac.tmp_dir + "/sessions"
}

func (ac *AppConfig) ImportsDir() string {
	return ac.tmp_dir + "/imports"
}

func (ac *AppConfig) ExportsDir() string {
	return ac.tmp_dir + "/exports"
}

func (ac *AppConfig) TilerDomain() string {
	return ac.tiler_domain
}

func MakeConfig(config_file_path string) (AppConfig, error) {
	file, err := os.Open(config_file_path)
	if err != nil {
		return AppConfig{}, err
	}

	defer file.Close()
	var c config

	decoder := xml.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		return AppConfig{}, err
	}

	app_config := AppConfig{
		db_host:       c.DB.Host,
		db_name:       c.DB.Name,
		db_ssl_mode:   c.DB.Sslmode,
		db_admin_user: c.DB.AdminUser,
		db_admin_pass: c.DB.AdminPass,
		db_user:       c.DB.User,
		db_pass:       c.DB.Pass,
		db_api_user:   c.DB.ApiUser,
		db_api_pass:   c.DB.ApiPass,
		views_dir:     c.ViewsDir,
		assets_dir:    c.AssetsDir,
		tmp_dir:       c.TmpDir,
		port:          c.Port,
		host:          c.Host,
		max_request_body_size: c.MaxRequestBodySize,
		auth_key:              []byte(c.AuthKey),
		enc_key:               []byte(c.EncKey),
		tiler_domain:          c.TilerDomain,
	}

	return app_config, nil
}

func MakeAppConfigTest(test_values map[string]interface{}) AppConfig {
	return AppConfig{}
}
