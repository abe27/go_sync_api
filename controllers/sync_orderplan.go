package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/abe27/syncapi/configs"
	"github.com/abe27/syncapi/models"
	_ "gopkg.in/goracle.v2"
)

func UpdateOrderPlan(id string, isSync bool) {
	// fmt.Printf("Update Order Plan ID: %s SYNC: %v\n", id, isSync)

	payload := strings.NewReader(fmt.Sprintf("is_sync=%v", isSync))
	client := &http.Client{}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/sync/orderplan/%s", configs.API_HOST, id), payload)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if _, err = io.ReadAll(res.Body); err != nil {
		panic(err)
	}
	// fmt.Println(string(body))
}

func GetSyncOrderPlan() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/sync/orderplan", configs.API_HOST), nil)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var SyncOrderPlan models.SyncOrderPlan
	if err := json.Unmarshal(body, &SyncOrderPlan); err != nil {
		panic(err)
	}

	// Connect OraDB
	db, err := sql.Open("goracle", configs.ORAC_USER+"/"+configs.ORAC_PASSWORD+"@"+configs.ORAC_HOST+"/"+configs.ORAC_SERVICE)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	x := 1
	for _, v := range SyncOrderPlan.Data {
		sqlQuery := fmt.Sprintf("SELECT ORDERID FROM TXP_ORDERPLAN WHERE FACTORY='%s' AND SHIPTYPE='%s' AND AFFCODE='%s' AND PONO='%s' AND ETDTAP=TO_DATE('%s','YYYY-MM-DD') AND PARTNO='%s' AND PC='%s' AND COMMERCIAL='%s' AND BISHPC='%s'", v.Vendor, v.Shipment.Title, v.Biac, v.Pono, v.EtdTap.Format("2006-01-02"), v.PartNo, v.Pc.Title, v.Commercial.Title, v.Bishpc)
		// fmt.Println(sqlQuery)
		rows, err := db.Query(sqlQuery)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		txtQuery := "INSERT"
		sqlExecute := fmt.Sprintf("INSERT INTO TXP_ORDERPLAN(FACTORY, SHIPTYPE, AFFCODE, PONO, ETDTAP, PARTNO, PARTNAME, ORDERMONTH, ORDERORGI, ORDERROUND, BALQTY, SHIPPEDFLG, SHIPPEDQTY, PC, COMMERCIAL, SAMPFLG, CARRIERCODE, ORDERTYPE, UPDDTE, ALLOCATEQTY, BIDRFL, DELETEFLG, REASONCD, BIOABT, FIRMFLG, BICOMD, BISTDP, BINEWT, BIGRWT, BISHPC, BIIVPX, BISAFN, BILENG, BIWIDT, BIHIGH, CURINV, OLDINV, SYSDTE, POUPDFLAG, CREATEDBY, MODIFIEDBY, LOTNO, ORDERSTATUS, ORDERID, STATUS, ORDSYNC)VALUES('%s', '%s', '%s', '%s', TO_DATE('%s', 'YYYY-MM-DD'), '%s', '%s', TO_DATE('%s', 'YYYY-MM-DD'), %f, %f, %f, '%s', %f, '%s', '%s', '%s', '%s', '%s', SYSDATE, %f, '%s', '%s', '%s', '%d', '%s', '%s', %f, %f, %f, '%s','%s', '%s', %f, %f, %f, null,null,SYSDATE, null,null,null, '%s', 0, '%s', 0, 0)", v.Vendor, v.Shipment.Title, v.Biac, v.Pono, v.EtdTap.Format("2006-01-02"), v.PartNo, strings.ReplaceAll(v.PartName, "'", "''"), v.Ordermonth.Format("2006-01-02"), v.Orderorgi, v.Orderround, v.BalQty, v.ShippedFlg, v.ShippedQty, v.Pc.Title, v.Commercial.Title, v.SampFlg, v.CarrierCode, v.OrderType.Title, v.AllocateQty, v.Bidrfl, v.DeleteFlg, v.Reasoncd, v.Bioabt, v.FirmFlg, v.Bicomd, v.Bistdp, v.Binewt, v.Bigrwt, v.Bishpc, v.Biivpx, v.Bisafn, v.Bileng, v.Biwidt, v.Bihigh, v.LotNo, v.ID)
		if rows.Next() {
			sqlExecute = fmt.Sprintf("UPDATE TXP_ORDERPLAN SET ORDERMONTH=TO_DATE('%s','YYYY-MM-DD'), ORDERORGI=%f, ORDERROUND=%f, BALQTY=%f, SHIPPEDFLG='%s', SHIPPEDQTY=%f, SAMPFLG='%s', CARRIERCODE='%s', ORDERTYPE='%s', UPDDTE=CURRENT_TIMESTAMP, ALLOCATEQTY=%f, BIDRFL='%s', DELETEFLG='%s', REASONCD='%s', BIOABT='%d', FIRMFLG='%s', BICOMD='%s', BISTDP=%f, BINEWT=%f, BIGRWT=%f, BIIVPX='%s', BISAFN='%s', BILENG=%f, BIWIDT=%f, BIHIGH=%f,LOTNO='%s', ORDERID='%s',ORDSYNC=0 WHERE FACTORY='%s' AND SHIPTYPE='%s' AND AFFCODE='%s' AND PONO='%s' AND ETDTAP=TO_DATE('%s','YYYY-MM-DD') AND PARTNO='%s' AND PC='%s' AND COMMERCIAL='%s' AND BISHPC='%s'", v.EtdTap.Format("2006-01-02"), v.Orderorgi, v.Orderround, v.BalQty, v.ShippedFlg, v.ShippedQty, v.SampFlg, v.CarrierCode, v.OrderType.Title, v.AllocateQty, v.Bidrfl, v.DeleteFlg, v.Reasoncd, v.Bioabt, v.FirmFlg, v.Bicomd, v.Bistdp, v.Binewt, v.Bigrwt, v.Biivpx, v.Bisafn, v.Bileng, v.Biwidt, v.Bihigh, v.LotNo, v.ID, v.Vendor, v.Shipment.Title, v.Biac, v.Pono, v.EtdTap.Format("2006-01-02"), v.PartNo, v.Pc.Title, v.Commercial.Title, v.Bishpc)
			txtQuery = "UPDATE"
		}

		// fmt.Printf("%s\n", sqlExecute)
		_, err = db.Exec(sqlExecute)
		if err != nil {
			panic(err)
		}

		sqlInsertLogs := fmt.Sprintf("INSERT INTO TXP_ORDERRECORDS(FACTORY, SHIPTYPE, AFFCODE, PONO, ETDTAP, PARTNO, PARTNAME, ORDERMONTH, ORDERORGI, ORDERROUND, BALQTY, SHIPPEDFLG, SHIPPEDQTY, PC, COMMERCIAL, SAMPFLG, CARRIERCODE, ORDERTYPE, UPDDTE, ALLOCATEQTY, BIDRFL, DELETEFLG, REASONCD, BIOABT, FIRMFLG, BICOMD, BISTDP, BINEWT, BIGRWT, BISHPC, BIIVPX, BISAFN, BILENG, BIWIDT, BIHIGH, CURINV, OLDINV, SYSDTE, POUPDFLAG, CREATEDBY, MODIFIEDBY, LOTNO, ORDERSTATUS, ORDERID, STATUS, ORDSYNC)VALUES('%s', '%s', '%s', '%s', TO_DATE('%s', 'YYYY-MM-DD'), '%s', '%s', TO_DATE('%s', 'YYYY-MM-DD'), %f, %f, %f, '%s', %f, '%s', '%s', '%s', '%s', '%s', SYSDATE, %f, '%s', '%s', '%s', '%d', '%s', '%s', %f, %f, %f, '%s','%s', '%s', %f, %f, %f, null,null,SYSDATE, null,null,null, '%s', 0, '%s', 0, 0)", v.Vendor, v.Shipment.Title, v.Biac, v.Pono, v.EtdTap.Format("2006-01-02"), v.PartNo, strings.ReplaceAll(v.PartName, "'", "''"), v.Ordermonth.Format("2006-01-02"), v.Orderorgi, v.Orderround, v.BalQty, v.ShippedFlg, v.ShippedQty, v.Pc.Title, v.Commercial.Title, v.SampFlg, v.CarrierCode, v.OrderType.Title, v.AllocateQty, v.Bidrfl, v.DeleteFlg, v.Reasoncd, v.Bioabt, v.FirmFlg, v.Bicomd, v.Bistdp, v.Binewt, v.Bigrwt, v.Bishpc, v.Biivpx, v.Bisafn, v.Bileng, v.Biwidt, v.Bihigh, v.LotNo, v.ID)
		_, err = db.Exec(sqlInsertLogs)
		if err != nil {
			panic(err)
		}
		UpdateOrderPlan(v.ID, true)
		fmt.Printf("%d %s ORDER ID: %s\n", x, txtQuery, v.ID)
		x++
	}
}
