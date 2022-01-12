package models

import (
	"encoding/json"
	"strconv"
	"strings"
)

func stringToBool(s string) bool {
	var result bool

	switch s = strings.ToLower(s); {
	case s == "yes", s == "true":
		result = true
	case s == "no", s == "false":
		result = false
	default:
		result = false
	}

	return result
}

type NCServerInfo struct {
	Ocs Ocs `json:"ocs"`
}

type NCError struct {
	Ocs OcsWithError `json:"ocs"`
}

type Ocs struct {
	Meta Meta `json:"meta"`
	Data Data `json:"data"`
}

type OcsWithError struct {
	Meta Meta `json:"meta"`
}

type Meta struct {
	Status     string `json:"status"`
	StatusCode uint64 `json:"statuscode"`
	Message    string `json:"message"`
}

type Data struct {
	NextCloud   NextCloud   `json:"nextcloud"`
	Server      Server      `json:"server"`
	ActiveUsers ActiveUsers `json:"activeUsers"`
}

type NextCloud struct {
	System  System  `json:"system"`
	Storage Storage `json:"storage"`
	Shares  Shares  `json:"shares"`
}

type System struct {
	Version             string  `json:"version" metric:"nc_version"`
	Theme               string  `json:"theme"`
	EnableAvatars       bool    `json:"enable_avatars" metric:"avatars_enabled"`
	EnablePreviews      bool    `json:"enable_previews" metric:"previews_enabled"`
	MemcacheLocal       string  `json:"memcache.local" metric:"memcache_type" label:"local"`
	MemcacheDistributed string  `json:"memcache.distributed" metric:"memcache_type" label:"distributed"`
	FileLockingEnabled  bool    `json:"filelocking.enabled" metric:"file_locking_enabled"`
	MemcacheLocking     string  `json:"memcahe.locking" metric:"memcache_locking_type"`
	Debug               bool    `json:"debug" metric:"debug_mode_enabled"`
	FreeSpace           float64 `json:"freespace" metric:"free_space_bytes"`
	CPULoad             CPULoad `json:"cpuload"`
	MemTotal            float64 `json:"mem_total"`
	MemFree             float64 `json:"mem_free"`
	SwapTotal           float64 `json:"swap_total"`
	SwapFree            float64 `json:"swap_free"`
	Apps                Apps    `json:"apps"`
}

// Raw system data returned from API.
// Requires mapping string "yes"/"no" to bool
type intermediateSystem struct {
	Version             string     `json:"version"`
	Theme               string     `json:"theme"`
	EnableAvatars       string     `json:"enable_avatars"`
	EnablePreviews      string     `json:"enable_previews"`
	MemcacheLocal       string     `json:"memcache.local"`
	MemcacheDistributed string     `json:"memcache.distributed"`
	FileLockingEnabled  string     `json:"filelocking.enabled"`
	MemcacheLocking     string     `json:"memcahe.locking"`
	Debug               string     `json:"debug"`
	FreeSpace           float64    `json:"freespace"`
	CPULoad             [3]float64 `json:"cpuload"`
	MemTotal            float64    `json:"mem_total"`
	MemFree             float64    `json:"mem_free"`
	SwapTotal           float64    `json:"swap_total"`
	SwapFree            float64    `json:"swap_free"`
	Apps                Apps       `json:"apps"`
}

func (sys *System) UnmarshalJSON(data []byte) error {
	var inter intermediateSystem
	err := json.Unmarshal(data, &inter)
	if err != nil {
		return err
	}

	sys.Version = inter.Version
	sys.Theme = inter.Theme
	sys.EnableAvatars = stringToBool(inter.EnableAvatars)
	sys.EnablePreviews = stringToBool(inter.EnablePreviews)
	sys.MemcacheLocal = inter.MemcacheLocal
	sys.MemcacheDistributed = inter.MemcacheDistributed
	sys.FileLockingEnabled = stringToBool(inter.FileLockingEnabled)
	sys.MemcacheLocking = inter.MemcacheLocking
	sys.Debug = stringToBool(inter.Debug)
	sys.FreeSpace = inter.FreeSpace
	sys.CPULoad = CPULoad{
		OneMinuteAverage:     inter.CPULoad[0],
		FiveMinuteAverage:    inter.CPULoad[1],
		FifteenMinuteAverage: inter.CPULoad[2],
	}
	sys.MemTotal = inter.MemTotal
	sys.MemFree = inter.MemFree
	sys.SwapTotal = inter.SwapTotal
	sys.SwapFree = inter.SwapFree
	sys.Apps = inter.Apps

	return nil
}

type CPULoad struct {
	OneMinuteAverage     float64
	FiveMinuteAverage    float64
	FifteenMinuteAverage float64
}

type Apps struct {
	NumInstalled        float64     `json:"num_installed" metric:"installed_apps"`
	NumUpdatesAvailable float64     `json:"num_updates_available" metric:"app_updates_available"`
	AppUpdates          interface{} `json:"app_updates"`
}

type Storage struct {
	NumUsers         float64 `json:"num_users" metric:"users"`
	NumFiles         float64 `json:"num_files" metric:"files"`
	NumStorages      float64 `json:"num_storages"`
	NumStoragesLocal float64 `json:"num_storages_local" metric:"storages" label:"local"`
	NumStoragesHome  float64 `json:"num_storages_home" metric:"storages" label:"home"`
	NumStoragesOther float64 `json:"num_storages_other" metric:"storages" label:"other"`
}

type Shares struct {
	NumShares               float64 `json:"num_shares"`
	NumSharesUser           float64 `json:"num_shares_user" metric:"shares" label:"user"`
	NumSharesGroups         float64 `json:"num_shares_groups" metric:"shares" label:"groups"`
	NumSharesLink           float64 `json:"num_shares_link" metric:"shares" label:"link"`
	NumSharesMail           float64 `json:"num_shares_mail" metric:"shares" label:"mail"`
	NumSharesRoom           float64 `json:"num_shares_room" metric:"shares" label:"room"`
	NumSharesLinkNoPassword float64 `json:"num_shares_link_no_password" metric:"shares" label:"no_password"`
	NumFedSharesSent        float64 `json:"num_fed_shares_sent" metric:"fed_shares_sent"`
	NumFedSharesReceived    float64 `json:"num_fed_shares_received" metric:"fed_shares_received"`
}

type Server struct {
	WebServer string   `json:"webserver" metric:"web_server_type"`
	PHP       PHP      `json:"php"`
	Database  Database `json:"database"`
}

type PHP struct {
	Version           string  `json:"version" metric:"php_version"`
	MemoryLimit       float64 `json:"memory_limit" metric:"php_memory_limit_bytes"`
	MaxExecutionTime  float64 `json:"max_execution_time" metric:"php_max_execution_time_seconds"`
	UploadMaxFileSize float64 `json:"upload_max_filesize" metric:"php_upload_max_file_size_bytes"`
	Opcache           Opcache `json:"opcache"`
	APCU              APCU    `json:"apcu"`
}

type Opcache struct {
	OpcacheEnabled       bool                 `json:"opcache_enabled" metric:"php_opcache_enabled"`
	CacheFull            bool                 `json:"cache_full" metric:"php_opcache_full"`
	RestartPending       bool                 `json:"restart_pending" metric:"php_opcache_restart_pending"`
	RestartInProgress    bool                 `json:"restart_in_progress" metric:"php_opcache_restart_in_progress"`
	MemoryUsage          MemoryUsage          `json:"memory_usage"`
	InternedStringsUsage InternedStringsUsage `json:"interned_strings_usage"`
	OpcacheStatistics    OpcacheStatistics    `json:"opcache_statistics"`
	JIT                  JIT                  `json:"jiit"`
}

type MemoryUsage struct {
	UsedMemory              float64 `json:"used_memory" metric:"php_opcache_memory_used_bytes"`
	FreeMemory              float64 `json:"free_memory" metric:"php_opcache_memory_free_bytes"`
	WastedMemory            float64 `json:"waster_memory"`
	CurrentWastedPercentage float64 `json:"current_wasted_percentage"`
}

type InternedStringsUsage struct {
	BufferSize      float64 `json:"buffer_size" metric:"php_opcache_interned_strings_buffer_size_bytes"`
	UsedMemory      float64 `json:"used_memory" metric:"php_opcache_interned_strings_memory_used_bytes"`
	FreeMemory      float64 `json:"free_memory" metric:"php_opcache_interned_strings_memory_free_bytes"`
	NumberOfStrings float64 `json:"number_of_strings" metric:"php_opcache_interned_strings_count"`
}

type OpcacheStatistics struct {
	NumCachedScripts   float64 `json:"num_cached_scripts" metric:"php_opcache_cached_scripts_count"`
	NumCachedKeys      float64 `json:"num_cached_keys" metric:"php_opcache_cached_keys_count"`
	MaxCachedKeys      float64 `json:"max_cached_keys"`
	Hits               float64 `json:"hits" metric:"php_opcache_hits_count"`
	StartTime          float64 `json:"start_time" metric:"php_opcache_start_time_ticks"`
	LastRestartTime    float64 `json:"last_restart_time" metric:"php_opcache_last_restart_time_ticks"`
	OOMRestarts        float64 `json:"oom_restarts" metric:"php_opcache_restart_count" label:"oom"`
	HashRestarts       float64 `json:"hash_restarts" metric:"php_opcache_restart_count" label:"hash"`
	ManualRestarts     float64 `json:"manual_restarts" metric:"php_opcache_restart_count" label:"manual"`
	Misses             float64 `json:"misses" metric:"php_opcache_misses_count"`
	BlacklistMisses    float64 `json:"blacklist_misses" metric:"php_opcache_blacklist_misses_count"`
	BlacklistMissRatio float64 `json:"blacklist_miss_ratio"`
	OpcacheHitRate     float64 `json:"opcache_hit_rate"`
}

type JIT struct {
	Enabled    bool    `json:"enabled" metric:"php_jit_enabled"`
	On         bool    `json:"on" metric:"php_jit_on"`
	Kind       float64 `json:"kind" metric:"php_jit_kind"`
	OptLevel   float64 `json:"opt_level" metric:"php_jit_optimization_level"`
	OptFlags   float64 `json:"opt_flags" metric:"php_jit_optimization_flags"`
	BufferSize float64 `json:"buffer_size" metric:"php_jit_buffer_size_bytes"`
	BufferFree float64 `json:"buffer_free" metric:"php_jit_buffer_free_bytes"`
}

type APCU struct {
	Cache Cache `json:"cache"`
	SMA   SMA   `json:"sma"`
}

type Cache struct {
	NumSlots   float64 `json:"num_slots" metric:"php_apcu_cache_slots"`
	TTL        float64 `json:"ttl" metric:"php_apcu_cache_ttl"`
	NumHits    float64 `json:"num_hits" metric:"php_apcu_cache_hits_count"`
	NumMisses  float64 `json:"num_misses" metric:"php_apcu_cache_misses_count"`
	NumInserts float64 `json:"num_inserts" metric:"php_apcu_cache_inserts_count"`
	NumEntries float64 `json:"num_entries" metric:"php_apcu_cache_entries"`
	Expunges   float64 `json:"expunges" metric:"php_apcu_cache_expunges_count"`
	StartTime  float64 `json:"start_time" metric:"php_apcu_cache_start_time_ticks"`
	MemSize    float64 `json:"mem_size" metric:"php_apcu_cache_memory_free_bytes"`
	MemoryType string  `json:"memory_type" metric:"php_apcu_cache_memory_type"`
}

type SMA struct {
	NumSeg   float64 `json:"num_seg" metric:"php_apcu_sma_seg"`
	SegSize  float64 `json:"seg_size" metric:"php_apcu_sma_seg_size_bytes"`
	AvailMem float64 `json:"avail_mem" metric:"php_apcu_sma_memory_free_bytes"`
}

type Database struct {
	Type    string  `json:"type" metric:"database_type"`
	Version string  `json:"version" metric:"database_version"`
	Size    float64 `json:"size" metric:"database_size_bytes"`
}

// Raw database data.
// Required to convert size to float64
type intermediateDatabase struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Size    string `json:"size"`
}

func (db *Database) UnmarshalJSON(data []byte) error {
	var inter intermediateDatabase
	err := json.Unmarshal(data, &inter)
	if err != nil {
		return nil
	}

	size, err := strconv.ParseFloat(inter.Size, 64)
	if err != nil {
		return nil
	}

	db.Type = inter.Type
	db.Version = inter.Version
	db.Size = size

	return nil
}

type ActiveUsers struct {
	Last5Minutes float64 `json:"last5minutes" metric:"active_users" label:"5min"`
	Last1Hour    float64 `json:"last1hour" metric:"active_users" label:"60min"`
	Last24Hours  float64 `json:"last24hours" metric:"active_users" label:"1440min"`
}
