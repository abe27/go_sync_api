package models

import "time"

type ResponseOrder struct {
	Data []Order `json:"data"`
}

type User struct {
	ID        string    `gorm:"primaryKey;size:21;" json:"id,omitempty"`
	UserName  string    `gorm:"not null;column:username;index;unique;size:10" json:"user_name,omitempty" form:"user_name"`
	Email     string    `gorm:"not null;unique;size:50;" json:"email,omitempty" form:"email"`
	Password  string    `gorm:"not null;unique;size:60;" json:"-" form:"password"`
	IsActive  bool      `gorm:"null" json:"is_active,omitempty" form:"is_active" default:"false"`
	CreatedAt time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt time.Time `json:"updated_at,omitempty" default:"now"`
}

type FileEdi struct {
	ID         string    `gorm:"primaryKey;size:21" json:"id,omitempty"`
	FactoryID  *string   `gorm:"not null;" json:"factory_id,omitempty" form:"factory_id"`
	MailboxID  *string   `gorm:"not null;" json:"mailbox_id,omitempty" form:"mailbox_id" binding:"required"`
	FileTypeID *string   `gorm:"not null;" json:"file_type_id,omitempty" form:"file_type_id"`
	BatchNo    string    `gorm:"not null;unique;size:10" json:"batch_no,omitempty" form:"batch_no" binding:"required"`
	Size       int64     `json:"size,omitempty" form:"size"`
	BatchName  string    `gorm:"size:50" json:"batch_name,omitempty" form:"batch_name"`
	CreationOn time.Time `json:"creation_on,omitempty" form:"creation_on"`
	Flags      string    `gorm:"size:5" json:"flags,omitempty" form:"flags" binding:"required"`
	FormatType string    `gorm:"size:5" json:"format_type,omitempty" form:"format_type" binding:"required"`
	Originator string    `gorm:"size:10" json:"originator,omitempty" form:"originator" binding:"required"`
	BatchPath  string    `gorm:"size:255" json:"batch_path,omitempty"`
	IsDownload bool      `json:"is_download,omitempty" form:"is_download" binding:"required"`
	IsActive   bool      `json:"is_active,omitempty" form:"is_active" binding:"required"`
	CreatedAt  time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt  time.Time `json:"updated_at,omitempty" default:"now"`
	Factory    Factory   `gorm:"foreignKey:FactoryID;references:ID" json:"factory,omitempty"`
	Mailbox    Mailbox   `gorm:"foreignKey:MailboxID;references:ID" json:"mailbox,omitempty"`
	FileType   FileType  `gorm:"foreignKey:FileTypeID;references:ID" json:"file_type,omitempty"`
}

type Part struct {
	ID          string    `gorm:"primaryKey;size:21" json:"id,omitempty"`
	Slug        string    `gorm:"size:50;unique;not null;" json:"slug,omitempty" form:"slug" binding:"required"`
	Title       string    `gorm:"size:50;unique;not null;" json:"title,omitempty" form:"title" binding:"required"`
	Description string    `json:"description,omitempty" form:"description" binding:"required"`
	IsActive    bool      `json:"is_active,omitempty" form:"is_active" default:"true"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
}

type Ledger struct {
	ID          string    `gorm:"primaryKey;size:21" json:"id,omitempty"`
	WhsID       *string   `gorm:"not null;" json:"whs_id,omitempty" form:"whs_id" binding:"required"`
	FactoryID   *string   `gorm:"not null;" json:"factory_id,omitempty" form:"factory_id" binding:"required"`
	PartID      *string   `gorm:"not null;" json:"part_id,omitempty" form:"part_id" binding:"required"`
	PartTypeID  *string   `gorm:"not null;" json:"part_type_id,omitempty" form:"part_type_id" binding:"required"`
	UnitID      *string   `gorm:"not null;" json:"unit_id,omitempty" form:"unit_id" binding:"required"`
	DimWidth    float64   `json:"dim_width,omitempty" form:"dim_width" default:"0"`
	DimLength   float64   `json:"dim_length,omitempty" form:"dim_length" default:"0"`
	DimHeight   float64   `json:"dim_height,omitempty" form:"dim_height" default:"0"`
	GrossWeight float64   `json:"gross_weight,omitempty" form:"gross_weight" default:"0"`
	NetWeight   float64   `json:"net_weight,omitempty" form:"net_weight" default:"0"`
	Qty         float64   `json:"qty,omitempty" form:"qty" default:"0"`
	Ctn         float64   `json:"ctn,omitempty" form:"ctn" default:"0"`
	IsActive    bool      `json:"is_active,omitempty" form:"is_active" default:"true"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
	Whs         Whs       `gorm:"foreignKey:WhsID;references:ID" json:"whs,omitempty"`
	Factory     Factory   `gorm:"foreignKey:FactoryID;references:ID" json:"factory,omitempty"`
	Part        Part      `gorm:"foreignKey:PartID;references:ID;" json:"part,omitempty"`
	PartType    PartType  `gorm:"foreignKey:PartTypeID;references:ID" json:"part_type,omitempty"`
	Unit        Unit      `gorm:"foreignKey:UnitID;references:ID" json:"unit,omitempty"`
}

type OrderZone struct {
	ID          string    `gorm:"primaryKey;size:21" json:"id,omitempty"`
	Value       int64     `gorm:"not null;" json:"value,omitempty" form:"value" binding:"required"`
	FactoryID   *string   `gorm:"not null;" json:"factory_id,omitempty" form:"factory_id" binding:"required"`
	WhsID       *string   `gorm:"not null;" json:"whs_id,omitempty" form:"whs_id" binding:"required"`
	Description string    `json:"description,omitempty" form:"description" binding:"required"`
	IsActive    bool      `json:"is_active,omitempty" form:"is_active" binding:"required"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
	Factory     Factory   `gorm:"foreignKey:FactoryID;references:ID" json:"factory,omitempty"`
	Whs         Whs       `gorm:"foreignKey:WhsID;references:ID" json:"whs,omitempty"`
}

type Order struct {
	ID           string         `gorm:"primaryKey;unique;index;size:21" json:"id"`
	RowID        string         `gorm:"null;size:18" json:"row_id" form:"row_id"`
	ConsigneeID  *string        `gorm:"not null" json:"consignee_id" form:"consignee_id" binding:"required"`
	ShipmentID   *string        `gorm:"not null" json:"shipment_id" form:"shipment_id" binding:"required"`
	EtdDate      *time.Time     `gorm:"not null;type:date;" json:"etd_date" form:"etd_date" binding:"required"`
	PcID         *string        `gorm:"not null" json:"pc_id" form:"pc_id" binding:"required"`
	CommercialID *string        `gorm:"not null" json:"commercial_id" form:"commercial_id" binding:"required"`
	SampleFlgID  *string        `gorm:"not null" json:"sample_flg_id" form:"sample_flg_id" binding:"required"`
	OrderTitleID *string        `gorm:"not null" json:"title_id" form:"title_id" binding:"required"`
	Bioabt       int64          `json:"bioabt" form:"bioabt" binding:"required"`
	CarrierCode  string         `gorm:"size:255;" json:"carrier_code" form:"carrier_code" binding:"required"`
	ShipForm     string         `gorm:"size:255;" json:"ship_form" form:"ship_form" binding:"required"`
	ShipTo       string         `gorm:"size:255;" json:"ship_to" form:"ship_to" binding:"required"`
	ShipVia      string         `gorm:"size:255;" json:"ship_via" form:"ship_via" binding:"required"`
	ShipDer      string         `gorm:"size:255;" json:"ship_der" form:"ship_der" binding:"required"`
	LoadingArea  string         `gorm:"size:255;" json:"loading_area" form:"loading_area" binding:"required"`
	Privilege    string         `gorm:"size:50;" json:"privilege" form:"privilege" binding:"required"`
	ZoneCode     string         `gorm:"size:10;unique;" json:"zone_code" form:"zone_code" binding:"required"`
	RunningSeq   int64          `json:"running_seq" form:"running_seq" binding:"required"`
	IsChecked    bool           `json:"is_checked" form:"is_checked"`
	IsInvoice    bool           `json:"is_invoice" form:"is_invoice"`
	IsSync       bool           `json:"is_sync" form:"is_sync"`
	IsActive     bool           `json:"is_active" form:"is_active"`
	CreatedAt    time.Time      `json:"created_at" default:"now"`
	UpdatedAt    time.Time      `json:"updated_at" default:"now"`
	Consignee    Consignee      `gorm:"foreignKey:ConsigneeID;references:ID" json:"consignee"`
	Shipment     Shipment       `gorm:"foreignKey:ShipmentID;references:ID" json:"shipment"`
	Pc           Pc             `gorm:"foreignKey:PcID;references:ID" json:"pc"`
	Commercial   Commercial     `gorm:"foreignKey:CommercialID;references:ID" json:"commercial"`
	SampleFlg    SampleFlg      `gorm:"foreignKey:SampleFlgID;references:ID" json:"sample_flg"`
	OrderTitle   OrderTitle     `gorm:"foreignKey:OrderTitleID;references:ID" json:"order_title"`
	OrderDetail  []*OrderDetail `json:"order_detail"`
	Pallet       []*Pallet      `json:"pallet"`
}

type OrderDetail struct {
	ID            string    `gorm:"primaryKey;unique;index;size:21" json:"id"`
	RowID         string    `gorm:"null;size:18" json:"row_id" form:"row_id"`
	OrderID       *string   `gorm:"not null;" json:"order_id" form:"order_id" binding:"required"`
	Pono          *string   `gorm:"not null;size:25" json:"pono" form:"pono" binding:"required"`
	LedgerID      *string   `gorm:"not null;" json:"ledger_id" form:"ledger_id" binding:"required"`
	OrderPlanID   *string   `gorm:"not null;unique;" json:"order_plan_id" form:"order_plan_id" binding:"required"`
	OrderCtn      int64     `json:"order_ctn" form:"order_ctn" binding:"required"`
	TotalOnPallet int64     `json:"total_on_pallet" form:"total_on_pallet" binding:"required"`
	IsChecked     bool      `json:"is_checked" form:"is_checked" default:"false"`
	IsMatched     bool      `json:"is_matched" form:"is_matched" default:"false"`
	IsSync        bool      `json:"is_sync" form:"is_sync" default:"false"`
	IsActive      bool      `json:"is_active" form:"is_active"`
	CreatedAt     time.Time `json:"created_at" default:"now"`
	UpdatedAt     time.Time `json:"updated_at" default:"now"`
	Order         Order     `gorm:"foreignKey:OrderID;references:ID" json:"order"`
	Ledger        Ledger    `gorm:"foreignKey:LedgerID;references:ID" json:"ledger"`
	OrderPlan     OrderPlan `gorm:"foreignKey:OrderPlanID;references:ID" json:"orderplan"`
}

type SyncOrderPlan struct {
	Message string      `json:"message"`
	Data    []OrderPlan `json:"data"`
}

type OrderPlan struct {
	ID               string        `gorm:"primaryKey;size:21;" json:"id"`
	RowID            string        `gorm:"null;size:18" json:"row_id,omitempty" form:"row_id"`
	FileEdiID        *string       `gorm:"not null;" json:"file_edi_id" form:"file_edi_id"`
	WhsID            *string       `gorm:"not null;" json:"whs_id" form:"whs_id"`
	OrderZoneID      *string       `gorm:"not null;" json:"order_zone_id" form:"order_type_id" binding:"required"`
	ConsigneeID      *string       `gorm:"not null;" json:"consignee_id" form:"consignee_id"`
	ReviseOrderID    *string       `gorm:"null;" json:"revise_order_id" form:"revise_order_id" binding:"required"`
	LedgerID         *string       `gorm:"not null;" json:"ledger_id" form:"ledger_id" binding:"required"`
	PcID             *string       `gorm:"not null;" json:"pc_id" form:"pc_id" binding:"required"`
	CommercialID     *string       `gorm:"not null;" json:"commercial_id" form:"commercial_id" binding:"required"`
	OrderTypeID      *string       `gorm:"not null;" json:"order_type_id" form:"order_type_id" binding:"required"`
	ShipmentID       *string       `gorm:"not null;" json:"shipment_id" form:"shipment_id" binding:"required"`
	SampleFlgID      *string       `gorm:"not null;" json:"sample_flg_id" form:"sample_flg_id" binding:"required"`
	Seq              int64         `form:"seq" json:"seq"`
	Vendor           string        `gorm:"size:5;" form:"vendor" json:"vendor"`
	Cd               string        `gorm:"size:5;" form:"cd" json:"cd"`
	Tagrp            string        `gorm:"size:5;" form:"tagrp" json:"tagrp"`
	Sortg1           string        `gorm:"size:25" form:"sortg1" json:"sortg1"`
	Sortg2           string        `gorm:"size:25" form:"sortg2" json:"sortg2"`
	Sortg3           string        `gorm:"size:25" form:"sortg3" json:"sortg3"`
	PlanType         string        `gorm:"size:25" form:"plan_type" json:"plan_type"`
	OrderGroup       string        `gorm:"size:25" form:"order_group" json:"order_groups"`
	Pono             string        `gorm:"size:25" form:"pono" json:"pono"`
	RecId            string        `gorm:"size:25" form:"rec_id" json:"rec_id"`
	Biac             string        `gorm:"size:25" form:"biac" json:"biac"`
	EtdTap           time.Time     `gorm:"type:date;" form:"etd_tap" json:"etd_tap"`
	PartNo           string        `gorm:"size:25" form:"part_no" json:"part_no"`
	PartName         string        `gorm:"size:50" form:"part_name" json:"part_name"`
	SampFlg          string        `gorm:"column:sample_flg;size:2" form:"sample_flg" json:"sample_flg"`
	Orderorgi        float64       `form:"orderorgi" json:"orderorgi"`
	Orderround       float64       `form:"orderround" json:"orderround"`
	FirmFlg          string        `gorm:"size:2" form:"firm_flg" json:"firm_flg"`
	ShippedFlg       string        `gorm:"size:2" form:"shipped_flg" json:"shipped_flg"`
	ShippedQty       float64       `form:"shipped_qty" json:"shipped_qty"`
	Ordermonth       time.Time     `gorm:"type:date;" form:"ordermonth" json:"ordermonth"`
	BalQty           float64       `form:"balqty" json:"balqty"`
	Bidrfl           string        `gorm:"size:2" form:"bidrfl" json:"bidrfl"`
	DeleteFlg        string        `gorm:"size:2" form:"delete_flg" json:"delete_flg"`
	Reasoncd         string        `gorm:"size:5" orm:"reasoncd" json:"reasoncd"`
	Upddte           time.Time     `gorm:"type:date;" form:"upddte" json:"upddte"`
	Updtime          time.Time     `gorm:"type:Time;" form:"updtime" json:"updtime"`
	CarrierCode      string        `gorm:"size:5" form:"carrier_code" json:"carrier_code"`
	Bioabt           int64         `form:"bioabt" json:"bioabt"`
	Bicomd           string        `gorm:"size:2" form:"bicomd" json:"bicomd"`
	Bistdp           float64       `form:"bistdp" json:"bistdp"`
	Binewt           float64       `form:"binewt" json:"binewt"`
	Bigrwt           float64       `form:"bigrwt" json:"bigrwt"`
	Bishpc           string        `gorm:"size:25" form:"bishpc" json:"bishpc"`
	Biivpx           string        `gorm:"size:5" form:"biivpx" json:"biivpx"`
	Bisafn           string        `gorm:"size:25" form:"bisafn" json:"bisafn"`
	Biwidt           float64       `form:"biwidt" json:"biwidt"`
	Bihigh           float64       `form:"bihigh" json:"bihigh"`
	Bileng           float64       `form:"bileng" json:"bileng"`
	LotNo            string        `gorm:"size:25" form:"lotno" json:"lotno"`
	Minimum          int64         `form:"minimum" json:"minimum"`
	Maximum          int64         `form:"maximum" json:"maximum"`
	Picshelfbin      string        `gorm:"size:25" form:"picshelfbin" json:"picshelfbin"`
	Stkshelfbin      string        `gorm:"size:25" form:"stkshelfbin" json:"stkshelfbin"`
	Ovsshelfbin      string        `gorm:"size:25" form:"ovsshelfbin" json:"ovsshelfbin"`
	PicshelfbasicQty float64       `form:"picshelfbasicqty" json:"picshelfbasicqty"`
	OuterPcs         float64       `form:"outerpcs" json:"outerpcs"`
	AllocateQty      float64       `json:"allocate_qty" form:"allocate_qty"`
	Description      string        `json:"description" form:"description"`
	IsReviseError    bool          `json:"is_revise_error" form:"is_revise_error" default:"false"`
	IsGenerate       bool          `json:"is_generate" form:"is_generate" default:"false"`
	ByManually       bool          `json:"by_manually" form:"by_manually" default:"false"`
	IsSync           bool          `json:"is_sync" form:"is_sync" default:"false"`
	IsActive         bool          `json:"is_active" form:"is_active" binding:"required"`
	CreatedAt        time.Time     `json:"created_at" form:"created_at" default:"now"`
	UpdatedAt        time.Time     `json:"updated_at" form:"updated_at" default:"now"`
	FileEdi          FileEdi       `gorm:"foreignKey:FileEdiID;references:ID;" json:"file_edi"`
	Whs              Whs           `gorm:"foreignKey:WhsID;references:ID;" json:"whs"`
	Consignee        Consignee     `gorm:"foreignKey:ConsigneeID;references:ID" json:"consignee"`
	ReviseOrder      ReviseOrder   `gorm:"foreignKey:ReviseOrderID;references:ID" json:"revise_order"`
	Ledger           Ledger        `gorm:"foreignKey:LedgerID;references:ID" json:"ledger"`
	Pc               Pc            `gorm:"foreignKey:PcID;references:ID" json:"pc"`
	Commercial       Commercial    `gorm:"foreignKey:CommercialID;references:ID" json:"commercial"`
	OrderType        OrderType     `gorm:"foreignKey:OrderTypeID;references:ID" json:"order_type"`
	Shipment         Shipment      `gorm:"foreignKey:ShipmentID;references:ID" json:"shipment"`
	OrderZone        OrderZone     `gorm:"foreignKey:OrderZoneID;references:ID" json:"orderzone"`
	SampleFlg        SampleFlg     `gorm:"foreignKey:SampleFlgID;references:ID" json:"sampleflg"`
	OrderDetail      []OrderDetail `json:"order_details"`
}

type Affcode struct {
	ID          string    `gorm:"primaryKey;size:21" json:"id,omitempty"`
	Title       string    `gorm:"not null;unique;size:15" json:"title,omitempty" form:"title" binding:"required"`
	Description string    `json:"description,omitempty" form:"description" binding:"required"`
	IsActive    bool      `json:"is_active,omitempty" form:"is_active" binding:"required"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
}

type Customer struct {
	ID          string    `gorm:"primaryKey;size:21" json:"id,omitempty"`
	Title       string    `gorm:"not null;unique;size:15" json:"title,omitempty" form:"title" binding:"required"`
	Description string    `json:"description,omitempty" form:"description" binding:"required"`
	IsActive    bool      `json:"is_active,omitempty" form:"is_active" binding:"required"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
}

type CustomerAddress struct {
	ID          string    `gorm:"primaryKey;size:21" json:"id,omitempty"`
	Title       string    `gorm:"not null;unique;size:50" json:"title,omitempty" form:"title" binding:"required"`
	Description string    `json:"description,omitempty" form:"description" binding:"required"`
	IsActive    bool      `json:"is_active,omitempty" form:"is_active" binding:"required"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
}

type Consignee struct {
	ID                string          `gorm:"primaryKey;size:21" json:"id,omitempty"`
	WhsID             *string         `gorm:"not null;" json:"whs_id,omitempty" form:"whs_id"`
	FactoryID         *string         `gorm:"not null;" json:"factory_id,omitempty" form:"factory_id"`
	AffcodeID         *string         `gorm:"not null;" json:"affcode_id,omitempty" form:"affcode_id" binding:"required"`
	CustomerID        *string         `gorm:"not null;" json:"customer_id,omitempty" form:"customer_id" binding:"required"`
	CustomerAddressID *string         `gorm:"null;" json:"customer_ddress_id,omitempty" form:"customer_address_id"`
	Prefix            string          `gorm:"not null" json:"prefix,omitempty" form:"prefix" binding:"required"`
	IsActive          bool            `json:"is_active,omitempty" form:"is_active" binding:"required"`
	CreatedAt         time.Time       `json:"created_at,omitempty" default:"now"`
	UpdatedAt         time.Time       `json:"updated_at,omitempty" default:"now"`
	Whs               Whs             `gorm:"foreignKey:WhsID;references:ID" json:"whs,omitempty"`
	Factory           Factory         `gorm:"foreignKey:FactoryID;references:ID" json:"factory,omitempty"`
	Affcode           Affcode         `gorm:"foreignKey:AffcodeID;references:ID" json:"affcode,omitempty"`
	Customer          Customer        `gorm:"foreignKey:CustomerID;references:ID" json:"customer,omitempty"`
	CustomerAddress   CustomerAddress `gorm:"foreignKey:CustomerAddressID;references:ID" json:"customer_address,omitempty"`
	OrderGroup        []*OrderGroup   `json:"order_group,omitempty"`
}

type Shipment struct {
	ID          string    `gorm:"primaryKey;size:21" json:"id,omitempty"`
	Title       string    `gorm:"not null;unique;size:15" json:"title,omitempty" form:"title" binding:"required"`
	Description string    `json:"description,omitempty" form:"description" binding:"required"`
	IsActive    bool      `json:"is_active,omitempty" form:"is_active" binding:"required"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
}

type Pc struct {
	ID          string    `gorm:"primaryKey;size:21" json:"id,omitempty"`
	Title       string    `gorm:"not null;unique;size:50" json:"title,omitempty" form:"title" binding:"required"`
	Description string    `json:"description,omitempty" form:"description" binding:"required"`
	IsActive    bool      `json:"is_active,omitempty" form:"is_active" binding:"required"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
}

type Commercial struct {
	ID          string    `gorm:"primaryKey;size:21" json:"id,omitempty"`
	Title       string    `gorm:"not null;unique;size:50" json:"title,omitempty" form:"title" binding:"required"`
	Prefix      string    `gorm:"size:5" json:"prefix,omitempty" form:"prefix" binding:"required"`
	Description string    `json:"description,omitempty" form:"description" binding:"required"`
	IsActive    bool      `json:"is_active,omitempty" form:"is_active" binding:"required"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
}

type SampleFlg struct {
	ID          string    `gorm:"primaryKey;size:21" json:"id,omitempty"`
	Title       string    `gorm:"not null;unique;size:50" json:"title,omitempty" form:"title" binding:"required"`
	Description string    `json:"description,omitempty" form:"description" binding:"required"`
	IsActive    bool      `json:"is_active,omitempty" form:"is_active" binding:"required"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
}

type OrderTitle struct {
	ID          string    `gorm:"primaryKey;size:21" json:"id"`
	Title       string    `gorm:"not null;unique;size:15" json:"title" form:"title" binding:"required"`
	Description string    `json:"description" form:"description" binding:"required"`
	IsActive    bool      `json:"is_active" form:"is_active" binding:"required"`
	CreatedAt   time.Time `json:"created_at" default:"now"`
	UpdatedAt   time.Time `json:"updated_at" default:"now"`
}

type Pallet struct {
	ID           string         `gorm:"primaryKey;size:21" json:"id,omitempty"`
	OrderID      *string        `gorm:"not null;" json:"order_id,omitempty" form:"order_id" binding:"required"`
	PalletTypeID *string        `gorm:"not null;" json:"pallet_type_id,omitempty" form:"pallet_type_id" binding:"required"`
	PalletPrefix string         `gorm:"not null;size:1;" json:"pallet_prefix,omitempty" form:"pallet_prefix" default:"P"`
	PalletNo     int64          `gorm:"not null;" json:"pallet_no,omitempty" form:"pallet_no" binding:"required"`
	PalletTotal  int64          `gorm:"not null;" json:"pallet_total,omitempty" form:"pallet_total" binding:"required"`
	IsSync       bool           `json:"is_sync,omitempty" form:"is_sync" binding:"required"`
	IsActive     bool           `json:"is_active,omitempty" form:"is_active" binding:"required"`
	CreatedAt    time.Time      `json:"created_at,omitempty" default:"now"`
	UpdatedAt    time.Time      `json:"updated_at,omitempty" default:"now"`
	Order        Order          `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"`
	PalletType   PalletType     `gorm:"foreignKey:PalletTypeID;references:ID"  json:"pallet_type"`
	PalletDetail []PalletDetail `json:"pallet_detail,omitempty"`
}

type PalletDetail struct {
	ID            string      `gorm:"primaryKey;size:21" json:"id,omitempty"`
	PalletID      *string     `json:"pallet_id,omitempty" form:"pallet_id" binding:"required"`
	OrderDetailID *string     `gorm:"not null;" json:"order_detail_id,omitempty" form:"order_detail_id" binding:"required"`
	SeqNo         int64       `gorm:"not null;" json:"seq_no" form:"seq_no" binding:"required"`
	IsPrintLabel  bool        `json:"is_print_label,omitempty" form:"is_print_label" binding:"required"`
	IsActive      bool        `json:"is_active,omitempty" form:"is_active" binding:"required"`
	CreatedAt     time.Time   `json:"created_at,omitempty" default:"now"`
	UpdatedAt     time.Time   `json:"updated_at,omitempty" default:"now"`
	Pallet        Pallet      `gorm:"foreignKey:PalletID;references:ID" json:"pallet,omitempty"`
	OrderDetail   OrderDetail `gorm:"foreignKey:OrderDetailID;references:ID"  json:"order_detail,omitempty"`
}

type Area struct {
	ID          string    `gorm:"primaryKey,unique;size:21;" json:"id,omitempty"`
	Title       string    `gorm:"not null;size:50;unique;" json:"title,omitempty" form:"title" binding:"required"`
	Description string    `gorm:"null" json:"description,omitempty" form:"description"`
	IsActive    bool      `gorm:"null" json:"is_active,omitempty" form:"is_active" default:"false"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
}
type Whs struct {
	ID          string    `gorm:"primaryKey,unique;size:21;" json:"id,omitempty"`
	Title       string    `gorm:"not null;size:50;unique;" json:"title,omitempty" form:"title" binding:"required"`
	Value       string    `gorm:"size:5;" json:"value,omitempty" form:"value" binding:"required"`
	Description string    `gorm:"null" json:"description,omitempty" form:"description"`
	IsActive    bool      `gorm:"null" json:"is_active,omitempty" form:"is_active" default:"false"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
}

type Factory struct {
	ID          string    `gorm:"primaryKey,unique;size:21;" json:"id,omitempty"`
	Title       string    `gorm:"not null;size:50;unique;" json:"title,omitempty" form:"title" binding:"required"`
	InvPrefix   string    `gorm:"size:5;" json:"inv_prefix,omitempty" form:"inv_prefix" binding:"required"`
	LabelPrefix string    `gorm:"size:5;" json:"label_prefix,omitempty" form:"label_prefix" binding:"required"`
	PartUnit    string    `gorm:"size:50;" json:"part_unit,omitempty" form:"part_unit" binding:"required"`
	CdCode      string    `gorm:"size:5;" json:"cd_code,omitempty" form:"cd_code" binding:"required"`
	PartType    string    `gorm:"size:50;" json:"part_type,omitempty" form:"part_type" binding:"required"`
	Sortg1      string    `gorm:"size:50;" json:"sortg1,omitempty" form:"sortg1" binding:"required"`
	Sortg2      string    `gorm:"size:50;" json:"sortg2,omitempty" form:"sortg2" binding:"required"`
	Sortg3      string    `gorm:"size:50;" json:"sortg3,omitempty" form:"sortg3" binding:"required"`
	Description string    `gorm:"null" json:"description,omitempty" form:"description"`
	IsActive    bool      `gorm:"null" json:"is_active,omitempty" form:"is_active" default:"false"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
}

type PrefixName struct {
	ID          string    `gorm:"primaryKey,unique;size:21;" json:"id,omitempty"`
	Title       string    `gorm:"not null;size:50;unique;" json:"title,omitempty" form:"title" binding:"required"`
	Description string    `gorm:"null" json:"description,omitempty" form:"description"`
	IsActive    bool      `gorm:"null" json:"is_active,omitempty" form:"is_active" default:"false"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
}

type Position struct {
	ID          string    `gorm:"primaryKey,unique;size:21;" json:"id,omitempty"`
	Title       string    `gorm:"not null;size:50;unique;" json:"title,omitempty" form:"title" binding:"required"`
	Description string    `gorm:"null" json:"description,omitempty" form:"description"`
	IsActive    bool      `gorm:"null" json:"is_active,omitempty" form:"is_active" default:"false"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
}

type Department struct {
	ID          string    `gorm:"primaryKey,unique;size:21;" json:"id,omitempty"`
	Title       string    `gorm:"not null;size:50;unique;" json:"title,omitempty" form:"title" binding:"required"`
	Description string    `gorm:"null" json:"description,omitempty" form:"description"`
	IsActive    bool      `gorm:"null" json:"is_active,omitempty" form:"is_active" default:"false"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
}

type Unit struct {
	ID          string    `gorm:"primaryKey,unique;size:21;" json:"id,omitempty"`
	Title       string    `gorm:"not null;size:50;unique;" json:"title,omitempty" form:"title" binding:"required"`
	Description string    `gorm:"null" json:"description,omitempty" form:"description"`
	IsActive    bool      `gorm:"null" json:"is_active,omitempty" form:"is_active" default:"false"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
}

type PartType struct {
	ID          string    `gorm:"primaryKey,unique;size:21;" json:"id,omitempty"`
	Title       string    `gorm:"not null;size:50;unique;" json:"title,omitempty" form:"title" binding:"required"`
	Description string    `gorm:"null" json:"description,omitempty" form:"description"`
	IsActive    bool      `gorm:"null" json:"is_active,omitempty" form:"is_active" default:"false"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
}

type ReceiveType struct {
	ID          string    `gorm:"primaryKey,unique;size:21;" json:"id,omitempty"`
	WhsID       string    `gorm:"not null;" json:"whs_id,omitempty" form:"whs_id"`
	Title       string    `gorm:"not null;size:50;unique;" json:"title,omitempty" form:"title" binding:"required"`
	Description string    `gorm:"null" json:"description,omitempty" form:"description"`
	IsActive    bool      `gorm:"null" json:"is_active,omitempty" form:"is_active" default:"false"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
	Whs         Whs       `gorm:"foreignKey:WhsID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"whs,omitempty"`
}

type FileType struct {
	ID          string    `gorm:"primaryKey;size:21" json:"id,omitempty"`
	Title       string    `gorm:"not null;unique;size:50" json:"title,omitempty" form:"title" binding:"required"`
	Description string    `json:"description,omitempty" form:"description" binding:"required"`
	IsActive    bool      `json:"is_active,omitempty" form:"is_active" binding:"required"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
}
type Mailbox struct {
	ID        string    `gorm:"primaryKey;size:21" json:"id,omitempty"`
	Mailbox   string    `gorm:"not null;unique;size:21" json:"mailbox,omitempty" form:"mailbox" binding:"required"`
	Password  string    `gorm:"size:50" json:"password,omitempty" form:"password" binding:"required"`
	HostUrl   string    `json:"host_url,omitempty" form:"host_url" binding:"required"`
	AreaID    *string   `json:"area_id,omitempty" form:"area_id" binding:"required"`
	IsActive  bool      `json:"is_active,omitempty" form:"is_active" binding:"required"`
	CreatedAt time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt time.Time `json:"updated_at,omitempty" default:"now"`
	Area      Area      `gorm:"foreignKey:AreaID;references:ID;constraint:OnDelete:CASCADE;" json:"area,omitempty"`
}

type ReviseOrder struct {
	ID          string    `gorm:"primaryKey;size:21" json:"id,omitempty"`
	Title       string    `gorm:"not null;unique;size:15" json:"title,omitempty" form:"title" binding:"required"`
	Description string    `json:"description,omitempty" form:"description" binding:"required"`
	IsActive    bool      `json:"is_active,omitempty" form:"is_active" binding:"required"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
}

type OrderType struct {
	ID          string    `gorm:"primaryKey;size:21" json:"id,omitempty"`
	Title       string    `gorm:"not null;unique;size:15" json:"title,omitempty" form:"title" binding:"required"`
	Description string    `json:"description,omitempty" form:"description" binding:"required"`
	IsActive    bool      `json:"is_active,omitempty" form:"is_active" binding:"required"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
}

// N=All,F=3 Front,E=3 End,O=Sprit Order
type OrderGroupType struct {
	ID          string    `gorm:"primaryKey;size:21" json:"id,omitempty"`
	Title       string    `gorm:"not null;unique;size:15" json:"title,omitempty" form:"title" binding:"required"`
	Description string    `json:"description,omitempty" form:"description" binding:"required"`
	IsActive    bool      `json:"is_active,omitempty" form:"is_active" binding:"required"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
}

type OrderGroup struct {
	ID               string          `gorm:"primaryKey;size:21" json:"id,omitempty"`
	UserID           *string         `gorm:"not null;" json:"user_id,omitempty" form:"user_id" binding:"required"`
	ConsigneeID      *string         `gorm:"not null;" json:"consignee_id,omitempty" form:"consignee_id" binding:"required"`
	OrderGroupTypeID *string         `gorm:"not null;" json:"order_group_type_id,omitempty" form:"order_group_type_id" binding:"required"`
	SubOrder         string          `gorm:"not null;size:15" json:"sub_order,omitempty" form:"sub_order" binding:"required"`
	Description      string          `json:"description,omitempty" form:"description" binding:"required"`
	IsActive         bool            `json:"is_active,omitempty" form:"is_active" binding:"required"`
	CreatedAt        time.Time       `json:"created_at,omitempty" default:"now"`
	UpdatedAt        time.Time       `json:"updated_at,omitempty" default:"now"`
	User             *User           `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Consignee        *Consignee      `gorm:"foreignKey:ConsigneeID;references:ID" json:"consignee,omitempty"`
	OrderGroupType   *OrderGroupType `gorm:"foreignKey:OrderGroupTypeID;references:ID" json:"order_group_type,omitempty"`
}

type FormGroupConsignee struct {
	UserID           string `json:"user_id,omitempty" form:"user_id" binding:"required"`
	WhsID            string `json:"whs_id,omitempty" form:"whs_id"`
	FactoryID        string `json:"factory_id,omitempty" form:"factory_id"`
	AffcodeID        string `json:"affcode_id,omitempty" form:"affcode_id" binding:"required"`
	CustcodeID       string `json:"custcode_id,omitempty" form:"custcode_id" binding:"required"`
	OrderGroupTypeID string `json:"order_group_type_id,omitempty" form:"order_group_type_id" binding:"required"`
	SubOrder         string `gorm:"size:15" json:"sub_order,omitempty" form:"sub_order" binding:"required"`
	Description      string `json:"description,omitempty" form:"description" binding:"required"`
	IsActive         bool   `json:"is_active,omitempty" form:"is_active" binding:"required"`
}

type Location struct {
	ID          string    `gorm:"primaryKey,unique;size:21;" json:"id,omitempty"`
	Title       string    `gorm:"not null;size:50;unique;" json:"title,omitempty" form:"title" binding:"required"`
	Description string    `gorm:"null" json:"description,omitempty" form:"description"`
	MaxLimit    int64     `gorm:"null" json:"max_limit,omitempty" form:"max_limit" binding:"required"`
	IsActive    bool      `gorm:"null" json:"is_active,omitempty" form:"is_active" default:"true"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
}

type AutoGenerateInvoice struct {
	ID         string    `gorm:"primaryKey,unique;size:21;" json:"id,omitempty"`
	FactoryID  *string   `gorm:"not null;unique;" json:"factory_id,omitempty" form:"factory_id,omitempty"`
	IsGenerate bool      `json:"is_generate" form:"is_generate" default:"true"`
	IsActive   bool      `json:"is_active,omitempty" form:"is_active" default:"true"`
	CreatedAt  time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt  time.Time `json:"updated_at,omitempty" default:"now"`
	Factory    Factory   `gorm:"foreignKey:FactoryID;references:ID" json:"factory,omitempty"`
}

type LineNotifyToken struct {
	ID        string    `gorm:"primaryKey,unique;size:21;" json:"id,omitempty"`
	WhsID     *string   `gorm:"not null;" json:"whs_id,omitempty" form:"whs_id,omitempty"`
	FactoryID *string   `gorm:"not null;" json:"factory_id,omitempty" form:"factory_id,omitempty"`
	Token     string    `gorm:"not null;unique;" json:"token,omitempty" form:"token"`
	IsActive  bool      `json:"is_active,omitempty" form:"is_active" default:"true"`
	CreatedAt time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt time.Time `json:"updated_at,omitempty" default:"now"`
	Whs       Whs       `gorm:"foreignKey:WhsID;references:ID" json:"whs,omitempty"`
	Factory   Factory   `gorm:"foreignKey:FactoryID;references:ID" json:"factory,omitempty"`
}

type PalletType struct {
	ID               string    `gorm:"primaryKey,unique;size:21;" json:"id,omitempty"`
	Type             string    `gorm:"not null;unique;" json:"type,omitempty" form:"type"`
	Floors           int64     `json:"floors" form:"floors" default:"0"`
	BoxSizeWidth     float64   `json:"box_size_width" form:"box_size_width" default:"0"`
	BoxSizeLength    float64   `json:"box_size_length" form:"box_size_length" default:"0"`
	BoxSizeHight     float64   `json:"box_size_hight" form:"box_size_hight" default:"0"`
	PalletSizeWidth  float64   `json:"pallet_size_width" form:"pallet_size_width" default:"0"`
	PalletSizeLength float64   `json:"pallet_size_length" form:"pallet_size_length" default:"0"`
	PalletSizeHight  float64   `json:"pallet_size_hight" form:"pallet_size_hight" default:"0"`
	LimitTotal       int64     `json:"limit_total" form:"limit_total" default:"0"`
	IsActive         bool      `json:"is_activey" form:"is_active" default:"true"`
	CreatedAt        time.Time `json:"created_at" form:"created_at" default:"now"`
	UpdatedAt        time.Time `json:"updated_at" form:"updated_at" default:"now"`
	// ชนดิ กลอ่ ง จ านวนชนั้ BOX SIZE PALLET SIZE BOX/PALLET
}

type LastFticket struct {
	ID          string    `gorm:"primaryKey;size:21" json:"id,omitempty"`
	FactoryID   *string   `gorm:"not null;unique;" json:"factory_id,omitempty" form:"factory_id" binding:"required"`
	OnYear      int64     `gorm:"not null;" json:"on_year,omitempty" form:"on_year"`
	LastRunning int64     `json:"last_running,omitempty" form:"last_running" binding:"required"`
	IsActive    bool      `json:"is_active,omitempty" form:"is_active" binding:"required"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
	Factory     Factory   `gorm:"foreignKey:FactoryID;references:ID" json:"factory,omitempty"`
}

type PlanningDay struct {
	ID          string    `gorm:"primaryKey,unique;size:21;" json:"id,omitempty"`
	Title       string    `gorm:"not null;size:50;unique;" json:"title,omitempty" form:"title" binding:"required"`
	Description string    `gorm:"null" json:"description,omitempty" form:"description"`
	IsActive    bool      `gorm:"null" json:"is_active,omitempty" form:"is_active" default:"false"`
	CreatedAt   time.Time `json:"created_at,omitempty" default:"now"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" default:"now"`
}
