// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package orm

import (
"fmt"
"time"

"appsite-go/internal/core/setting"

"gorm.io/driver/mysql"
"gorm.io/driver/sqlite"
"gorm.io/gorm"
"gorm.io/gorm/logger"
"gorm.io/gorm/schema"
)

// InitDB initializes the database based on configuration
func InitDB(cfg *setting.DatabaseConfig) (*gorm.DB, error) {
switch cfg.Type {
case "mysql":
return NewMySQLConnection(cfg)
case "sqlite", "sqlite3":
return NewSQLiteConnection(cfg)
default:
return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
}
}

// NewMySQLConnection initializes a MySQL connection using GORM.
func NewMySQLConnection(cfg *setting.DatabaseConfig) (*gorm.DB, error) {
if cfg == nil {
return nil, fmt.Errorf("database config is nil")
}

dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
cfg.User,
cfg.Password,
cfg.Host,
cfg.Port,
cfg.Name,
cfg.Charset,
)

dialector := mysql.New(mysql.Config{
DSN:                       dsn,
DefaultStringSize:         256,
DisableDatetimePrecision:  true,
DontSupportRenameIndex:    true,
DontSupportRenameColumn:   true,
SkipInitializeWithVersion: false,
})

return NewConnection(dialector, cfg)
}

// NewSQLiteConnection initializes a SQLite connection using GORM.
func NewSQLiteConnection(cfg *setting.DatabaseConfig) (*gorm.DB, error) {
if cfg == nil {
return nil, fmt.Errorf("database config is nil")
}

// For SQLite, Host or Name can be used as the file path
dbPath := cfg.Name
if dbPath == "" {
dbPath = "appsite.db"
}

dialector := sqlite.Open(dbPath)
return NewConnection(dialector, cfg)
}

// NewConnection initializes GORM with a specific dialector.
func NewConnection(dialector gorm.Dialector, cfg *setting.DatabaseConfig) (*gorm.DB, error) {
gormConfig := &gorm.Config{
NamingStrategy: schema.NamingStrategy{
SingularTable: true,
},
Logger: logger.Default.LogMode(logger.Info),
}

db, err := gorm.Open(dialector, gormConfig)
if err != nil {
return nil, fmt.Errorf("failed to open database connection: %w", err)
}

sqlDB, err := db.DB()
if err != nil {
return nil, err
}

// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
sqlDB.SetMaxIdleConns(10)

// SetMaxOpenConns sets the maximum number of open connections to the database.
sqlDB.SetMaxOpenConns(100)

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
sqlDB.SetConnMaxLifetime(time.Hour)

return db, nil
}
