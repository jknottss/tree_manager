package trees

import (
	"tree_manager/internal/storage"
)

func toEquipmentDb(equip *Equipment) *storage.EquipmentDb {
	return &storage.EquipmentDb{
		StationId:        &equip.StationId,
		UnitId:           &equip.UnitId,
		UnitShortKey:     &equip.UnitShortKey,
		UnitType:         &equip.UnitType,
		AnalyticsEnabled: &equip.AnalyticsEnabled,
		ParamId:          &equip.ParamId,
		StopParamId:      &equip.StopParamId,
		InfoId:           &equip.InfoId,
		MonIds:           &equip.MonIds,
		UnitMonitored:    &equip.UnitMonitored,
		StopCondition:    &equip.StopCondition,
		OperationModes:   &equip.OperationModes,
		MarkId:           &equip.MarkId,
		UniqueUnitId:     &equip.UniqueUnitId,
		ExternalId:       &equip.ExternalId,
		ParentType:       &equip.ParentType,
		AgrParentId:      equip.AgrParentId,
		OrgParentId:      equip.OrgParentId,
	}
}

func toEquipment(db *storage.EquipmentDbParse) *Equipment {
	result := &Equipment{
		StationId:        db.StationId.Int64,
		UnitId:           db.UnitId.Int64,
		UnitShortKey:     db.UnitShortKey.String,
		UnitType:         db.UnitType.String,
		AnalyticsEnabled: db.AnalyticsEnabled.Bool,
		ParamId:          db.ParamId.String,
		StopParamId:      db.StopParamId.String,
		InfoId:           db.InfoId.String,
		MonIds:           db.MonIds.String,
		UnitMonitored:    db.UnitMonitored.Bool,
		StopCondition:    db.StopCondition.String,
		OperationModes:   db.OperationModes.String,
		MarkId:           db.MarkId.Int64,
		UniqueUnitId:     db.UniqueUnitId.Int64,
		ExternalId:       db.ExternalId.Int64,
		ParentType:       db.ParentType.String,
	}
	if db.AgrParentId.Valid {
		result.AgrParentId = &db.AgrParentId.Int64
	}
	if db.OrgParentId.Valid {
		result.OrgParentId = &db.OrgParentId.Int64
	}
	return result
}

func getEquipmentsMap(equips []*Equipment) EquipMap {
	equipsMap := EquipMap{}
	for _, node := range equips {
		if node.AgrParentId != nil {
			equipsMap[*node.AgrParentId] = append(equipsMap[*node.AgrParentId], node)
		} else {
			equipsMap[rootParentId] = append(equipsMap[rootParentId], node)
		}
	}
	return equipsMap
}

func convertEquipments(db []*storage.EquipmentDbParse) []*Equipment {
	result := make([]*Equipment, 0, len(db))

	for _, equipDb := range db {
		result = append(result, toEquipment(equipDb))
	}
	return result
}

func buildAgrTree(equips []*Equipment, data EquipMap) []*Equipment {
	for _, node := range equips {
		node.Equipments = data[node.UniqueUnitId]
	}
	return data[rootParentId]
}

func (t *TreeManager) GetAllEquipment() ([]*Equipment, error) {
	equipsFromDb, err := t.Repo.GetAllEquipment()
	if err != nil {
		return nil, err
	}
	equip := convertEquipments(equipsFromDb)
	maps := getEquipmentsMap(equip)
	result := buildAgrTree(equip, maps)

	return result, nil
}

func (t *TreeManager) GetEquipmentById(id int64) ([]*Equipment, error) {
	target, err := t.Repo.GetEquipmentById(id)
	if err != nil {
		return nil, err
	}
	allEquips, err := t.Repo.GetAllEquipment()
	if err != nil {
		return nil, err
	}
	equips := convertEquipments(allEquips)

	equipMap := getEquipmentsMap(equips)

	_ = buildAgrTree(equips, equipMap)

	result := convertEquipments(target)
	for _, node := range result {
		node.Equipments = equipMap[node.UniqueUnitId]
	}
	return result, nil
}

func (t *TreeManager) CreateEquipment(equip *Equipment) (int64, error) {
	equipDbIn := toEquipmentDb(equip)
	id, err := t.Repo.CreateEquipment(equipDbIn)
	return id, err

}

func (t *TreeManager) UpdateEquipment(equip *Equipment) (*Equipment, error) {
	equipDbIn := toEquipmentDb(equip)
	equipDb, err := t.Repo.UpdateEquipment(equipDbIn)
	if err != nil {
		return nil, err
	}
	return toEquipment(equipDb), nil

}

func (t *TreeManager) DeleteEquipment(id int64) error {
	err := t.Repo.DeleteEquipment(id)
	return err
}
