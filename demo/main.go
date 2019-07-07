package demo

import (
	"log"
        //"time"
        "sync"
        "bytes"
        //"strings"
        "encoding/binary"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	PROJECT_NAME_DB   *leveldb.DB
	PROJECT_MARK_DB   *leveldb.DB
	PROJECT_TAGS_DB   *leveldb.DB // 为项目打TAG?
	PROJECT_TIME_DB   *leveldb.DB
	PROJECT_MASTER_DB *leveldb.DB

	USER_NAME_DB      *leveldb.DB
	USER_MARK_DB      *leveldb.DB
	USER_PROJECT_DB   *leveldb.DB

	LIST_PROJECT_DB   *leveldb.DB // 项目列表 所有

	AUTOID_DB         *leveldb.DB // 每次都将ID+1并写入

	AUTOID_USER_CH    chan int64
	AUTOID_PROJECT_CH chan int64

	TOKEN_MAP         sync.Map
)

var Texts = "adclisaj"
func init() {
	// 初始化前检查剩余空间与权限

	var err error

	// USER
	USER_NAME_DB, err = leveldb.OpenFile("../data/user_name", nil)
	if err != nil { panic("USER_NAME_DB INIT ERROR") }

	USER_MARK_DB, err = leveldb.OpenFile("../data/user_mark", nil)
	if err != nil { panic("USER_MARK_DB INIT ERROR") }

	USER_PROJECT_DB, err = leveldb.OpenFile("../data/user_project", nil)
	if err != nil { panic("USER_PROJECT_DB INIT ERROR") }

	// PROJECT
	PROJECT_NAME_DB, err = leveldb.OpenFile("../data/project_name", nil)
	if err != nil { panic("PROJECT_NAME_DB INIT ERROR") }

	PROJECT_MARK_DB, err = leveldb.OpenFile("../data/project_mark", nil)
	if err != nil { panic("PROJECT_MARK_DB INIT ERROR") }

	PROJECT_TAGS_DB, err = leveldb.OpenFile("../data/project_tags", nil)
	if err != nil { panic("PROJECT_TAGS_DB INIT ERROR") }

	PROJECT_MASTER_DB, err = leveldb.OpenFile("../data/project_master", nil)
	if err != nil { panic("PROJECT_MASTER_DB INIT ERROR") }

	PROJECT_TIME_DB, err = leveldb.OpenFile("../data/project_time", nil)
	if err != nil { panic("PROJECT_TIME_DB INIT ERROR") }

	// LIST
	LIST_PROJECT_DB, err = leveldb.OpenFile("../data/list_project", nil)
	if err != nil { panic("LIST_PROJECT_DB INIT ERROR") }

	AUTOID_DB, err = leveldb.OpenFile("../data/autoid", nil)
	if err != nil { panic("AUTOID_DB INIT ERROR") }

	// 通道初始化
	AUTOID_USER_CH = make(chan int64)
	AUTOID_PROJECT_CH = make(chan int64)

	// 自增数值独立进程初始化
	go autoid("user", AUTOID_USER_CH)
	go autoid("project", AUTOID_PROJECT_CH)

	//USER_NAME_DB.Put([]byte("233"),[]byte("Last"), nil)
	//data, err := db.Get([]byte("key"), nil)
	//err = db.Put([]byte("key"), []byte("value"), nil)
	//err = db.Delete([]byte("key"), nil)
	//defer USER_DB.Close()
}


func autoid(name string, c chan int64) {
	buf := new(bytes.Buffer)

	data, err := AUTOID_DB.Get([]byte(name), nil)
	if err != nil {
		binary.Write(buf, binary.BigEndian, 0)
		data = buf.Bytes()
		err = AUTOID_DB.Put([]byte(name), data, nil)
		if err != nil { panic("AUTOID_DB " + name + " INIT ERROR") }
		log.Println("计数器初始化", name)
	}

	var sum int64
	binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &sum)
	log.Println("计数器", name, sum)
	for {
		sum++
		c <- sum
		buf = new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, sum)
		err = AUTOID_DB.Put([]byte(name), buf.Bytes(), nil)
		if err != nil { panic("AUTOID ++ ERROR") }
	}
}

type Object interface {
	Delete() bool
	Create(name []byte)
	Updata()
	Load()
}
type User    []byte
type Project []byte
type Tag     []byte

// user 固有的属性是不能增加删除的
// 作为一体化信息, 删除 user 整体痕迹
// 才能一并删除所有属性
// 但修改单个属性呢?
func (u User)Delete() bool {
	var err error
	err = USER_NAME_DB.Delete(u, nil)
	if err != nil {}
	err = USER_MARK_DB.Delete(u, nil)
	if err != nil {}
	return true
}
func (u User)Load() {}
func (u User)Updata() {}
func (u User)Rewrite(item, name []byte) {
	// ??? 修改并不是重写
	// 修改是对单一属性对象的改变
}
func (u User)Create(name []byte) {
	// 并不存在明确目标的事物创建
	// 操作本身不具有签名, 事物签名是后置的, 也可以选择不
	// 留言 刻画
	// 对象是有明确结构的 可复数的 生物性的
	var err error
	err = USER_NAME_DB.Put(u, name, nil)
	if err != nil {}
	err = USER_MARK_DB.Put(u, name, nil)
	if err != nil {}
}
func (u User)Token(name []byte) bool {
	return true
}
func (p Project)Delete() bool {
	var err error
	err = PROJECT_NAME_DB.Delete(p, nil)
	if err != nil {
		log.Println(err)
	}
	//解除关联, 收藏, 关注, 上级
	//删除子级
	return true
}
func (p Project)Load() {}
func (p Project)Create(name []byte) {}
func (p Project)Updata() {}
func (t Tag)Delete() bool {
	return true
}
func (t Tag)Create(name []byte) {}
