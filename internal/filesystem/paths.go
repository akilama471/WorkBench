package filesystem

import (
	"path/filepath"
)

type Paths struct {
	root string
}

func NewPaths(root string) *Paths {
	return &Paths{root: filepath.Clean(root)}
}

func (p *Paths) Root() string {
	return p.root
}

func (p *Paths) Bin() string {
	return filepath.Join(p.root, "bin")
}

func (p *Paths) Active() string {
	return filepath.Join(p.root, "active")
}

func (p *Paths) WWW() string {
	return filepath.Join(p.root, "www")
}

func (p *Paths) Data() string {
	return filepath.Join(p.root, "data")
}

func (p *Paths) Etc() string {
	return filepath.Join(p.root, "etc")
}

func (p *Paths) Logs() string {
	return filepath.Join(p.root, "logs")
}

func (p *Paths) Backup() string {
	return filepath.Join(p.root, "backup")
}

func (p *Paths) Packages() string {
	return filepath.Join(p.root, "packages")
}

func (p *Paths) Cache() string {
	return filepath.Join(p.root, "cache")
}

func (p *Paths) Database() string {
	return filepath.Join(p.root, "devbox.db")
}

func (p *Paths) ApacheBin(version string) string {
	return filepath.Join(p.Bin(), "apache", version)
}

func (p *Paths) PHPBin(version string) string {
	return filepath.Join(p.Bin(), "php", version)
}

func (p *Paths) MariaDBBin(version string) string {
	return filepath.Join(p.Bin(), "mariadb", version)
}

func (p *Paths) ActivePHP() string {
	return filepath.Join(p.Active(), "php")
}

func (p *Paths) ApacheConfig() string {
	return filepath.Join(p.Etc(), "apache")
}

func (p *Paths) PHPConfig() string {
	return filepath.Join(p.Etc(), "php")
}

func (p *Paths) MariaDBConfig() string {
	return filepath.Join(p.Etc(), "mariadb")
}

func (p *Paths) MariaDBData() string {
	return filepath.Join(p.Data(), "mariadb")
}

func (p *Paths) ApacheLogs() string {
	return filepath.Join(p.Logs(), "apache")
}

func (p *Paths) PHPLogs() string {
	return filepath.Join(p.Logs(), "php")
}

func (p *Paths) MariaDBLogs() string {
	return filepath.Join(p.Logs(), "mariadb")
}

func (p *Paths) CacheDownloads() string {
	return filepath.Join(p.Cache(), "downloads")
}

func (p *Paths) CacheExtract() string {
	return filepath.Join(p.Cache(), "extract")
}
