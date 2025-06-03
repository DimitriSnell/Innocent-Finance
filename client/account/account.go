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

// Add transaction to deleted transactions list even if its not found in memory if the transaction id exists nowhere then the server will just ignore it
func (a *Account) DeleteTransactionFromMemory(t Transaction) {
	a.changes.DeletedTransactions = append(a.changes.DeletedTransactions, t)
	idx := FindIndexByID(t.ID, a.GetData().Transactions)
	if idx == -1 {
		return
	}
	newList := append(a.GetData().Transactions[:idx], a.GetData().Transactions[idx+1:]...)
	a.SetTransactionData(newList)
}

func (a *Account) AppendTransaction(t Transaction) {
	a.data.Transactions = append(a.data.Transactions, t)
	a.changes.AddedTransactions = append(a.changes.AddedTransactions, t)
}

func (a *Account) GetData() Data {
	return a.data
}

func (a *Account) SetSyncToken(st int64) {
	a.data.SyncToken = st
}

func (a *Account) SetData(d Data) {
	a.data = d
}

func (a *Account) GetChanges() Changes {
	return a.changes
}

func (a *Account) SetTransactionData(Tlist []Transaction) {
	a.data.Transactions = Tlist
}
