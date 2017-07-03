package model

import (
	"github.com/EGaaS/go-egaas-mvp/packages/dal/types"
	"github.com/shopspring/decimal"
)

type KeyStatus string

const (
	Pending  KeyStatus = "my_pending"
	Approved           = "approved"
)

type dltTransaction struct {
	Model
	SenderWalletID         types.DALInt64
	RecipientWalletID      types.DALInt64
	RecipientWalletAddress types.DALString
	Amount                 types.DALDecimal
	Comission              types.DALDecimal
	Time                   types.DALInt32
	Comment                types.DALString
	BlockID                types.DALInt32
	RbID                   types.DALInt32
}

type MyKeys struct {
	AddTime      int32
	Notification int8
	PublicKey    []byte
	PrivateKey   []byte
	PasswordHash string
	Status       KeyStatus
	MyTime       int32
	Time         int32
	BlockID      int32
}

type MyNodeKeys struct {
	AddTime    int32
	PublicKey  []byte
	PrivateKey []byte
	Status     KeyStatus
	MyTime     int32
	Time       int64
	BlockID    int32
	RbID       int32
}

type TransactionStatus struct {
	Hash      []byte
	Time      int32
	Type      int32
	WalletID  int64
	CitizenID int64
	BlockID   int32
	Error     string
}

type Confirmations struct {
	BlockID      int64
	Good         int32
	Bad          int32
	Time         int
	Transaction  int32
	Cur0lMinerId int32
	MaxMinerId   int32
}

type Currency struct {
	Name     string
	FullName string
	RbID     int32
}

type InfoBlock struct {
	Hash           []byte
	BlockID        int32
	StateID        int32
	WalletID       int64
	Time           int32
	Level          int8
	CurrentVersion string
	Sent           int8
}

type LogTransactions struct {
	Hash []byte
	Time int32
}

type MainLock struct {
	LockTime  int32
	SciptName string
	Info      string
	Uniq      int8
}

type FullNode struct {
	Host                  string
	WalletID              int64
	StateID               int32
	FinalDelegateWalletID int64
	FinalDelegateStateID  int64
	RbID                  int32
}

type RbFullNodes struct {
	RbID               int64
	FullNodeWalletJson []byte
	BlockID            int32
	PrevRbID           int64
}

type UpdFullNodes struct {
	Time int32
	RbID int64
}

type RbUpdFullNodes struct {
	Time     int32
	BlockID  int32
	PrevRbID int64
}

type QueueBlocks struct {
	Hash       []byte
	FullNodeId int32
	BlockID    int32
}

type QueueTx struct {
	Hash          []byte
	Data          []byte
	FromGate      int32
	TmpNodeUserId string
}

type Transactions struct {
	Hash       []byte
	Data       []byte
	Verified   int8
	Used       int8
	HignRate   int8
	ForSelfUse int8
	Type       int8
	WalletID   int64
	CitizenID  int64
	ThirdVar   int32
	Counter    int8
	Sent       int8
}

type DltWallets struct {
	WalletID           int64
	PublickKey0        []byte
	NodePublicKey      []byte
	LastForgingDataUpd int64
	Amount             decimal.Decimal
	Host               string
	AddressVote        string
	FuelRate           int64
	SpendingContract   string
	ConditionsChange   string
	RbID               int64
}

type GlobalApps struct {
	Name   string
	Done   int32
	Blocks string
}

type SystemRecognizedStates struct {
	Name             string
	StateID          int64
	Host             string
	NodePublicKey    []byte
	DelegateWalletID int64
	DelegateStateID  int32
	RbID             int64
}

type Install struct {
	Progress string
}

type Config struct {
	MyBlockID              int32
	DltWalletID            int64
	StateID                int32
	CitizenID              int64
	BadBlocks              string
	PoolTechWorks          int8
	AutoReload             int32
	SetupPassword          string
	SqliteDbURL            string
	FirstLoadBlockchainURL string
	FirstLoadBlockchain    string
	HTTPHost               string
	AutoUpdate             int8
	AutoUpdateURL          string
	AnalyticsDisabled      int8
	StatHost               string
}

type StopDaemons struct {
	StopTime int32
}

type IncorrectTx struct {
	Hash  []byte
	Time  int32
	Error string
}

type MigrationsHistory struct {
	ID          int32
	Version     int32
	DateApplied int32
}

type DltWalletsBuffer struct {
	Hash       []byte
	DelBlockID int64
	WalletID   int64
	Amount     decimal.Decimal
	BlockID    int64
}

type President struct {
	StateID   int32
	CitizenID int64
	StartTime int64
}

type CbHead struct {
	StateCode string
	CitizenID int64
}

type RollbackTransactions struct {
	BlockID   int64
	TxHash    []byte
	TableName string
	TableID   string
}

type UpdateFullNodes struct {
	Time int32
	RbID int64
}

type RbUpdateFullNodes struct {
	Time     int32
	BlockID  int32
	PrevRbID int64
}

type RollbackRbID struct {
	RbID    int64
	BlockID int64
	Data    string
}

type SystemParameters struct {
	Name       string
	Value      map[string]interface{}
	Conditions string
	RbID       int64
}

type GlobalMenu struct {
	Name       string
	Value      string
	Conditions string
	RbID       int64
}

type GlobalPages struct {
	Name       string
	Value      string
	Menu       string
	Conditions string
	RbID       int64
}

type GlobalSignatures struct {
	Name       string
	Value      map[string]interface{}
	Conditions string
	RbID       int64
}

type GlobalSmartContract struct {
	Name       string
	Value      []byte
	WalletID   int64
	Active     rune
	Conditions string
	Variables  []byte
	RbID       int64
}

type GlobalTables struct {
	Name                  []byte
	ColumnsAndPermissions map[string]interface{}
	Conditions            string
	RbID                  int64
}

type SystemStates struct {
	RbID int64
}

type SystemRestoreAccess struct {
	CitizenID int64
	StateID   int64
	Active    int64
	Time      int64
	Close     int64
	Secret    string
	RbID      int64
}
