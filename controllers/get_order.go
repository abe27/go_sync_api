package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/abe27/syncapi/configs"
	"github.com/abe27/syncapi/models"
	_ "gopkg.in/goracle.v2"
)

func FetchTest() {
	fmt.Println(configs.ORAC_DNS)
	db, err := sql.Open("goracle", configs.ORAC_DNS)
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func GenerateInvNo(obj models.Order) string {
	return fmt.Sprintf("%s%s%s%04d%s", obj.Consignee.Factory.InvPrefix, obj.Consignee.Prefix, (strconv.Itoa(obj.EtdDate.Year()))[3:], obj.RunningSeq, obj.Shipment.Title)
}

func FetchAll() (*models.ResponseOrder, error) {
	var data models.ResponseOrder
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/sync/order", configs.API_HOST), nil)
	if err != nil {
		return &data, err
	}

	res, err := client.Do(req)
	if err != nil {
		return &data, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &data, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return &data, err
	}
	fmt.Printf("Found Order Size: %d\n", len(data.Data))
	return &data, nil
}

func CreateIssueEnt(data *[]models.Order) (bool, error) {
	fmt.Println(configs.ORAC_DNS)
	db, err := sql.Open("goracle", configs.ORAC_DNS)
	if err != nil {
		return false, err
	}
	defer db.Close()

	runn := 1
	for _, v := range *data {
		// fmt.Println(v.ID)
		dbQuery := fmt.Sprintf("SELECT UUID FROM TXP_ISSTRANSENT WHERE UUID='%s'", v.ID)
		rows, err := db.Query(dbQuery)
		if err != nil {
			fmt.Println(".....Error processing query")
			fmt.Println(err)
			return false, err
		}
		defer rows.Close()

		// Loop Issue
		var uuID string
		for rows.Next() {
			rows.Scan(&uuID)
		}

		empID := "SKTSYS"
		if len(v.Consignee.OrderGroup) > 0 {
			empID = v.Consignee.OrderGroup[0].User.UserName
		}
		invNo := GenerateInvNo(v)
		strExecuteOrderEnt := fmt.Sprintf("INSERT INTO TXP_ISSTRANSENT(ISSUINGKEY, ETDDTE, FACTORY, AFFCODE, BISHPC, CUSTNAME, COMERCIAL, ZONEID, SHIPTYPE, COMBINV, SHIPDTE, PC, SHIPFROM, ZONECODE, NOTE1, NOTE2, ISSUINGMAX, ISSUINGSTATUS, RECISSTYPE,UPDDTE, SYSDTE, UUID, CREATEDBY, MODIFIEDBY, REFINVOICE)VALUES('%s', to_date('%s', 'YYYY-MM-DD'), '%s', '%s', '%s', '%s', '%s', '%d', '%s', '%s', to_date('%s', 'YYYY-MM-DD'), '%s', '%s', '%s', '%s', '%s', %d, 0, '%s',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,'%s','%s', '%s', '%s')", invNo, v.EtdDate.Format("2006-01-02"), v.Consignee.Factory.Title, v.Consignee.Affcode.Title, v.Consignee.Customer.Title, v.Consignee.Customer.Description, v.Commercial.Title, v.Bioabt, v.Shipment.Title, v.OrderDetail[0].OrderPlan.OrderType.Title, v.EtdDate.Format("2006-01-02"), v.Pc.Title, v.ShipForm, v.ZoneCode, v.LoadingArea, v.Privilege, len(v.OrderDetail), v.Consignee.Factory.CdCode, v.ID, empID, empID, invNo)
		if uuID != "" {
			strExecuteOrderEnt = fmt.Sprintf("UPDATE TXP_ISSTRANSENT SET ETDDTE=to_date('%s', 'YYYY-MM-DD'),SHIPTYPE='%s',COMERCIAL='%s',PC='%s',NOTE1='%s', NOTE2='%s',ISSUINGMAX=%d, UPDDTE=CURRENT_TIMESTAMP WHERE UUID='%s'", v.EtdDate.Format("2006-01-02"), v.Shipment.Title, v.Commercial.Title, v.Pc.Title, v.LoadingArea, v.Privilege, len(v.OrderDetail), uuID)
		}

		_, err = db.Exec(strExecuteOrderEnt)
		if err != nil {
			fmt.Println(err.Error())
			return false, err
		}
		// fmt.Println(strExecuteOrderEnt)

		// Order Detail
		seq := 1
		for _, b := range v.OrderDetail {
			// Fetch Body Data
			p := b.OrderPlan
			dbQuery := fmt.Sprintf("SELECT UUID FROM TXP_ISSTRANSBODY WHERE ISSUINGKEY='%s' AND PONO='%s' AND PARTNO='%s'", invNo, p.Pono, p.PartNo)
			rows, err := db.Query(dbQuery)
			if err != nil {
				fmt.Println(".....Error processing query")
				fmt.Println(err)
				return false, err
			}
			defer rows.Close()

			// Loop Issue
			var bodyID string
			for rows.Next() {
				rows.Scan(&bodyID)
			}

			strExecuteOrderBody := fmt.Sprintf("INSERT INTO TXP_ISSTRANSBODY(ISSUINGKEY, ISSUINGSEQ, PONO, TAGRP, PARTNO, STDPACK, ORDERQTY, ISSUINGSTATUS, BWIDE, BLENG, BHIGHT, NEWEIGHT, GTWEIGHT, UPDDTE, SYSDTE, PARTTYPE, PARTNAME, SHIPTYPE, EDTDTE, UUID, CREATEDBY, MODIFIEDBY, ORDERTYPE, LOTNO, ORDERID, REFINV)VALUES('%s',%d,'%s','C','%s',%f,%f,0,%f,%f,%f,%f,%f,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,'%s','%s','%s',TO_DATE('%s', 'YYYY-MM-DD'),'%s','%s','%s','%s','%s','%s','%s')", invNo, seq, p.Pono, p.PartNo, p.Bistdp, p.BalQty, p.Biwidt, p.Bileng, p.Bihigh, p.Binewt, p.Bigrwt, (b.Ledger.PartType.Title)[:1], strings.ReplaceAll(p.PartName, "'", "''"), v.Shipment.Title, p.EtdTap.Format("2006-01-02"), p.ID, empID, empID, p.OrderType.Title, p.LotNo, v.ID, invNo)
			if bodyID != "" {
				strExecuteOrderBody = fmt.Sprintf("UPDATE TXP_ISSTRANSBODY SET ORDERQTY=%f,NEWEIGHT=%f,GTWEIGHT=%f,UPDDTE=current_timestamp,LOTNO='%s' WHERE ISSUINGKEY='%s' AND PONO='%s' AND PARTNO='%s'", p.BalQty, p.Binewt, p.Biwidt, p.LotNo, invNo, p.Pono, strings.ReplaceAll(p.PartName, "'", "''"))
			}
			// fmt.Println(strExecuteOrderBody)
			_, err = db.Exec(strExecuteOrderBody)
			if err != nil {
				fmt.Println(err.Error())
				return false, err
			}

			// Update OrderPlan
			_, err = db.Exec(fmt.Sprintf("UPDATE TXP_ORDERPLAN SET MODIFIEDBY='%s',CURINV='%s',OLDINV=CURINV,UPDDTE=CURRENT_TIMESTAMP WHERE ORDERID='%s'\n", empID, invNo, p.ID))
			if err != nil {
				fmt.Println(err.Error())
				return false, err
			}
			seq++
		}

		// Order Pallet
		for _, p := range v.Pallet {
			pNum := fmt.Sprintf("%s%03d", p.PalletPrefix, p.PalletNo)
			dbQuery := fmt.Sprintf("SELECT ISSUINGKEY FROM TXP_ISSPALLET WHERE ISSUINGKEY='%s' AND PALLETNO='%s'", invNo, pNum)
			rows, err := db.Query(dbQuery)
			if err != nil {
				fmt.Println(err.Error())
				return false, err
			}
			defer rows.Close()

			// Loop Issue
			var uuID string
			for rows.Next() {
				rows.Scan(&uuID)
			}

			strExecuteOrderPallet := fmt.Sprintf("INSERT INTO TXP_ISSPALLET(FACTORY, ISSUINGKEY, PALLETNO, CUSTNAME, PLTYPE, PLOUTSTS, UPDDTE, SYSDTE, PLTOTAL,PLWIDE, PLLENG, PLHIGHT)VALUES('%s','%s','%s','%s','%s',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,%d,%f,%f,%f)", v.Consignee.Factory.Title, invNo, pNum, v.Consignee.Customer.Description, p.PalletPrefix, p.PalletTotal, p.PalletType.PalletSizeWidth, p.PalletType.PalletSizeLength, p.PalletType.BoxSizeHight)
			if uuID != "" {
				strExecuteOrderPallet = fmt.Sprintf("UPDATE TXP_ISSPALLET SET UPDDTE=current_timestamp,PLTOTAL=%d WHERE ISSUINGKEY='%s' AND PALLETNO='%s'", p.PalletTotal, invNo, pNum)
			}
			_, err = db.Exec(strExecuteOrderPallet)
			if err != nil {
				fmt.Println(err.Error())
				return false, err
			}

			// Loop Issue Detail
			for _, pl := range p.PalletDetail {
				fticket_no := fmt.Sprintf("%s%s%08d", v.Consignee.Factory.LabelPrefix, (v.EtdDate.Format("2006-01-02"))[3:4], pl.SeqNo)
				rows, err := db.Query(fmt.Sprintf("SELECT FTICKETNO FROM TXP_ISSPACKDETAIL WHERE FTICKETNO='%s'", fticket_no))
				if err != nil {
					fmt.Println(err.Error())
					return false, err
				}
				defer rows.Close()

				// Loop Issue
				var fID string
				for rows.Next() {
					rows.Scan(&fID)
				}

				strExecuteFTicketNo := fmt.Sprintf("INSERT INTO TXP_ISSPACKDETAIL(ISSUINGKEY, PONO, TAGRP, PARTNO, FTICKETNO, SHIPPLNO,ISSUINGSTATUS, UPDDTE, SYSDTE, UUID, CREATEDBY, MODIFEDBY)VALUES('%s','%s','C','%s','%s','%s',0,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,'%s','%s','%s')", invNo, *pl.OrderDetail.Pono, pl.OrderDetail.Ledger.Part.Title, fticket_no, pNum, pl.ID, empID, empID)
				if fID != "" {
					strExecuteFTicketNo = fmt.Sprintf("UPDATE TXP_ISSPACKDETAIL SET UPDDTE=current_timestamp WHERE FTICKETNO='%s'", fticket_no)
				}
				_, err = db.Exec(strExecuteFTicketNo)
				if err != nil {
					fmt.Println(err.Error())
					return false, err
				}
			}
		}

		fmt.Printf("%d SYNC ORDER ID: %s\n", runn, v.ID)

		payload := strings.NewReader("is_sync=false")

		client := &http.Client{}
		req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/sync/order/%s", configs.API_HOST, v.ID), payload)

		if err != nil {
			return false, err
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		res, err := client.Do(req)
		if err != nil {
			return false, err
		}
		defer res.Body.Close()

		_, err = io.ReadAll(res.Body)
		if err != nil {
			return false, err
		}
		// fmt.Println(string(body))
		runn++
	}

	return true, nil
}
