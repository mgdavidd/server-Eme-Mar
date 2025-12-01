package models

type Account struct {
    Balance    float64 `json:"balance"`      // saldo que tengo
    AmountOwed float64 `json:"amount_owed"`  // saldo que me deben
}