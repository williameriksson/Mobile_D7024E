package common

const (
  CMD_FOUND_FILE = "FOUND_FILE"
  CMD_RETRIEVE_FILE = "RETRIEVE_FILE"
  CMD_REMOVE_FILE = "REMOVE_FILE"
)

type Handle struct {
  Command string
  Hash  string
  Ip    string
}

func NewHandle(cmd string, hash string, ip string) Handle {
  return Handle{cmd, hash, ip}
}
