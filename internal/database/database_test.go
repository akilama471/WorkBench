package database

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDatabaseOpenAndMigrate(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() failed: %v", err)
	}
	defer db.Close()

	if err := db.Initialize(); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	info, err := os.Stat(dbPath)
	if err != nil {
		t.Fatalf("database file not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("database file is empty")
	}
}

func TestDatabaseSettings(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() failed: %v", err)
	}
	defer db.Close()

	if err := db.Initialize(); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	if err := db.SetSetting("key1", "value1"); err != nil {
		t.Fatalf("SetSetting() failed: %v", err)
	}

	val, err := db.GetSetting("key1")
	if err != nil {
		t.Fatalf("GetSetting() failed: %v", err)
	}
	if val != "value1" {
		t.Errorf("GetSetting(\"key1\") = %q, want \"value1\"", val)
	}

	if err := db.SetSetting("key1", "updated"); err != nil {
		t.Fatalf("SetSetting() update failed: %v", err)
	}

	val, err = db.GetSetting("key1")
	if err != nil {
		t.Fatalf("GetSetting() after update failed: %v", err)
	}
	if val != "updated" {
		t.Errorf("GetSetting(\"key1\") after update = %q, want \"updated\"", val)
	}
}

func TestDatabaseGetSettingNotExist(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() failed: %v", err)
	}
	defer db.Close()

	if err := db.Initialize(); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	val, err := db.GetSetting("nonexistent")
	if err != nil {
		t.Fatalf("GetSetting(\"nonexistent\") failed: %v", err)
	}
	if val != "" {
		t.Errorf("GetSetting(\"nonexistent\") = %q, want empty", val)
	}
}

func TestDatabaseDeleteSetting(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() failed: %v", err)
	}
	defer db.Close()

	if err := db.Initialize(); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	_ = db.SetSetting("key1", "value1")

	if err := db.DeleteSetting("key1"); err != nil {
		t.Fatalf("DeleteSetting() failed: %v", err)
	}

	val, err := db.GetSetting("key1")
	if err != nil {
		t.Fatalf("GetSetting() after delete failed: %v", err)
	}
	if val != "" {
		t.Errorf("GetSetting(\"key1\") after delete = %q, want empty", val)
	}
}

func TestDatabaseGetAllSettings(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() failed: %v", err)
	}
	defer db.Close()

	if err := db.Initialize(); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	_ = db.SetSetting("a", "1")
	_ = db.SetSetting("b", "2")

	all, err := db.GetAllSettings()
	if err != nil {
		t.Fatalf("GetAllSettings() failed: %v", err)
	}

	if len(all) != 2 {
		t.Errorf("GetAllSettings() returned %d, want 2", len(all))
	}
	if all["a"] != "1" {
		t.Errorf("all[\"a\"] = %q, want \"1\"", all["a"])
	}
	if all["b"] != "2" {
		t.Errorf("all[\"b\"] = %q, want \"2\"", all["b"])
	}
}

func TestDatabaseIdempotentMigrations(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() failed: %v", err)
	}
	defer db.Close()

	if err := db.Initialize(); err != nil {
		t.Fatalf("first Initialize() failed: %v", err)
	}

	_ = db.SetSetting("persist", "data")

	if err := db.Initialize(); err != nil {
		t.Fatalf("second Initialize() failed: %v", err)
	}

	val, err := db.GetSetting("persist")
	if err != nil {
		t.Fatalf("GetSetting() after second init failed: %v", err)
	}
	if val != "data" {
		t.Errorf("data lost after second migration: got %q", val)
	}
}
