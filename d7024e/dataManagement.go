package d7024e

import (
  "time"
)

type PurgeInformation struct {
  Key             string
  Pinned          bool
  PurgeTimeStamp  time.Time
}

type DataInformation struct {
  MyKeys       []string
  PurgeInfos   []PurgeInformation
}

func (kademlia *Kademlia) RepublishData() {

}

func (kademlia *Kademlia) PurgeData() {

}

func (kademlia *Kademlia) SetPurgeStamp(purgeInfo PurgeInformation) {
  duration := time.Duration(60)*time.Second
  purgeInfo.PurgeTimeStamp = time.Now().Add(duration)
}
