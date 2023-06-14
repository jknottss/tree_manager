package trees

const (
	rootParentId = -1
)

type EquipMap map[int64][]*Equipment

type OrgMap map[int64][]*Organization

type Equipment struct {
	StationId        int64  `json:"station_id"`
	UnitId           int64  `json:"unit_id"`
	UnitShortKey     string `json:"unit_short_key"`
	UnitType         string `json:"unit_type"`
	AnalyticsEnabled bool   `json:"analytics_enabled"`
	ParamId          string `json:"param_id"`
	StopParamId      string `json:"stop_param_id"`
	InfoId           string `json:"info_id"`
	MonIds           string `json:"mon_ids"`
	UnitMonitored    bool   `json:"unit_monitored"`
	StopCondition    string `json:"stop_condition"`
	OperationModes   string `json:"operation_modes"`
	MarkId           int64  `json:"mark_id"`
	UniqueUnitId     int64  `json:"unique_unit_id"`
	ExternalId       int64  `json:"external_id"`
	ParentType       string `json:"parent_type"`
	AgrParentId      *int64 `json:"agr_parent_id"`
	OrgParentId      *int64 `json:"org_parent_id"`
	Equipments       []*Equipment
}

type Organization struct {
	Id            int64  `json:"organization_id"`
	Type          int64  `json:"organization_type"`
	ParentId      int64  `json:"parent_id"`
	CountryCode   string `json:"country_code"`
	RegionCode    string `json:"region_code"`
	TimeZone      string `json:"time_zone"`
	PowerIds      string `json:"power_ids"`
	Status        string `json:"status"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
	Name          string `json:"organization_name"`
	ShortName     string `json:"organization_short_name"`
	LangCode      int64  `json:"lang_code"`
	TypeName      string `json:"type_name"`
	Organizations []*Organization
	Equipments    []*Equipment
}
