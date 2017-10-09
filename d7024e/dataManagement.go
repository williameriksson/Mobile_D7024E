package d7024e

import (
  "time"
  "fmt"
)

const (
  TTL time.Duration = time.Duration(60)*time.Second
  REPUBLISH_MY_FILES_INTERVAL time.Duration = time.Duration(20)*time.Second
  REPUBLISH_INTERVAL time.Duration = time.Duration(20)*time.Second
  PURGE_INTERVAL time.Duration = time.Duration(10)*time.Second
)

type PurgeInformation struct {
  Key             string
  Pinned          bool
  PurgeTimeStamp  time.Time
  TimeToLive      time.Duration
  LastPublished   time.Time
}

type DataInformation struct {
  MyKeys       map[string]bool
  PurgeInfos   map[string]PurgeInformation
}

func (kademlia *Kademlia) RepublishMyDataOnce() {
  // dataInfoMutex.Lock()
  // defer dataInfoMutex.Unlock()

  for myKey, _ := range kademlia.Datainfo.MyKeys {
    if myKey == "" {
      delete(kademlia.Datainfo.MyKeys, myKey)
      continue
    }
    tmp := kademlia.Datainfo.PurgeInfos[myKey]
    tmp.TimeToLive = TTL
    kademlia.Datainfo.PurgeInfos[myKey] = tmp
    fmt.Println("")
    fmt.Println("IN REPUBLISHMYDATA, The purgeInfo: ", kademlia.Datainfo.PurgeInfos[myKey], "The filepath: ", kademlia.files[myKey])
    fmt.Println("")
    kademlia.PublishData(kademlia.Datainfo.PurgeInfos[myKey], kademlia.files[myKey], true)
  }
}

func (kademlia *Kademlia) RepublishMyData() {
  kademlia.RepublishMyDataOnce()
  time.AfterFunc(REPUBLISH_MY_FILES_INTERVAL, kademlia.RepublishMyData)
}

func (kademlia *Kademlia) RepublishData() {
  // dataInfoMutex.Lock()

  for _, purgeInfo := range kademlia.Datainfo.PurgeInfos {
    if purgeInfo.Key == "" {
      continue
    }
    purgeInfo.TimeToLive = purgeInfo.PurgeTimeStamp.Sub(time.Now())

    filesMutex.Lock()
    if _, exists := kademlia.files[purgeInfo.Key]; exists {
      fmt.Println("\n IN REPUBLISH DATA: exists in kademlia.files, the key is: ", purgeInfo.Key, "\n")
      if !kademlia.Datainfo.MyKeys[purgeInfo.Key] {
        if (time.Now().Sub(purgeInfo.LastPublished) > REPUBLISH_INTERVAL ) {
          purgeInfo.LastPublished = time.Now()
          kademlia.PublishData(purgeInfo, kademlia.files[purgeInfo.Key], false)
        } else {
          fmt.Println("NOPE HAS ALDREADY BEEN REPUBLISHED, the diff is:", time.Now().Sub(purgeInfo.LastPublished)  )
        }
      }
  	} else {
      fmt.Println("Wanted to Republish this key: ", purgeInfo.Key, " but the key was not found in kademlia.files.")
    }


    filesMutex.Unlock()
  }
  // dataInfoMutex.Unlock()
  time.AfterFunc(REPUBLISH_INTERVAL, kademlia.RepublishData)
}


// Should periodically call itself (needs testing), could be changed to trigger on next event
// if sorting mechanism is implemented
func (kademlia *Kademlia) PurgeData() {
  dataInfoMutex.Lock()
  for key, purgeInfo := range kademlia.Datainfo.PurgeInfos {
    if !purgeInfo.Pinned && time.Now().After(purgeInfo.PurgeTimeStamp) && !kademlia.Datainfo.MyKeys[purgeInfo.Key]{
      filesMutex.Lock()
      select {
      case kademlia.ServerChannel <- NewHandle(CMD_REMOVE_FILE, purgeInfo, kademlia.files[purgeInfo.Key]):
        delete(kademlia.files, purgeInfo.Key)
        delete(kademlia.Datainfo.PurgeInfos, key)
        delete(kademlia.Datainfo.MyKeys, key)
      case <-time.After(time.Millisecond * 50):
        fmt.Println("Could not purge the data, handler did not read the channel")
      }
      filesMutex.Unlock()
    }
  }
  dataInfoMutex.Unlock()
  time.AfterFunc(PURGE_INTERVAL, kademlia.PurgeData)

}

func (kademlia *Kademlia) SetPurgeStamp(purgeInfo *PurgeInformation) {
  newTimeStamp := time.Now().Add(purgeInfo.TimeToLive)
  if newTimeStamp.After(purgeInfo.PurgeTimeStamp) {
    purgeInfo.PurgeTimeStamp = newTimeStamp
  }
}
