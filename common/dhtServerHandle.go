package common

const (
  CMD_FOUND_FILE = "FOUND_FILE"
  CMD_RETRIEVE_FILE = "RETRIEVE_FILE"
)

type Handle struct {
  command string
  hash  string
  ip    string
}

func NewHandle(cmd string, hash string, ip string) Handle {
  return Handle{cmd, hash, ip}
}
