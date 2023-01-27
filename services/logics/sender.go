package logics

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"user/domains/models"
	"user/onesender"
	"user/services/repositories"
	"user/utils"

	"github.com/sirupsen/logrus"
)

type ISenderService interface {
	GetSender(req models.Sender)
}

type SenderService struct {
	senderRepository repositories.ISenderRepository
	useTransaction   bool
}

func InitSenderService(senderRepo repositories.ISenderRepository) *SenderService {
	if utils.IsNil(senderRepo) {
		senderRepo = repositories.InitSenderRepository(nil, nil)
	}

	service := SenderService{
		senderRepository: senderRepo,
		useTransaction:   false,
	}
	return &service
}

func (service *SenderService) GetSender(req models.Sender) {
	onesender.ApiUrl = os.Getenv("API_URL_ONESENDER")
	onesender.ApiKey = os.Getenv("API_KEY_ONESENDER")
	// onesender.SendTextMessage(req.Sender_Phone, "Selamat Datang di BOT Garasi Pulsa yang menjual Pulsa, Pulsa Data, Token Listrik, PDAM, Dan Lain Lain, untuk melihat menu anda dapat mengetikan kata /Menu atau /menu")
	if strings.Contains(req.Message_Text, "/Menu") || strings.Contains(req.Message_Text, "/menu") {

		jsonData := models.APISenderWithButton{}
		jsonData.Recipient_type = "individual"
		jsonData.To = req.Sender_Phone
		jsonData.Type = "interactive"
		jsonData.Interactive.Type = "button"
		jsonData.Interactive.Header.Text = "Layanan"
		jsonData.Interactive.Body.Text = "Menu Layanan"
		jsonData.Interactive.Footer.Text = "Pilihan Layanan"

		// button 1
		button := models.Button{}
		button.Type = "button"
		button.Reply.ID = "prabayar"
		button.Reply.Title = "/prabayar"
		jsonData.Interactive.Action.Buttons = append(jsonData.Interactive.Action.Buttons, button)
		// button 2
		button2 := models.Button{}
		button2.Type = "button"
		button2.Reply.ID = "paskabayar"
		button2.Reply.Title = "/paskaBayar"
		jsonData.Interactive.Action.Buttons = append(jsonData.Interactive.Action.Buttons, button2)

		bearer := fmt.Sprintf("Bearer %s", os.Getenv("API_KEY_ONESENDER"))
		headers := map[string]string{
			"Authorization": bearer,
		}
		result, err := utils.CallAPI(http.MethodPost, os.Getenv("API_URL_ONESENDER"), &jsonData, headers, nil)
		if err != nil {
			utils.PrintLog("error [services][logics][sender][CallAPI] ", err)
			logrus.Error("error [services][logics][sender][CallAPI] ", err)
		}
		defer result.Body.Close()
		bytes, err := io.ReadAll(result.Body)
		if err != nil {
			utils.PrintLog("error [services][logics][sender][ReadAll Looping CallAPI] ", err)
			logrus.Error("error [services][logics][sender][ReadAll Looping CallAPI] ", err)
		}
		res := models.APIReceiver{}
		err = json.Unmarshal(bytes, &res)
		if err != nil {
			utils.PrintLog("error [services][logics][sender][ReadAll Looping CallAPI] ", err)
			logrus.Error("error [services][logics][sender][ReadAll Looping CallAPI] ", err)
		}
		fmt.Println(res)
	} else if strings.Contains(req.Message_Text, "/prabayar") {
		url := "https://api.digiflazz.com/v1/price-list"
		jsonData := models.CheckHargaDigiflazz{}
		username := os.Getenv("USERNAME_BELI_PULSA")
		apikey := os.Getenv("API_KEY_BELI_PULSA")
		sign := md5.Sum([]byte(username + apikey + "pricelist"))
		pass := fmt.Sprintf("%x", sign)
		jsonData.Command = "prepaid"
		jsonData.Username = username
		jsonData.Sign = pass
		result, err := utils.CallAPI(http.MethodPost, url, &jsonData, nil, nil)
		if err != nil {
			utils.PrintLog("error [services][logics][sender][CallAPI] ", err)
			logrus.Error("error [services][logics][sender][CallAPI] ", err)
		}
		defer result.Body.Close()
		bytes, err := io.ReadAll(result.Body)
		if err != nil {
			utils.PrintLog("error [services][logics][sender][ReadAll Looping CallAPI] ", err)
			logrus.Error("error [services][logics][sender][ReadAll Looping CallAPI] ", err)
		}
		res := models.DaftarHarga{}
		// messageToWa := []models.SenderToWa{}
		err = json.Unmarshal(bytes, &res)
		if err != nil {
			utils.PrintLog("error [services][logics][sender][ReadAll Looping CallAPI] ", err)
			logrus.Error("error [services][logics][sender][ReadAll Looping CallAPI] ", err)
		}
		jsonData2 := models.APISenderWithButton{}
		jsonData2.Recipient_type = "individual"
		jsonData2.To = req.Sender_Phone
		jsonData2.Type = "interactive"
		jsonData2.Interactive.Type = "list"
		jsonData2.Interactive.Header.Text = "Daftar Produk"
		jsonData2.Interactive.Body.Text = "Berikut daftar produk yang tersedia untuk di jual"
		jsonData2.Interactive.Footer.Text = "Pilihlah produk yang akan di beli"
		jsonData2.Interactive.Action.Button = "Daftar Produk"

		section1 := models.Section{}
		section1.Title = "PLN"
		section2 := models.Section{}
		section2.Title = "Pulsa"
		section3 := models.Section{}
		section3.Title = "Data"
		section4 := models.Section{}
		section4.Title = "Masa Aktif"
		section5 := models.Section{}
		section5.Title = "E-Money"
		row1 := []models.Row{}
		row2 := []models.Row{}
		row3 := []models.Row{}
		row4 := []models.Row{}
		row5 := []models.Row{}
		for _, data := range res.Data {
			if data.Category == "PLN" {
				row1 = append(row1, models.Row{ID: data.Buyer_SKU_Code, Title: data.Product_Name, Description: "/" + data.Buyer_SKU_Code})
				section1.Rows = row1
			}
			if data.Category == "Pulsa" {
				row2 = append(row2, models.Row{ID: data.Buyer_SKU_Code, Title: data.Product_Name, Description: "/" + data.Buyer_SKU_Code})
				section2.Rows = row2
			}
			if data.Category == "Data" {
				row3 = append(row3, models.Row{ID: data.Buyer_SKU_Code, Title: data.Product_Name, Description: "/" + data.Buyer_SKU_Code})
				section3.Rows = row3
			}
			if data.Category == "Masa Aktif" {
				row4 = append(row4, models.Row{ID: data.Buyer_SKU_Code, Title: data.Product_Name, Description: "/" + data.Buyer_SKU_Code})
				section4.Rows = row4
			}
			if data.Category == "E-Money" {
				row5 = append(row5, models.Row{ID: data.Buyer_SKU_Code, Title: data.Product_Name, Description: "/" + data.Buyer_SKU_Code})
				section5.Rows = row5
			}
		}
		jsonData2.Interactive.Action.Sections = append(jsonData2.Interactive.Action.Sections, section1)
		jsonData2.Interactive.Action.Sections = append(jsonData2.Interactive.Action.Sections, section2)
		jsonData2.Interactive.Action.Sections = append(jsonData2.Interactive.Action.Sections, section3)
		jsonData2.Interactive.Action.Sections = append(jsonData2.Interactive.Action.Sections, section4)
		jsonData2.Interactive.Action.Sections = append(jsonData2.Interactive.Action.Sections, section5)

		bearer := fmt.Sprintf("Bearer %s", os.Getenv("API_KEY_ONESENDER"))
		headers := map[string]string{
			"Authorization": bearer,
		}

		result2, err2 := utils.CallAPI(http.MethodPost, os.Getenv("API_URL_ONESENDER"), &jsonData2, headers, nil)
		if err2 != nil {
			utils.PrintLog("error [services][logics][sender][CallAPI] ", err2)
			logrus.Error("error [services][logics][sender][CallAPI] ", err2)
		}
		defer result.Body.Close()
		bytes, err2 = io.ReadAll(result2.Body)
		if err != nil {
			utils.PrintLog("error [services][logics][sender][ReadAll Looping CallAPI] ", err2)
			logrus.Error("error [services][logics][sender][ReadAll Looping CallAPI] ", err2)
		}
		res2 := models.APIReceiver{}
		err2 = json.Unmarshal(bytes, &res2)
		if err2 != nil {
			utils.PrintLog("error [services][logics][sender][ReadAll Looping CallAPI] ", err2)
			logrus.Error("error [services][logics][sender][ReadAll Looping CallAPI] ", err2)
		}
		fmt.Println(res2)
	} else if strings.Contains(req.Message_Text, "/PLN20") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan PLN20<spasi><nomor pelanggan pln anda> contoh PLN20 2343039230394923049")
	} else if strings.Contains(req.Message_Text, "/PLN50") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan PLN50<spasi><nomor pelanggan pln anda> contoh PLN50 2343039230394923049")
	} else if strings.Contains(req.Message_Text, "/PLN100") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan PLN100<spasi><nomor pelanggan pln anda> contoh PLN100 2343039230394923049")
	} else if strings.Contains(req.Message_Text, "/TELKOM10") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan TELKOM10<spasi><nomor pelanggan telkomsel anda> contoh TELKOM10 081392381293")
	} else if strings.Contains(req.Message_Text, "/TELKOM20") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan TELKOM20<spasi><nomor pelanggan telkomsel anda> contoh TELKOM20 081392381293")
	} else if strings.Contains(req.Message_Text, "/TELKOM25") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan TELKOM25<spasi><nomor pelanggan telkomsel anda> contoh TELKOM25 081392381293")
	} else if strings.Contains(req.Message_Text, "/TELKOM50") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan TELKOM50<spasi><nomor pelanggan telkomsel anda> contoh TELKOM50 081392381293")
	} else if strings.Contains(req.Message_Text, "/TELKOM100") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan TELKOM100<spasi><nomor pelanggan telkomsel anda> contoh TELKOM100 081392381293")
	} else if strings.Contains(req.Message_Text, "/TRI10") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan TRI10<spasi><nomor pelanggan three anda> contoh TRI10 0898123912812")
	} else if strings.Contains(req.Message_Text, "/TRI20") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan TRI20<spasi><nomor pelanggan three anda> contoh TRI20 0898123912812")
	} else if strings.Contains(req.Message_Text, "/TRI25") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan TRI25<spasi><nomor pelanggan three anda> contoh TRI25 0898123912812")
	} else if strings.Contains(req.Message_Text, "/TRI50") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan TRI50<spasi><nomor pelanggan three anda> contoh TRI50 0898123912812")
	} else if strings.Contains(req.Message_Text, "/TRI100") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan TRI100<spasi><nomor pelanggan three anda> contoh TRI100 0898123912812")
	} else if strings.Contains(req.Message_Text, "/IND10") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan IND10<spasi><nomor pelanggan indosat anda> contoh IND10 0814294812398")
	} else if strings.Contains(req.Message_Text, "/IND20") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan IND20<spasi><nomor pelanggan indosat anda> contoh IND20 0814294812398")
	} else if strings.Contains(req.Message_Text, "/IND25") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan IND25<spasi><nomor pelanggan indosat anda> contoh IND25 0814294812398")
	} else if strings.Contains(req.Message_Text, "/IND50") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan IND50<spasi><nomor pelanggan indosat anda> contoh IND50 0814294812398")
	} else if strings.Contains(req.Message_Text, "/SMRTF10") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan SMRTF10<spasi><nomor pelanggan smartfren anda> contoh SMRTF10 08819239198328")
	} else if strings.Contains(req.Message_Text, "/SMRTF20") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan SMRTF20<spasi><nomor pelanggan smartfren anda> contoh SMRTF20 08819239198328")
	} else if strings.Contains(req.Message_Text, "/SMRTF50") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan SMRTF50<spasi><nomor pelanggan smartfren anda> contoh SMRTF50 08819239198328")
	} else if strings.Contains(req.Message_Text, "/SMRTF100") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan SMRTF100<spasi><nomor pelanggan smartfren anda> contoh SMRTF100 08819239198328")
	} else if strings.Contains(req.Message_Text, "/AXIS10") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan AXIS10<spasi><nomor pelanggan axis anda> contoh AXIS10 0838102391029")
	} else if strings.Contains(req.Message_Text, "/AXIS25") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan AXIS25<spasi><nomor pelanggan axis anda> contoh AXIS25 0838102391029")
	} else if strings.Contains(req.Message_Text, "/AXIS50") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan AXIS50<spasi><nomor pelanggan axis anda> contoh AXIS50 0838102391029")
	} else if strings.Contains(req.Message_Text, "/XL10") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan XL10<spasi><nomor pelanggan xl anda> contoh XL10 08179123128312")
	} else if strings.Contains(req.Message_Text, "/XL25") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan XL25<spasi><nomor pelanggan xl anda> contoh XL25 08179123128312")
	} else if strings.Contains(req.Message_Text, "/XL50") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan XL50<spasi><nomor pelanggan xl anda> contoh XL50 08179123128312")
	} else if strings.Contains(req.Message_Text, "/TELKOMDATA1GB30HARI") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan TELKOMDATA1GB30HARI<spasi><nomor pelanggan telkomsel anda> contoh TELKOMDATA1GB30HARI 0813239189289")
	} else if strings.Contains(req.Message_Text, "/TELKOMDATA3GB30HARI") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan TELKOMDATA3GB30HARI<spasi><nomor pelanggan telkomsel anda> contoh TELKOMDATA3GB30HARI 0813239189289")
	} else if strings.Contains(req.Message_Text, "/TELKOMDATA25GB30HARI") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan TELKOMDATA25GB30HARI<spasi><nomor pelanggan telkomsel anda> contoh TELKOMDATA25GB30HARI 0813239189289")
	} else if strings.Contains(req.Message_Text, "/INDODATA1GB30HARI") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan INDODATA1GB30HARI<spasi><nomor pelanggan indosat anda> contoh INDODATA1GB30HARI 0814912930122")
	} else if strings.Contains(req.Message_Text, "/INDODATA2GB30HARI") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan INDODATA2GB30HARI<spasi><nomor pelanggan indosat anda> contoh INDODATA2GB30HARI 0814912930122")
	} else if strings.Contains(req.Message_Text, "/INDODATA3GB30HARI") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan INDODATA3GB30HARI<spasi><nomor pelanggan indosat anda> contoh INDODATA3GB30HARI 0814912930122")
	} else if strings.Contains(req.Message_Text, "/INDODATA6GB30HARI") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan INDODATA6GB30HARI<spasi><nomor pelanggan indosat anda> contoh INDODATA6GB30HARI 0814912930122")
	} else if strings.Contains(req.Message_Text, "/INDODATA10GB30HARI") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan INDODATA10GB30HARI<spasi><nomor pelanggan indosat anda> contoh INDODATA10GB30HARI 0814912930122")
	} else if strings.Contains(req.Message_Text, "/AXISDATA300MB7HARI") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan AXISDATA300MB7HARI<spasi><nomor pelanggan axis anda> contoh AXISDATA300MB7HARI 0838192391291")
	} else if strings.Contains(req.Message_Text, "/AXISDATA1GB30HARI") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan AXISDATA1GB30HARIspasi><nomor pelanggan axis anda> contoh AXISDATA1GB30HARI 0838192391291")
	} else if strings.Contains(req.Message_Text, "/AXISDATA2GB30HARI") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan AXISDATA2GB30HARI<spasi><nomor pelanggan axis anda> contoh AXISDATA2GB30HARI 0838192391291")
	} else if strings.Contains(req.Message_Text, "/AXISDATA5GB30HARI") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan AXISDATA5GB30HARI<spasi><nomor pelanggan axis anda> contoh AXISDATA5GB30HARI 0838192391291")
	} else if strings.Contains(req.Message_Text, "/AXISDATA20GB30HARI") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan AXISDATA20GB30HARI<spasi><nomor pelanggan axis anda> contoh AXISDATA20GB30HARI 0838192391291")
	} else if strings.Contains(req.Message_Text, "/AXISDATA50GB30HARI") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan AXISDATA50GB30HARI<spasi><nomor pelanggan axis anda> contoh AXISDATA50GB30HARI 0838192391291")
	} else if strings.Contains(req.Message_Text, "/TRI4G10GB30HARI") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan TRI4G10GB30HARI<spasi><nomor pelanggan three anda> contoh TRI4G10GB30HARI 089812318238128")
	} else if strings.Contains(req.Message_Text, "/TRI1GB7HARI") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan TRI1GB7HARI<spasi><nomor pelanggan three anda> contoh TRI1GB7HARI 089812318238128")
	} else if strings.Contains(req.Message_Text, "/TRI1GB14HARI") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan TRI1GB14HARI<spasi><nomor pelanggan three anda> contoh TRI1GB14HARI 089812318238128")
	} else if strings.Contains(req.Message_Text, "/TRIDATA2GB30HARI") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan TRIDATA2GB30HARI<spasi><nomor pelanggan three anda> contoh TRIDATA2GB30HARI 089812318238128")
	} else if strings.Contains(req.Message_Text, "/TRI10GBSEMUA10GByoutubene") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan TRI10GBSEMUA10GByoutubene<spasi><nomor pelanggan three anda> contoh TRI10GBSEMUA10GByoutubene 089812318238128")
	} else if strings.Contains(req.Message_Text, "/TELKOM30HARIAKTIF") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan TELKOM30HARIAKTIF<spasi><nomor pelanggan telkomsel anda> contoh TELKOM30HARIAKTIF 08131923912283")
	} else if strings.Contains(req.Message_Text, "/GOPAY10") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan GOPAY10<spasi><nomor pelanggan gopay anda> contoh GOPAY10 0813123123123")
	} else if strings.Contains(req.Message_Text, "/GOPAY20") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan GOPAY20<spasi><nomor pelanggan gopay anda> contoh GOPAY20 0813123123123")
	} else if strings.Contains(req.Message_Text, "/GOPAY50") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan GOPAY50<spasi><nomor pelanggan gopay anda> contoh GOPAY50 0813123123123")
	} else if strings.Contains(req.Message_Text, "/GOPAY100") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan GOPAY100<spasi><nomor pelanggan gopay anda> contoh GOPAY100 0813123123123")
	} else if strings.Contains(req.Message_Text, "/OVO20") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan OVO20<spasi><nomor pelanggan ovo anda> contoh OVO20 0813123123123")
	} else if strings.Contains(req.Message_Text, "/OVO50") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan OVO50<spasi><nomor pelanggan ovo anda> contoh OVO50 0813123123123")
	} else if strings.Contains(req.Message_Text, "/OVO100") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan OVO100<spasi><nomor pelanggan ovo anda> contoh OVO100 0813123123123")
	} else if strings.Contains(req.Message_Text, "/DANA10") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan DANA10<spasi><nomor pelanggan dana anda> contoh DANA10 0813123123123")
	} else if strings.Contains(req.Message_Text, "/DANA20") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan DANA20<spasi><nomor pelanggan dana anda> contoh DANA20 0813123123123")
	} else if strings.Contains(req.Message_Text, "/DANA100") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan DANA100<spasi><nomor pelanggan dana anda> contoh DANA100 0813123123123")
	} else if strings.Contains(req.Message_Text, "/LINKAJA10") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan LINKAJA10<spasi><nomor pelanggan linkaja anda> contoh LINKAJA10 0813123123123")
	} else if strings.Contains(req.Message_Text, "/LINKAJA20") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan LINKAJA20<spasi><nomor pelanggan linkaja anda> contoh LINKAJA20 0813123123123")
	} else if strings.Contains(req.Message_Text, "/LINKAJA50") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan LINKAJA50<spasi><nomor pelanggan linkaja anda> contoh LINKAJA50 0813123123123")
	} else if strings.Contains(req.Message_Text, "/LINKAJA100") {
		onesender.SendTextMessage(req.Sender_Phone, "ketikan LINKAJA100<spasi><nomor pelanggan linkaja anda> contoh LINKAJA100 0813123123123")
	}
}
