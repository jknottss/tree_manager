package trees

import (
	"tree_manager/internal/storage"
)

func toOrganizationDb(org *Organization) *storage.OrganizationDb {
	return &storage.OrganizationDb{
		Id:          &org.Id,
		Type:        &org.Type,
		ParentId:    &org.ParentId,
		CountryCode: &org.CountryCode,
		TimeZone:    &org.TimeZone,
		PowerIds:    &org.PowerIds,
		Status:      &org.Status,
		Latitude:    &org.Latitude,
		Longitude:   &org.Longitude,
		Name:        &org.Name,
		ShortName:   &org.ShortName,
		TypeName:    &org.TypeName,
		LangCode:    &org.LangCode,
	}
}

func toOrganization(db *storage.OrganizationDbParse) *Organization {
	return &Organization{
		Id:          db.Id.Int64,
		Type:        db.Type.Int64,
		ParentId:    db.ParentId.Int64,
		CountryCode: db.CountryCode.String,
		TimeZone:    db.TimeZone.String,
		PowerIds:    db.PowerIds.String,
		Status:      db.Status.String,
		Latitude:    db.Latitude.String,
		Longitude:   db.Longitude.String,
		Name:        db.Name.String,
		ShortName:   db.ShortName.String,
		TypeName:    db.TypeName.String,
		LangCode:    db.LangCode.Int64,
	}
}

func (t *TreeManager) addEquipToOrgMap(orgs []*Organization) (EquipMap, error) {
	equipMap := EquipMap{}
	equips, err := t.GetAllEquipment()
	if err != nil {
		t.Logger.Error().Err(err).Msg("error while get all equipments")
		return nil, err
	}

	for _, equip := range equips {
		if equip.OrgParentId != nil {
			equipMap[*equip.OrgParentId] = append(equipMap[*equip.OrgParentId], equip)
		}
	}

	for _, node := range orgs {
		node.Equipments = equipMap[node.Id]
	}
	return equipMap, nil
}

func convertOrganizations(db []*storage.OrganizationDbParse) []*Organization {
	result := make([]*Organization, 0, len(db))
	for _, val := range db {
		result = append(result, toOrganization(val))
	}
	return result
}

func getOrganizationsMap(orgs []*Organization) OrgMap {
	result := OrgMap{}
	for _, node := range orgs {
		if node.ParentId != 0 {
			result[node.ParentId] = append(result[node.ParentId], node)
		} else {
			result[rootParentId] = append(result[rootParentId], node)
		}
	}
	return result
}

func buildOrgTree(orgs []*Organization, childs OrgMap) []*Organization {

	for _, node := range orgs {
		node.Organizations = childs[node.Id]
	}
	return childs[rootParentId]
}

func (t *TreeManager) GetAllOrganizations() ([]*Organization, error) {
	orgsFromDb, err := t.Repo.GetAllOrganizations()
	if err != nil {
		return nil, err
	}

	orgs := convertOrganizations(orgsFromDb)

	// addEquipToOrgMap - set equipments to Organization
	_, err = t.addEquipToOrgMap(orgs)
	if err != nil {
		return nil, err
	}

	//set organizations map
	orgsMap := getOrganizationsMap(orgs)

	//full result slice, and set childs in struct. Return roots only
	result := buildOrgTree(orgs, orgsMap)

	return result, nil

}

func (t *TreeManager) GetOrganizationById(id int64) ([]*Organization, error) {
	target, err := t.Repo.GetOrganizationById(id)
	if err != nil {
		return nil, err
	}
	orgsFromDb, err := t.Repo.GetAllOrganizations()
	if err != nil {
		return nil, err
	}
	orgs := convertOrganizations(orgsFromDb)

	equipMap, err := t.addEquipToOrgMap(orgs)
	if err != nil {
		return nil, err
	}

	childs := getOrganizationsMap(orgs)

	// buildOrgTree - need for building tree, set childs
	_ = buildOrgTree(orgs, childs)

	result := convertOrganizations(target)
	for _, node := range result {
		node.Organizations = childs[id]
		node.Equipments = equipMap[id]
	}

	return result, nil
}

func (t *TreeManager) CreateOrganization(org *Organization) (int64, error) {
	orgDb := toOrganizationDb(org)
	id, err := t.Repo.CreateOrganization(orgDb)
	return id, err
}

func (t *TreeManager) UpdateOrganization(org *Organization) error {
	orgDb := toOrganizationDb(org)
	err := t.Repo.UpdateOrganization(orgDb)
	return err

}

func (t *TreeManager) DeleteOrganization(id int64) error {
	err := t.Repo.DeleteOrganization(id)
	return err
}
