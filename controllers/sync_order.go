package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/abe27/syncapi/configs"
	"github.com/abe27/syncapi/models"
	_ "gopkg.in/goracle.v2"
)

func SyncOrder() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/sync/order", configs.API_HOST), nil)

	if err != nil {
		panic(err)
	}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var ord models.SyncOrder
	if err := json.Unmarshal(body, &ord); err != nil {
		panic(err)
	}

	x := 1
	for _, v := range ord.Data {
		// fmt.Printf("%d: %s\n", x,obj.ID)
		CreateInvoice(&x, &v)
		x++
	}
}

func CreateInvoice(x *int, obj *models.Order) {
	if len(obj.OrderDetail) > 0 {
		invNo := fmt.Sprintf("%s%s%s%04d%s", obj.Consignee.Factory.InvPrefix, obj.Consignee.Prefix, (obj.EtdDate.Format("2006-01-02"))[3:4], obj.RunningSeq, obj.Shipment.Title)
		empID := strings.ToUpper(os.Getenv("USERNAME"))
		if len(obj.Consignee.OrderGroup) > 0 {
			empID = obj.Consignee.OrderGroup[0].User.UserName
		}

		db, err := sql.Open("goracle", configs.ORAC_USER+"/"+configs.ORAC_PASSWORD+"@"+configs.ORAC_HOST+"/"+configs.ORAC_SERVICE)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		// Get SPL LastInvoice
		dbQuerySplInv := fmt.Sprintf("SELECT count(UUID) + 1 FROM TXP_ISSTRANSENT WHERE REFINVOICE LIKE '%s'", obj.Consignee.Factory.InvPrefix+obj.Consignee.Prefix+"%")
		rows, err := db.Query(dbQuerySplInv)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		// Last Seq
		var lastSplNo int
		if rows.Next() {
			rows.Scan(&lastSplNo)
		}
		splRef := fmt.Sprintf("%s%s-%s-%04d", "S", obj.Consignee.Prefix, (obj.EtdDate.Format("20060102")), lastSplNo)

		// Create Ent
		dbQuery := fmt.Sprintf("SELECT ISSUINGKEY FROM TXP_ISSTRANSENT WHERE ETDDTE=TO_DATE('%s', 'YYYY-MM-DD') AND FACTORY='%s' AND AFFCODE='%s' AND BISHPC='%s' AND PC='%s' AND COMERCIAL='%s' AND ZONEID='%d' AND SHIPTYPE='%s'", obj.EtdDate.Format("2006-01-02"), obj.Consignee.Factory.Title, obj.Consignee.Affcode.Title, obj.Consignee.Customer.Title, obj.Pc.Title, obj.Commercial.Title, obj.Bioabt, obj.Shipment.Title)
		// fmt.Println(dbQuery)
		rows, err = db.Query(dbQuery)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		var issueKey string
		if rows.Next() {
			rows.Scan(&issueKey)
		}

		orderType := obj.OrderDetail[0].OrderPlan.OrderType.Title
		txtUpdate := "INSERT"
		sqlEntExecute := fmt.Sprintf("INSERT INTO TXP_ISSTRANSENT(ISSUINGKEY, ETDDTE, FACTORY, AFFCODE, BISHPC, CUSTNAME, COMERCIAL, ZONEID, SHIPTYPE, COMBINV, SHIPDTE, PC, SHIPFROM, ZONECODE, NOTE1, NOTE2, ISSUINGMAX, ISSUINGSTATUS, RECISSTYPE,UPDDTE, SYSDTE, UUID, CREATEDBY, MODIFIEDBY, REFINVOICE)VALUES('%s', to_date('%s', 'YYYY-MM-DD'), '%s', '%s', '%s', '%s', '%s', '%d', '%s', '%s', to_date('%s', 'YYYY-MM-DD'), '%s', '%s', '%s', '%s', '%s', %d, 0, '%s',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,'%s','%s', '%s', '%s')", splRef, obj.EtdDate.Format("2006-01-02"), obj.Consignee.Factory.Title, obj.Consignee.Affcode.Title, obj.Consignee.Customer.Title, obj.Consignee.Customer.Description, obj.Commercial.Title, obj.Bioabt, obj.Shipment.Title, orderType, obj.EtdDate.Format("2006-01-02"), obj.Pc.Title, obj.ShipForm, obj.ZoneCode, obj.LoadingArea, obj.Privilege, len(obj.OrderDetail), obj.Consignee.Factory.CdCode, obj.ID, empID, empID, invNo)
		if issueKey != "" {
			txtUpdate = "UPDATE"
			sqlEntExecute = fmt.Sprintf("UPDATE TXP_ISSTRANSENT SET ETDDTE=TO_DATE('%s', 'YYYY-MM-DD'),SHIPTYPE='%s',COMERCIAL='%s',PC='%s',NOTE1='%s', NOTE2='%s',ISSUINGMAX=%d, REFINVOICE='%s',UUID='%s',MODIFIEDBY='%s',UPDDTE=CURRENT_TIMESTAMP WHERE ISSUINGKEY='%s'", obj.EtdDate.Format("2006-01-02"), obj.Shipment.Title, obj.Commercial.Title, obj.Pc.Title, obj.LoadingArea, obj.Privilege, len(obj.OrderDetail), invNo, obj.ID, empID, issueKey)
		}
		_, err = db.Exec(sqlEntExecute)
		if err != nil {
			panic(err)
		}

		if issueKey == "" {
			issueKey = splRef
		}
		fmt.Printf("%d. %s ORDER ID: %s INV: %s REF: %s EmpID: %s\n", *x, txtUpdate, obj.ID, invNo, issueKey, empID)
		// Create Body
		seq := 1
		for _, v := range obj.OrderDetail {
			p := v.OrderPlan
			dbQuery := fmt.Sprintf("SELECT UUID FROM TXP_ISSTRANSBODY WHERE ISSUINGKEY='%s' AND PONO='%s' AND PARTNO='%s'", issueKey, p.Pono, p.PartNo)
			rows, err := db.Query(dbQuery)
			if err != nil {
				panic(err)
			}
			defer rows.Close()

			var bID string
			if rows.Next() {
				rows.Scan(&bID)
			}

			sqlBodyExecute := fmt.Sprintf("INSERT INTO TXP_ISSTRANSBODY(ISSUINGKEY, ISSUINGSEQ, PONO, TAGRP, PARTNO, STDPACK, ORDERQTY, ISSUINGSTATUS, BWIDE, BLENG, BHIGHT, NEWEIGHT, GTWEIGHT, UPDDTE, SYSDTE, PARTTYPE, PARTNAME, SHIPTYPE, EDTDTE, UUID, CREATEDBY, MODIFIEDBY, ORDERTYPE, LOTNO, ORDERID, REFINV)VALUES('%s',%d,'%s','C','%s',%f,%f,0,%f,%f,%f,%f,%f,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,'%s','%s','%s',TO_DATE('%s', 'YYYY-MM-DD'),'%s','%s','%s','%s','%s','%s','%s')", issueKey, seq, p.Pono, p.PartNo, p.Bistdp, p.BalQty, p.Biwidt, p.Bileng, p.Bihigh, p.Binewt, p.Bigrwt, (v.Ledger.PartType.Title)[:1], strings.ReplaceAll(p.PartName, "'", "''"), obj.Shipment.Title, p.EtdTap.Format("2006-01-02"), p.ID, empID, empID, p.OrderType.Title, p.LotNo, v.ID, invNo)
			txtIssueBody := "INSERT"
			if bID != "" {
				txtIssueBody = "UPDATE"
				sqlBodyExecute = fmt.Sprintf("UPDATE TXP_ISSTRANSBODY SET ORDERQTY=%f,NEWEIGHT=%f,GTWEIGHT=%f,UPDDTE=current_timestamp,LOTNO='%s' WHERE ISSUINGKEY='%s' AND PONO='%s' AND PARTNO='%s'", p.BalQty, p.Binewt, p.Biwidt, p.LotNo, invNo, p.Pono, p.PartNo)
			}
			fmt.Printf("==> %d.%s ISSUE: %s PONO: %s PART: %s \n", seq, txtIssueBody, issueKey, p.Pono, p.PartNo)
			_, err = db.Exec(sqlBodyExecute)
			if err != nil {
				panic(err)
			}

			sqlOrdPlan := fmt.Sprintf("UPDATE TXP_ORDERPLAN SET MODIFIEDBY='%s',CURINV='%s',OLDINV=CURINV,UPDDTE=CURRENT_TIMESTAMP WHERE ORDERID='%s'\n", empID, issueKey, p.ID)
			_, err = db.Exec(sqlOrdPlan)
			if err != nil {
				panic(err)
			}
			seq++
		}

		// Create Pallet
		seq = 1
		for _, v := range obj.Pallet {
			plNo := fmt.Sprintf("1%s%03d", v.PalletPrefix, v.PalletNo)
			dbQuery := fmt.Sprintf("SELECT ISSUINGKEY FROM TXP_ISSPALLET WHERE ISSUINGKEY='%s' AND PALLETNO='%s'", issueKey, plNo)
			rows, err := db.Query(dbQuery)
			if err != nil {
				panic(err)
			}
			defer rows.Close()

			var pID string
			if rows.Next() {
				rows.Scan(&pID)
			}

			plSqlExecute := fmt.Sprintf("INSERT INTO TXP_ISSPALLET(FACTORY, ISSUINGKEY, PALLETNO, CUSTNAME, PLTYPE, PLOUTSTS, UPDDTE, SYSDTE, PLTOTAL,PLWIDE, PLLENG, PLHIGHT)VALUES('%s','%s','%s','%s','%s',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,%d,%f,%f,%f)", obj.Consignee.Factory.Title, issueKey, plNo, obj.Consignee.Customer.Description, v.PalletPrefix, len(v.PalletDetail), v.PalletType.PalletSizeWidth, v.PalletType.PalletSizeLength, v.PalletType.BoxSizeHight)
			txtPallet := "INSERT"
			if pID != "" {
				txtPallet = "UPDATE"
				plSqlExecute = fmt.Sprintf("UPDATE TXP_ISSPALLET SET UPDDTE=current_timestamp,PLTOTAL=%d WHERE ISSUINGKEY='%s' AND PALLETNO='%s'", len(v.PalletDetail), issueKey, plNo)
			}

			_, err = db.Exec(plSqlExecute)
			if err != nil {
				panic(err)
			}

			fmt.Printf("====> %d. %s %s TOTAL: %d\n", seq, txtPallet, plNo, len(v.PalletDetail))
			// Create Issue Detail
			i := 1
			for _, p := range v.PalletDetail {
				FTicketNo := fmt.Sprintf("%s%s%08d", obj.Consignee.Factory.LabelPrefix, (obj.EtdDate.Format("2006-01-02"))[3:4], p.SeqNo)
				sqlFTicket := fmt.Sprintf("SELECT FTICKETNO FROM TXP_ISSPACKDETAIL WHERE FTICKETNO='%s'", FTicketNo)
				rows, err := db.Query(sqlFTicket)
				if err != nil {
					panic(err)
				}
				defer rows.Close()

				// Loop Issue
				var fID string
				for rows.Next() {
					rows.Scan(&fID)
				}

				strExecuteFTicketNo := fmt.Sprintf("INSERT INTO TXP_ISSPACKDETAIL(ISSUINGKEY, PONO, TAGRP, PARTNO, FTICKETNO, SHIPPLNO,ISSUINGSTATUS, UPDDTE, SYSDTE, UUID, CREATEDBY, MODIFEDBY)VALUES('%s','%s','C','%s','%s','%s',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,'%s','%s','%s')", issueKey, *p.OrderDetail.Pono, p.OrderDetail.Ledger.Part.Title, FTicketNo, plNo, p.ID, empID, empID)
				txtFTicket := "INSERT"
				if fID != "" {
					txtFTicket = "UPDATE"
					strExecuteFTicketNo = fmt.Sprintf("UPDATE TXP_ISSPACKDETAIL SET ISSUINGKEY='%s',PONO='%s',PARTNO='%s',SHIPPLNO='%s',UPDDTE=current_timestamp WHERE FTICKETNO='%s'", issueKey, *p.OrderDetail.Pono, p.OrderDetail.Ledger.Part.Title, plNo, FTicketNo)
				}
				_, err = db.Exec(strExecuteFTicketNo)
				if err != nil {
					panic(err)
				}
				fmt.Printf("======> %d. %s %s\n", i, txtFTicket, FTicketNo)
				i++
			}
			seq++
		}
	}

	// Update Order After Sync
	payload := strings.NewReader("is_sync=true")

	client := &http.Client{}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/sync/order/%s", configs.API_HOST, obj.ID), payload)

	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
}
