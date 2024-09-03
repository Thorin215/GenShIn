package model

import (
	"application/sql"
	"time"

	"gorm.io/gorm"
)

// Selling 销售要约
// 需要确定ObjectOfSale是否属于Seller
// 买家初始为空
// Seller和ObjectOfSale一起作为复合键,保证可以通过seller查询到名下所有发起的销售
type Selling struct {
	ObjectOfSale  string  `json:"objectOfSale"`  //销售对象(正在出售的房地产RealEstateID)
	Seller        string  `json:"seller"`        //发起销售人、卖家(卖家AccountId)
	Buyer         string  `json:"buyer"`         //参与销售人、买家(买家AccountId)
	Price         float64 `json:"price"`         //价格
	CreateTime    string  `json:"createTime"`    //创建时间
	SalePeriod    int     `json:"salePeriod"`    //智能合约的有效期(单位为天)
	SellingStatus string  `json:"sellingStatus"` //销售状态
}

// SellingStatusConstant 销售状态
var SellingStatusConstant = func() map[string]string {
	return map[string]string{
		"saleStart": "销售中", //正在销售状态,等待买家光顾
		"cancelled": "已取消", //被卖家取消销售或买家退款操作导致取消
		"expired":   "已过期", //销售期限到期
		"delivery":  "交付中", //买家买下并付款,处于等待卖家确认收款状态,如若卖家未能确认收款，买家可以取消并退款
		"done":      "完成",  //卖家确认接收资金，交易完成
	}
}

// Donating 捐赠要约
// 需要确定ObjectOfDonating是否属于Donor
// 需要指定受赠人Grantee，并等待受赠人同意接收
type Donating struct {
	ObjectOfDonating string `json:"objectOfDonating"` //捐赠对象(正在捐赠的房地产RealEstateID)
	Donor            string `json:"donor"`            //捐赠人(捐赠人AccountId)
	Grantee          string `json:"grantee"`          //受赠人(受赠人AccountId)
	CreateTime       string `json:"createTime"`       //创建时间
	DonatingStatus   string `json:"donatingStatus"`   //捐赠状态
}

// DonatingStatusConstant 捐赠状态
var DonatingStatusConstant = func() map[string]string {
	return map[string]string{
		"donatingStart": "捐赠中", //捐赠人发起捐赠合约，等待受赠人确认受赠
		"cancelled":     "已取消", //捐赠人在受赠人确认受赠之前取消捐赠或受赠人取消接收受赠
		"done":          "完成",  //受赠人确认接收，交易完成
	}
}

type TData_Language struct {
	Context string `json:"context"`
	Label   bool   `json:"label"`
	SetNo   string `json:"setNo"`
}

type TData_Set struct {
	SetNo    string `json:"setNo"`
	SetDonor string `json:"setDonor"`
	SetSize  string `json:"setSize"`
	SetNum   string `json:"setNum"`
}

type DataSet struct {
	Name       string    `json:"name"`
	AccountID  int32     `json:"account_id"`
	Stars      int32     `json:"stars"`
	DataSetID  int32     `json:"data_set_id" gorm:"primary_key"`
	CreateTime time.Time `json:"create_time"`
	ModifyTime time.Time `json:"modified_time"`
	Data       []byte    `json:"data"`
	Downloads  int32     `json:"downloads"`
}

type MetaData struct {
	DataSetID  int32  `json:"data_set_id" gorm:"primary_key"`
	Tasks      string `json:"tasks"`
	Modalities string `json:"modalities"`
	Formats    string `json:"formats"`
	SubTasks   string `json:"sub_tasks"`
	Languages  string `json:"languages"`
	Libararies string `json:"libraries"`
	Tags       string `json:"tags"`
	License    string `json:"license"`
	Rows       int32  `json:"rows"`
}

func CreateDataSet(dataSet *DataSet, metaData *MetaData) error {
	error := sql.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(dataSet).Error; err != nil {
			return err
		}
		if err := tx.Create(metaData).Error; err != nil {
			return err
		}
		return nil
	})
	return error
}

func GetDataSet(dataSetID int32) (*DataSet, *MetaData, error) {
	dataSet := &DataSet{}
	metaData := &MetaData{}
	err := sql.DB.Where("data_set_id = ?", dataSetID).First(dataSet).Error
	if err != nil {
		return nil, nil, err
	}
	err = sql.DB.Where("data_set_id = ?", dataSetID).First(metaData).Error
	if err != nil {
		return nil, nil, err
	}
	return dataSet, metaData, nil
}

func UpdateDataSet(dataSet *DataSet) error {
	error := sql.DB.Updates(dataSet).Error
	return error
}

func UpdateMetaData(metaData *MetaData) error {
	error := sql.DB.Updates(metaData).Error
	return error
}

func DeleteDataSet(dataSetID int32) error {
	error := sql.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("data_set_id = ?", dataSetID).Delete(&DataSet{}).Error; err != nil {
			return err
		}
		if err := tx.Where("data_set_id = ?", dataSetID).Delete(&MetaData{}).Error; err != nil {
			return err
		}
		return nil
	})
	return error
}

func GetDataSetList(account_id int32) ([]DataSet, error) {
	dataSetList := []DataSet{}
	err := sql.DB.Where("account_id = ?", account_id).Find(&dataSetList).Error
	if err != nil {
		return nil, err
	}
	return dataSetList, nil
}

func GetDataSetListByName(name string) ([]DataSet, error) {
	dataSetList := []DataSet{}
	err := sql.DB.Where("name = ?", name).Find(&dataSetList).Error
	if err != nil {
		return nil, err
	}
	return dataSetList, nil
}
