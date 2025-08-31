package clinics


type controlBody string

const(
  DHA controlBody = "DHA"
  HAAD controlBody = "HAAD"
  MOH controlBody = "MOH"

)


type Clinic struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `gorm:"type:varchar(255);not null" json:"name"`
	ControlBody controlBody `gorm:"type:enum('DHA','HAAD','MOH');not null" json:"control_body"`
	Phone     string `gorm:"type:varchar(20);not null" json:"phone"`
	Email    string `gorm:"type:varchar(100);not null" json:"email"`
	FacilityID string `gorm:"type:varchar(100);not null" json:"facility_id"`
	Licence string `gorm:"type:varchar(100);not null" json:"licence"`	
	CreatedAt string `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt string `gorm:"autoUpdateTime" json:"updated_at"`
}