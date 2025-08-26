package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"-"`
	Fullname   string `json:"fullname" form:"fullname" gorm:"not null"`
	Username   string `json:"username" form:"username" gorm:"not null, unique"`
	Password   string `json:"password" form:"password" gorm:"not null"`
	Role       string `json:"role" form:"role" gorm:"not null"`
}

type Project struct {
	gorm.Model  `json:"-"`
	ID          uint64 `json:"id" form:"id" gorm:"not null"`
	Name        string `json:"name" form:"name" gorm:"not null"`
	Description string `json:"description" form:"description" gorm:"not null"`
}

type ItemType string

const (
	ItemTypeFolder ItemType = "folder"
	ItemTypePdf    ItemType = "pdf"
	ItemTypeJpg    ItemType = "jpg"
	ItemTypePng    ItemType = "png"
	ItemTypeOther  ItemType = "other" // <<-- untuk semua jenis file selain 3 di atas
)


type Item struct {
	gorm.Model `json:"-"`
	ID         uint64   `json:"id" form:"id" gorm:"not null"`
	Name       string   `json:"name" form:"name" gorm:"type:varchar(255);not null"`
	ParentID   *uint    `json:"parentId" form:"parentId" gorm:"index"`               // Boleh NULL (untuk root)
	Parent     *Item    `json:"parent" form:"parent" gorm:"foreignKey:ParentID"` // Relasi ke folder induk
	Type       ItemType `json:"type" form:"type" gorm:"type:varchar(20);not null"`
	FilePath   *string  `json:"path" form:"path" gorm:"type:text"`         // NULL jika folder
	MimeType   *string  `json:"mimeType" form:"mimeType" gorm:"type:varchar(100)"` // NULL jika folder
	Size       *int64   `json:"size" form:"size"` // NULL jika folder
	UploadedBy *uint    `json:"uploadedBy" form:"uploadedBy" gorm:"index"` // FK ke tabel users (optional)
	IsFolder   bool     `json:"isFolder" form:"isFolder"`
	IsLink     bool     `json:"isLink" form:"isLink"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// type Shift struct {
// 	gorm.Model
// 	Name        string `json:"name" form:"name" gorm:"type:varchar(255);not null"`
// 	Description string `json:"description" form:"description" gorm:"type:text;not null"`
// }

// type CS struct {
// 	gorm.Model
// 	Name string `json:"name" form:"name" gorm:"type:text;not null"`
// 	NIK  string `json:"nik" form:"nik" gorm:"type:text;not null"`
// }

// type CSS struct {
// 	gorm.Model
// 	Name string `json:"name" form:"name" gorm:"type:text;not null"`
// 	NIK  string `json:"nik" form:"nik" gorm:"type:text;not null"`
// }

// type Gerbang struct {
// 	gorm.Model
// 	NamaGerbang string `json:"nama_gerbang" form:"nama_gerbang" gorm:"type:text;not null"`
// }

// type PicLimNTK struct {
// 	gorm.Model
// 	CS1      string `json:"cs1" form:"cs1" gorm:"type:text;not null"`
// 	CS1NIK   string `json:"cs1_nik" form:"cs1_nik" gorm:"type:text;not null"`
// 	CS2      string `json:"cs2" form:"cs2" gorm:"type:text;not null"`
// 	CS2NIK   string `json:"cs2_nik" form:"cs2_nik" gorm:"type:text;not null"`
// 	CSS      string `json:"css" form:"css" gorm:"type:text;not null"`
// 	CSSNIK   string `json:"css_nik" form:"css_nik" gorm:"type:text;not null"`
// 	Date     string `json:"date" form:"date" gorm:"type:text;not null"`
// 	Shift    string `json:"shift" form:"shift" gorm:"type:text;not null"`
// 	CS1Gardu string `json:"cs1_gardu" form:"cs1_gardu" gorm:"type:text;not null"`
// 	CS2Gardu string `json:"cs2_gardu" form:"cs2_gardu" gorm:"type:text;not null"`

// 	BeritaAcara []BeritaAcaraLimNTK `gorm:"foreignKey:IDPicLimNTK;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
// }

// type PicTNIPolri struct {
// 	gorm.Model
// 	CS     string `json:"cs" form:"cs" gorm:"type:text;not null"`
// 	CSNIK  string `json:"cs_nik" form:"cs_nik" gorm:"type:text;not null"`
// 	CSS    string `json:"css" form:"css" gorm:"type:text;not null"`
// 	CSSNIK string `json:"css_nik" form:"css_nik" gorm:"type:text;not null"`
// 	Date   string `json:"date" form:"date" gorm:"type:text;not null"`
// 	Shift  string `json:"shift" form:"shift" gorm:"type:text;not null"`

// 	BeritaAcara []BeritaAcaraTNIPolri `gorm:"foreignKey:IDPicTNIPolri;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
// }

// type BeritaAcaraLimNTK struct {
// 	gorm.Model
// 	Date          string `json:"date" form:"date" gorm:"type:text;not null"`
// 	Gardu         string `json:"gardu" form:"gardu" gorm:"type:text;not null"`
// 	Resi          string `json:"resi" form:"resi" gorm:"type:text;not null"`
// 	JumlahNotran  string `json:"jumlah_notran" form:"jumlah_notran" gorm:"type:text;not null"`
// 	JenisKejadian string `json:"jenis_kejadian" form:"jenis_kejadian" gorm:"type:text;not null"`
// 	Gerbang       string `json:"gerbang" form:"gerbang" gorm:"type:text;not null"`
// 	Nopol         string `json:"nopol" form:"nopol" gorm:"type:text;not null"`
// 	EtollKTP      string `json:"etoll_ktp" form:"etoll_ktp" gorm:"type:text;not null"`
// 	AsalGerbang   string `json:"asal_gerbang" form:"asal_gerbang" gorm:"type:text;not null"`
// 	Uraian        string `json:"uraian" form:"uraian" gorm:"type:text;not null"`

// 	IDPicLimNTK uint64 `json:"id_pic_lim_ntk" form:"id_pic_lim_ntk" gorm:"not null"`
// }

// type BeritaAcaraTNIPolri struct {
// 	gorm.Model
// 	Date              string `json:"date" form:"date" gorm:"type:text;not null"`
// 	Gardu             string `json:"gardu" form:"gardu" gorm:"type:text;not null"`
// 	Pukul             string `json:"pukul" form:"pukul" gorm:"type:text;not null"`
// 	NomorResi         string `json:"nomor_resi" form:"nomor_resi" gorm:"type:text;not null"`
// 	AsalTujuanGerbang string `json:"asal_tujuan_gerbang" form:"asal_tujuan_gerbang" gorm:"type:text;not null"`
// 	JumlahKendaraan   string `json:"jumlah_kendaraan" form:"jumlah_kendaraan" gorm:"type:text;not null"`
// 	NopolKendaraan    string `json:"nopol_kendaraan" form:"nopol_kendaraan" gorm:"type:text;not null"`
// 	Instansi          string `json:"instansi" form:"instansi" gorm:"type:text;not null"`
// 	PenanggungJawab   string `json:"penanggung_jawab" form:"penanggung_jawab" gorm:"type:text;not null"`
// 	JenisInstansi     string `json:"jenis_instansi" form:"jenis_instansi" gorm:"type:text;not null"`

// 	IDPicTNIPolri uint64 `json:"id_pic_tni_polri" form:"id_pic_tni_polri" gorm:"not null"`
// }

type ResponseError struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type ResponseSuccess struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
