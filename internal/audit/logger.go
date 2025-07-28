package audit

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func Log(ctx context.Context, db *gorm.DB, entity string, entityId uint, action string, changes map[string][2]interface{})error{

 useid:= GetUserid(ctx)

	changesJSON, err := json.Marshal(changes)
	if err != nil {
		return err
	}

	logEntry:= Auditlog{
		Entity:    entity,	
		EntityID:  entityId,
		Action:    action,
		ChangedBy: useid,
		Changes:   datatypes.JSON(changesJSON),
		CreatedAt: time.Now(),


	}
	return  db.Create(&logEntry).Error
}