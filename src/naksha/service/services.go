package service

import (
	"naksha/db"
)

type Repository struct {
	user_service      UserService
	table_service     TableService
	table_row_service TableRowService
	map_service       MapService
	export_service    ExportService
}

func (r *Repository) GetUserService() UserService {
	return r.user_service
}

func (r *Repository) GetTableService() TableService {
	return r.table_service
}

func (r *Repository) GetTableRowService() TableRowService {
	return r.table_row_service
}

func (r *Repository) GetMapService() MapService {
	return r.map_service
}

func (r *Repository) GetExportService() ExportService {
	return r.export_service
}

func GetServicesRepository(map_user_dao db.Dao, api_user_dao db.Dao) Repository {
	user_service := MakeUserService(map_user_dao)
	table_service := MakeTableService(map_user_dao)
	table_row_service := MakeTableRowService(map_user_dao)
	map_service := MakeMapService(map_user_dao, api_user_dao)
	export_service := MakeExportService(map_user_dao)
	repository := Repository{
		user_service, table_service, table_row_service, map_service, export_service,
	}

	return repository
}
