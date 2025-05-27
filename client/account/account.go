package account

type Budget struct {
	MonthlyLimit int64 `json:"monthly_limit"`
	TotalBudget  int64 `json:"total_budget"`
}

type Data struct {
	Transactions []Transaction `json:"transactions"`
	Budget       Budget        `json:"budget"`
	SyncToken    int64         `json:"sync_token"`
}

type Account struct {
	data    Data
	changes Changes
}

type Changes struct {
	AddedTransactions    []Transaction `json:"addedTransactions,omitempty"`
	DeletedTransactions  []Transaction `json:"deletedTransactions,omitempty"`
	ReplacedTransactions []Transaction `json:"replacedTransactions,omitempty"`
}

func NewAccount(d Data) (*Account, error) {
	a := Account{}
	a.data = d
	a.changes = Changes{}
	return &a, nil
}

func (a *Account) AppendTransaction(t Transaction) {
	a.data.Transactions = append(a.data.Transactions, t)
	a.changes.AddedTransactions = append(a.changes.AddedTransactions, t)
}

func (a *Account) GetData() Data {
	return a.data
}

func (a *Account) SetData(d Data) {
	a.data = d
}

func (a *Account) GetChanges() Changes {
	return a.changes
}
