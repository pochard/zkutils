package zkutils

import (
	"errors"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
)

func NewRegister(conn *zk.Conn, path string) *Register {
	return &Register{conn, path}
}

type Register struct {
	conn *zk.Conn
	path string
}

func (r *Register) Create() error {
	_, err := r.conn.Create(r.path, []byte{}, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	return err
}

func (r *Register) Delete() error {
	_, s, err := r.conn.Get(r.path)
	if err != nil {
		return err
	}
	return r.conn.Delete(r.path, s.Version)
}

type ZkConf struct {
	conn    *zk.Conn
	path    string
	confmap map[string][]byte
}

func NewZkConf(conn *zk.Conn, path string) (*ZkConf, error) {
	children, _, err := conn.Children(path)
	if err != nil {
		return nil, err
	}

	tmpmap := make(map[string][]byte)
	for _, child := range children {
		var byteArr []byte
		byteArr, _, err = conn.Get(strings.Join([]string{path, "/", child}, ""))
		if err != nil {
			return nil, err
		}
		tmpmap[child] = byteArr
	}
	return &ZkConf{conn, path, tmpmap}, nil
}

func (zc *ZkConf) Get(name string) ([]byte, error) {
	val, ok := zc.confmap[name]
	if !ok {
		return nil, errors.New(name + "not found")
	}

	return val, nil
}

type KeepWatcher struct {
	conn *zk.Conn
}

func NewKeepWatcher(conn *zk.Conn) *KeepWatcher {
	return &KeepWatcher{conn}
}

func (w *KeepWatcher) WatchChildren(path string, listener func(children []string, err error)) {
	for {
		children, _, child_ch, err := w.conn.ChildrenW(path)
		listener(children, err)
		<-child_ch
	}
}

func (w *KeepWatcher) WatchData(path string, listener func(data []byte, err error)) {
	for {
		dataBuf, _, events, err := w.conn.GetW(path)
		listener(dataBuf, err)
		<-events
	}
}
