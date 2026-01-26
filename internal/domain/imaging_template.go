package domain

import (
	"time"

	"gorm.io/gorm"
)

// ImagingModality represents the type of imaging modality
type ImagingModality string

const (
	ImagingModalityXRay        ImagingModality = "XRAY"
	ImagingModalityCT          ImagingModality = "CT"
	ImagingModalityMRI         ImagingModality = "MRI"
	ImagingModalityUltrasound  ImagingModality = "ULTRASOUND"
	ImagingModalityMammography ImagingModality = "MAMMOGRAPHY"
	ImagingModalityFluoroscopy ImagingModality = "FLUOROSCOPY"
)

// BodyPart represents the body part being imaged
type BodyPart string

const (
	BodyPartHead      BodyPart = "HEAD"
	BodyPartChest     BodyPart = "CHEST"
	BodyPartAbdomen   BodyPart = "ABDOMEN"
	BodyPartPelvis    BodyPart = "PELVIS"
	BodyPartSpine     BodyPart = "SPINE"
	BodyPartExtremity BodyPart = "EXTREMITY"
	BodyPartOther     BodyPart = "OTHER"
)

// ImagingTemplate represents an imaging template
type ImagingTemplate struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Template Information
	Code            string          `gorm:"uniqueIndex;size:20;not null" json:"code"`
	Name            string          `gorm:"size:200;not null;index" json:"name"`
	Modality        ImagingModality `gorm:"size:20;not null;index" json:"modality"`
	BodyPart        BodyPart        `gorm:"size:20;not null" json:"body_part"`
	Description     string          `gorm:"type:text" json:"description"`
	TemplateContent string          `gorm:"type:text" json:"template_content"` // Default findings template
	Price           float64         `gorm:"type:decimal(10,2)" json:"price"`
	IsActive        bool            `gorm:"default:true;index" json:"is_active"`
}

// TableName specifies the table name for ImagingTemplate model
func (ImagingTemplate) TableName() string {
	return "imaging_templates"
}
