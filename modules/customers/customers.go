package customers

import "database/sql"

type Customer struct {
	Id        int64  `json:"id,omitempty"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Contacted bool   `json:"contacted"`
}

type CustomerRepository struct {
	DB *sql.DB
}

func (cr CustomerRepository) OfId(id int) (*Customer, error) {
	customerRow := cr.DB.QueryRow("SELECT * FROM customers WHERE id = ?", id)
	customer := Customer{}

	err := customerRow.Scan(&customer.Id, &customer.Name, &customer.Role, &customer.Email, &customer.Phone, &customer.Contacted)

	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (cr CustomerRepository) All() ([]Customer, error) {
	rows, err := cr.DB.Query("SELECT * FROM customers")

	if err != nil {
		return nil, err
	}

	var customersSlice = []Customer{}

	for rows.Next() {
		var customerStruct Customer

		rows.Scan(&customerStruct.Id, &customerStruct.Name, &customerStruct.Role, &customerStruct.Email, &customerStruct.Phone, &customerStruct.Contacted)

		customersSlice = append(customersSlice, customerStruct)
	}

	return customersSlice, nil
}

func (cr CustomerRepository) Create(customer *Customer) error {
	result, err := cr.DB.Exec("INSERT INTO customers(name, role, email, phone, contacted) VALUES (?, ?, ?, ?, ?)", customer.Name, customer.Role, customer.Email, customer.Phone, customer.Contacted)

	if err != nil {
		return err
	}

	customer.Id, _ = result.LastInsertId()

	return nil
}

func (cr CustomerRepository) Delete(customer *Customer) error {
	_, error := cr.DB.Exec("DELETE FROM customers WHERE id = ?", customer.Id)
	return error
}

func (cr CustomerRepository) Update(customer *Customer) error {
	_, err := cr.DB.Exec("UPDATE customers SET name = ?, role = ?, email = ?, phone = ?, contacted = ? WHERE id = ?", customer.Name, customer.Role, customer.Email, customer.Phone, customer.Contacted, customer.Id)
	return err
}
