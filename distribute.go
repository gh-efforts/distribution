package main

import (
	"encoding/json"
	"errors"
	"os"
)

var (
	orgsJson     string
	dataSetsJson string

	// 副本数量，默认10
	duplicate int
	// 单个SP单个piece重复的次数，正常最大为0
	repeat int
)

type User struct {
	Org string   `json:"org"`
	Sps []string `json:"sps"`
}

type Users struct {
	List []*User `json:"list"`
}

func NewUsers() *Users {
	return new(Users)
}

// ReadUsersFromFile 从JSON文件中读取Users结构体
func (u *Users) ReadUsersFromFile() error {
	// 读取文件内容
	content, err := os.ReadFile(orgsJson)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			_, err := os.Create(orgsJson)
			if err != nil {
				return err
			}
			if err := u.WriteUsersToFile(); err != nil {
				return err
			}
			return nil
		}
	}
	err = json.Unmarshal(content, &u)
	if err != nil {
		return err
	}

	return nil
}

// WriteUsersToFile 将Users结构体写入到JSON文件中
func (u *Users) WriteUsersToFile() error {
	// 序列化JSON数据
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}

	// 写入到文件
	err = os.WriteFile(orgsJson, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (u *Users) Add(user *User) {
	u.List = append(u.List, user)
}

func (u *Users) Get(org string) *User {
	for _, user := range u.List {
		if user.Org == org {
			return user
		}
	}
	return nil
}

// GetSps 用已知的sp获取到所在org的全部sp
func (u *Users) GetSps(inputSp string) []string {
	for _, user := range u.List {
		for _, sp := range user.Sps {
			if sp == inputSp {
				return user.Sps
			}
		}
	}
	return nil
}

func (u *Users) Update(update *User) {
	for i, user := range u.List {
		if user.Org == update.Org {
			u.List[i] = update
			break

		}
	}
}

func (u *Users) Delete(org string) bool {
	for i, user := range u.List {
		if user.Org == org {
			u.List = append(u.List[:i], u.List[i+1:]...)
			return true

		}
	}
	return false
}

func (u *Users) View() (string, error) {
	data, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type SpInfo struct {
	Sp  string `json:"sp"`
	Num int    `json:"num"`
}

type Piece struct {
	PieceCid  string    `json:"pieceCid"`
	PieceSize int64     `json:"pieceSize"`
	CarSize   int64     `json:"carSize"`
	SpInfos   []*SpInfo `json:"spInfos"`
}
type DataSet struct {
	Duplicate   int      `json:"duplicate"`
	DataSetName string   `json:"dataSetName"`
	Pieces      []*Piece `json:"pieces"`
}
type DataSets struct {
	List []*DataSet `json:"list"`
}

func NewDataSet() *DataSet {
	return new(DataSet)
}
func NewDataSets() *DataSets {
	return new(DataSets)
}

// ReadDataSetsFromFile 从JSON文件中读取DataSets结构体
func (d *DataSets) ReadDataSetsFromFile() error {
	// 读取文件内容
	content, err := os.ReadFile(dataSetsJson)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			_, err := os.Create(dataSetsJson)
			if err != nil {
				return err
			}
			if err := d.WriteDataSetsToFile(); err != nil {
				return err
			}
			return nil
		}
	}
	err = json.Unmarshal(content, &d)
	if err != nil {
		return err
	}

	return nil
}

// WriteDataSetsToFile 将DataSets结构体写入到JSON文件中
func (d *DataSets) WriteDataSetsToFile() error {
	// 序列化JSON数据
	data, err := json.Marshal(d)
	if err != nil {
		return err
	}

	// 写入到文件
	err = os.WriteFile(dataSetsJson, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
func (d *DataSets) AddDataSet(dataSet *DataSet) {
	d.List = append(d.List, dataSet)
}

func (d *DataSets) GetDataset(dataSetName string) *DataSet {
	for _, dateSet := range d.List {
		if dateSet.DataSetName == dataSetName {
			return dateSet
		}
	}
	return nil
}
func (d *DataSets) UpdateDataSet(update *DataSet) {
	for i, dataSet := range d.List {
		if dataSet.DataSetName == update.DataSetName {
			d.List[i] = update
			break

		}
	}
}

func (d *DataSets) DeleteDataSet(dataSetName string) bool {
	for i, dataSet := range d.List {
		if dataSet.DataSetName == dataSetName {
			d.List = append(d.List[:i], d.List[i+1:]...)
			return true

		}
	}
	return false
}

func (d *DataSets) View() (string, error) {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (d *DataSet) Add(piece *Piece) {
	d.Pieces = append(d.Pieces, piece)
}

func (d *DataSet) Get(pieceCid string) *Piece {
	for _, piece := range d.Pieces {
		if piece.PieceCid == pieceCid {
			return piece
		}
	}
	return nil
}

// GetSize 返回符合条件的Pieces,总的pieceSize,总的carSize.会判断副本的数量，组织内其他sp是否已经发送，已经发送的次数是否小于等于repeat
func (d *DataSet) GetSize(inputSp string, size int64, sps []string) (*DataSet, int64, int64) {
	outPieces := NewDataSet()
	pieceSize := int64(0)
	carSize := int64(0)

	if duplicate == 0 {
		duplicate = d.Duplicate
	}

	for _, piece := range d.Pieces {
		if pieceSize >= size {
			continue
		}
		if len(piece.SpInfos) >= duplicate {
			continue
		}
		var haveSP bool
		for _, spInfo := range piece.SpInfos {
			if spInfo.Sp == inputSp {
				haveSP = true
				if spInfo.Num <= repeat {
					spInfo.Num += 1
					pieceSize += piece.PieceSize
					carSize += piece.CarSize
					outPieces.Add(piece)

				}
				break
			} else {
				for _, sp := range sps {
					if spInfo.Sp == sp && spInfo.Num > repeat {
						haveSP = true
						break
					}
				}
			}
		}

		if !haveSP {
			pieceSize += piece.PieceSize
			carSize += piece.CarSize
			outPieces.Add(piece)
			piece.SpInfos = append(piece.SpInfos, &SpInfo{inputSp, 1})
		}

	}
	return outPieces, pieceSize, carSize
}

func (d *DataSet) Update(update *Piece) {
	for i, piece := range d.Pieces {
		if piece.PieceCid == update.PieceCid {
			d.Pieces[i] = update
			break

		}
	}
}

func (d *DataSet) Delete(pieceCid string) bool {
	for i, piece := range d.Pieces {
		if piece.PieceCid == pieceCid {
			d.Pieces = append(d.Pieces[:i], d.Pieces[i+1:]...)
			return true

		}
	}
	return false
}
