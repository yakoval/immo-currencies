package update

import "immo-currency/project"

func convertToInternal(c *Currency) (currency *project.Currency, err error) {
	rate, err := c.Rate()
	if err != nil {
		return
	}

	currency = &project.Currency{
		Name: c.Name,
		Rate: rate,
	}
	return
}
