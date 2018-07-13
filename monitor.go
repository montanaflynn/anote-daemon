package main

import (
	"log"
	"time"

	"github.com/anonutopia/gowaves"
)

type WavesMonitor struct {
	StartedTime int64
}

func (wm *WavesMonitor) start() {
	wm.StartedTime = time.Now().Unix() * 1000
	for {
		// todo - make sure that everything is ok with 10 here
		pages, err := wnc.TransactionsAddressLimit("3PDb1ULFjazuzPeWkF2vqd1nomKh4ctq9y2", 100)
		if err != nil {
			log.Println(err)
		}
		if len(pages) > 0 {
			for _, t := range pages[0] {
				wm.checkTransaction(&t)
			}
		}
		time.Sleep(time.Second)
	}
}

func (wm *WavesMonitor) checkTransaction(t *gowaves.TransactionsAddressLimitResponse) {
	tr := Transaction{TxId: t.ID}
	db.FirstOrCreate(&tr, &tr)
	if tr.Processed != 1 {
		wm.processTransaction(&tr, t)
	}
}

func (wm *WavesMonitor) processTransaction(tr *Transaction, t *gowaves.TransactionsAddressLimitResponse) {
	if t.Type == 4 && t.Timestamp >= wm.StartedTime {
		if len(t.AssetID) == 0 {
			p, err := pc.DoRequest()
			if err == nil {
				amount := int(float64(t.Amount) * p.WAVES / 0.01)
				atr := &gowaves.AssetsTransferRequest{
					Amount:    amount,
					AssetID:   "4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf",
					Fee:       100000,
					Recipient: t.Sender,
					Sender:    conf.NodeAddress,
				}

				_, err := wnc.AssetsTransfer(atr)
				if err != nil {
					log.Printf("[WavesMonitor.processTransation] error assets transfer: %s", err)
				} else {
					log.Printf("Sent ANO: %s => %d", t.Sender, amount)
				}

				user := &User{Address: t.Sender}
				db.First(user, user)
				if len(user.Referral) > 0 {
					referral := &User{Address: user.Referral}
					db.First(referral, referral)
					splitToHolders := t.Amount / 2
					splitToHolders -= (t.Amount / 5)
					if referral.ID != 0 {
						newProfit := uint64(t.Amount / 5)
						referral.ReferralProfitWav += newProfit
						referral.ReferralProfitWavTotal += newProfit
						db.Save(referral)
						splitToHolders -= (t.Amount / 5)
					}
					log.Println(splitToHolders)
				}
			}
		}
	}

	tr.Processed = 1
	db.Save(tr)
}

func (wm *WavesMonitor) calculateAmount(trType int64, amount int64) int64 {
	var amountSending int64
	amountSending = 0
	return amountSending
}

func initMonitor() {
	wm := &WavesMonitor{}
	wm.start()
}
