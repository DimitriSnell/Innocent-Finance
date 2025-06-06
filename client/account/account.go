package account

type Budget struct {
	MonthlyLimit int64 `json:"monthly_limit"`
	TotalBudget  int64 `json:"total_budget"`
}

type Data struct {
	Donators     []*Donator    `json:"donators"`
	Transactions []Transaction `json:"transactions"`
	Budget       Budget        `json:"budget"`
	SyncToken    int64         `json:"sync_token"`
}

type Account struct {
	data                Data
	changes             Changes
	donatorMap          map[string]*Donator
	onSyncStatusChanged func(synced bool) //callback function for view
}

type Changes struct {
	AddedTransactions    []Transaction `json:"addedTransactions,omitempty"`
	DeletedTransactions  []Transaction `json:"deletedTransactions,omitempty"`
	ReplacedTransactions []Transaction `json:"replacedTransactions,omitempty"`
	AddedDonators        []Donator     `json:"addedDonator,omitempty"`
}

func NewAccount(d Data) (*Account, error) {
	a := Account{}
	a.data = d
	a.changes = Changes{}
	a.createDonatorMap()
	return &a, nil
}

func (a *Account) createDonatorMap() {
	a.donatorMap = make(map[string]*Donator)
	for _, d := range a.GetData().Donators {
		a.donatorMap[d.ID] = d
	}
}

func (a *Account) GetDonatorNameByID(id string) string {
	value, ok := a.donatorMap[id]
	if !ok {
		return "Unknown"
	}
	return value.Name
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

func (a *Account) AppendDonator(d *Donator) {
	a.data.Donators = append(a.data.Donators, d)
	a.changes.AddedDonators = append(a.changes.AddedDonators, *d)
	a.donatorMap[d.ID] = d
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

func (a *Account) ClearChanges() {
	a.changes = Changes{}
}

func (a *Account) MarkUnsynced() {
	if a.onSyncStatusChanged != nil {
		a.onSyncStatusChanged(false)
	}
}

func (a *Account) SetOnSyncStatusChanged(cb func(synced bool)) {
	a.onSyncStatusChanged = cb
}

func (a *Account) NotifySyncStatusChanged(b bool) {
	if a.onSyncStatusChanged != nil {
		a.onSyncStatusChanged(b)
	}
}
